package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/cloudquery/plugin-sdk/v4/message"
	"github.com/cloudquery/plugin-sdk/v4/plugin"
	"github.com/cloudquery/plugin-sdk/v4/scheduler"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/state"
	"github.com/poozlehq/cq-ticketing/client"
	"github.com/poozlehq/cq-ticketing/internal/ticketing"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/rs/zerolog"
)

type Client struct {
	logger      zerolog.Logger
	config      client.Spec
	tables      schema.Tables
	scheduler   *scheduler.Scheduler
	backendConn *grpc.ClientConn
	services    *ticketing.Client

	plugin.UnimplementedDestination
}

const (
	maxMsgSize = 100 * 1024 * 1024 // 100 MiB
)

func (c *Client) Logger() *zerolog.Logger {
	return &c.logger
}

func (c *Client) Sync(ctx context.Context, options plugin.SyncOptions, res chan<- message.SyncMessage) error {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return err
	}

	var stateClient state.Client
	if options.BackendOptions == nil {
		c.logger.Info().Msg("No backend options provided, using no state backend")
		stateClient = &state.NoOpClient{}
		c.backendConn = nil
	} else {
		c.backendConn, err = grpc.DialContext(ctx, options.BackendOptions.Connection,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultCallOptions(
				grpc.MaxCallRecvMsgSize(maxMsgSize),
				grpc.MaxCallSendMsgSize(maxMsgSize),
			),
		)
		if err != nil {
			return fmt.Errorf("failed to dial grpc source plugin at %s: %w", options.BackendOptions.Connection, err)
		}
		stateClient, err = state.NewClient(ctx, c.backendConn, options.BackendOptions.TableName)
		if err != nil {
			return fmt.Errorf("failed to create state client: %w", err)
		}
		c.logger.Info().Str("table_name", options.BackendOptions.TableName).Msg("Connected to state backend")
	}

	schedulerClient := client.New(c.logger, c.config, c.services, stateClient)
	return c.scheduler.Sync(ctx, schedulerClient, tt, res, scheduler.WithSyncDeterministicCQID(options.DeterministicCQID))
}

func (c *Client) Tables(_ context.Context, options plugin.TableOptions) (schema.Tables, error) {
	tt, err := c.tables.FilterDfs(options.Tables, options.SkipTables, options.SkipDependentTables)
	if err != nil {
		return nil, err
	}
	return tt, nil
}

func (*Client) Close(_ context.Context) error {
	return nil
}

func Configure(_ context.Context, logger zerolog.Logger, specBytes []byte, opts plugin.NewClientOptions) (plugin.Client, error) {
	if opts.NoConnection {
		return &Client{
			logger: logger,
			tables: getTables(),
		}, nil
	}

	config := &client.Spec{}
	if err := json.Unmarshal(specBytes, config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal spec: %w", err)
	}
	config.SetDefaults()
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate spec: %w", err)
	}

	services, err := ticketing.New(ticketing.ClientOptions{
		Log: logger.With().Str("source", "cq-ticketing").Logger(),
		HC: &http.Client{
			Timeout: time.Duration(config.Timeout) * time.Second,
		},
		ApiKey:               config.ApiKey,
		WorkspaceId:          config.WorkspaceId,
		IntegrationAccountId: config.IntegrationAccountId,
		StartDate:            config.StartDate,
		MaxRetries:           config.MaxRetries,
		PageSize:             int(config.PageSize),
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		config: *config,
		logger: logger,
		scheduler: scheduler.NewScheduler(
			scheduler.WithLogger(logger),
		),
		services: services,
		tables:   getTables(),
	}, nil
}
