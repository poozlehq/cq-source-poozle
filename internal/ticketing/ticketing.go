package ticketing

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/poozlehq/cq-source-ticketing/internal/httperror"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type Client struct {
	opts *ClientOptions

	lim *rate.Limiter
}

type ClientOptions struct {
	Log zerolog.Logger

	HC         HTTPDoer
	MaxRetries int64
	PageSize   int

	ApiKey               string
	WorkspaceId          string
	IntegrationAccountId string
	StartDate            string
}

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

func New(opts ClientOptions) (*Client, error) {
	return &Client{
		opts: &opts,
		lim:  rate.NewLimiter(rate.Limit(80), 120),
	}, nil
}

func (s *Client) request(ctx context.Context, integrationUrl string, params url.Values) (retResp *http.Response, retErr error) {
	if params == nil {
		params = url.Values{}
	}
	// params.Set("limit", strconv.FormatInt(int64(s.opts.PageSize), 10))
	params.Set("realtime", "true")

	tries := int64(0)

	log := s.opts.Log.With().Str("edge", integrationUrl).Interface("query_params", params).Logger()

	defer func() {
		if retErr != nil {
			log.Error().Err(retErr).Msg("request failed")
		} else if tries > 0 {
			log.Debug().Int64("num_tries", tries).Msg("success after tries")
		}
	}()

	for {
		if !s.lim.Allow() {
			log.Debug().Msg("waiting for rate limiter...")
			if err := s.lim.Wait(ctx); err != nil {
				return nil, err
			}
			log.Debug().Msg("wait complete")
		}

		r, wait, err := s.retryableRequest(ctx, integrationUrl, params)
		if err == nil {
			return r, nil
		}

		temporary := false
		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			temporary = true
		} else if he, ok := err.(httperror.Error); ok {
			temporary = he.Temporary()
		}
		if !temporary {
			return nil, fmt.Errorf("request failed with error: %w", err)
		}

		tries++
		if tries >= s.opts.MaxRetries {
			return nil, fmt.Errorf("exceeded max retries (%d): %w", s.opts.MaxRetries, err)
		}

		if wait == nil { // no retry-after returned, linear backoff
			w := time.Duration(tries) * 1 * time.Second
			wait = &w
		}

		log.Warn().Err(err).Float64("backoff_seconds", wait.Seconds()).Msg("retryable request failed, will retry")

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(*wait):
		}
	}
}

func (s *Client) retryableRequest(ctx context.Context, integrationUrl string, params url.Values) (*http.Response, *time.Duration, error) {
	log := s.opts.Log.With().Str("edge", integrationUrl).Interface("query_params", params).Logger()

	u := integrationUrl
	if strings.Contains(u, "?") {
		u += "&" + params.Encode()
	} else {
		u += "?" + params.Encode()
	}
	log.Trace().Str("url", u).Msg("requesting...")
	log.Info().Str("url", u).Msg("requesting...")

	var (
		body []byte
		err  error
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", s.opts.ApiKey))
	req.Header.Add("Content-type", "application/json")
	req.Header.Add("workspaceId", s.opts.WorkspaceId)
	req.Header.Add("integrationAccountId", s.opts.IntegrationAccountId)

	resp, err := s.opts.HC.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("do %s: %w", integrationUrl, err)
	}

	var wait *time.Duration
	if ra := resp.Header.Get("Retry-After"); ra != "" {
		rr, err := strconv.ParseFloat(ra, 64)
		if err != nil {
			log.Warn().Str("retry_after", ra).Err(err).Msg("Unknown Retry-After received")
		} else {
			t := time.Duration(rr) * time.Second
			wait = &t
		}
	}

	if resp.StatusCode != http.StatusOK {
		bdy, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		var bodyStr string
		if bdy != nil {
			bodyStr = string(bdy)
		}

		if bodyStr == "" {
			b, _ := json.Marshal(resp.Header)
			bodyStr = "headers: " + string(b)
		}

		return nil, wait, httperror.New(resp.StatusCode, http.MethodGet, integrationUrl, resp.Status, bodyStr)
	}

	return resp, wait, nil
}

func (s *Client) GetCollection(ctx context.Context, pageUrl string, params url.Values) (*CollectionResponse, url.Values, error) {
	var ret CollectionResponse

	resp, err := s.request(ctx, pageUrl, params)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		log.Warn().Err(err).Msg("Error decoding body response")
		return nil, nil, err
	}

	nextPage := getNextPage(ret.Meta, params)

	return &ret, nextPage, nil
}

func (s *Client) GetTicket(ctx context.Context, pageUrl string, params url.Values) (*TicketResponse, url.Values, error) {
	var ret TicketResponse

	log.Info().Str("cursor", pageUrl).Msg("This is the pageurl for tickets")

	resp, err := s.request(ctx, pageUrl, params)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		log.Warn().Err(err).Msg("Error decoding body response")
		return nil, nil, err
	}

	nextPage := getNextPage(ret.Meta, params)

	return &ret, nextPage, nil
}

func (s *Client) GetComment(ctx context.Context, pageUrl string, params url.Values) (*CommentResponse, url.Values, error) {
	var ret CommentResponse

	log.Info().Str("cursor", pageUrl).Msg("This is the pageurl for comment")

	resp, err := s.request(ctx, pageUrl, params)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		log.Warn().Err(err).Msg("Error decoding body response")
		return nil, nil, err
	}

	nextPage := getNextPage(ret.Meta, params)

	return &ret, nextPage, nil
}

func (s *Client) GetTag(ctx context.Context, pageUrl string, params url.Values) (*TagResponse, url.Values, error) {
	var ret TagResponse

	log.Info().Str("cursor", pageUrl).Msg("This is the pageurl for comment")

	resp, err := s.request(ctx, pageUrl, params)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		log.Warn().Err(err).Msg("Error decoding body response")
		return nil, nil, err
	}

	nextPage := getNextPage(ret.Meta, params)

	return &ret, nextPage, nil
}

func (s *Client) GetTeam(ctx context.Context, pageUrl string, params url.Values) (*TeamResponse, url.Values, error) {
	var ret TeamResponse

	log.Info().Str("cursor", pageUrl).Msg("This is the pageurl for comment")

	resp, err := s.request(ctx, pageUrl, params)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		log.Warn().Err(err).Msg("Error decoding body response")
		return nil, nil, err
	}

	nextPage := getNextPage(ret.Meta, params)

	return &ret, nextPage, nil
}

func (s *Client) GetUsers(ctx context.Context, pageUrl string, params url.Values) (*UsersResponse, url.Values, error) {
	var ret UsersResponse

	log.Info().Str("cursor", pageUrl).Msg("This is the pageurl for comment")

	resp, err := s.request(ctx, pageUrl, params)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		log.Warn().Err(err).Msg("Error decoding body response")
		return nil, nil, err
	}

	nextPage := getNextPage(ret.Meta, params)

	return &ret, nextPage, nil
}

func (s *Client) GetUser(ctx context.Context, pageUrl string, params url.Values) (*UserResponse, url.Values, error) {
	var ret UserResponse

	log.Info().Str("cursor", pageUrl).Msg("This is the pageurl for comment")

	resp, err := s.request(ctx, pageUrl, params)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&ret)
	if err != nil {
		log.Warn().Err(err).Msg("Error decoding body response")
		return nil, nil, err
	}

	return &ret, nil, nil
}

func getNextPage(meta Meta, params url.Values) url.Values {
	if meta.Cursors.Next != "" {
		params.Set("cursor", meta.Cursors.Next)
		return params
	}

	return nil
}
