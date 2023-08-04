package resources

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/poozlehq/cq-source-ticketing/client"
	"github.com/poozlehq/cq-source-ticketing/internal/ticketing"
)

func Ticket() *schema.Table {
	return &schema.Table{
		Name:          "ticketing_ticket",
		Resolver:      fetchTicket,
		Transform:     transformers.TransformWithStruct(&ticketing.Ticket{}),
		IsIncremental: true,
		Columns: []schema.Column{
			{
				Name:       "id",
				Type:       arrow.BinaryTypes.String,
				Resolver:   schema.PathResolver("Id"),
				PrimaryKey: true,
			},
			{
				Name:       "integration_account_id",
				Type:       arrow.BinaryTypes.String,
				Resolver:   schema.PathResolver("IntegrationAccountId"),
				PrimaryKey: true,
			},
			{
				Name:           "created_at",
				Type:           arrow.FixedWidthTypes.Timestamp_us,
				Resolver:       schema.PathResolver("CreatedAt"),
				IncrementalKey: true,
			},
			{
				Name:           "updated_at",
				Type:           arrow.FixedWidthTypes.Timestamp_us,
				Resolver:       schema.PathResolver("UpdatedAt"),
				IncrementalKey: true,
			},
		},
		Relations: []*schema.Table{
			Comment(),
		},
	}
}

func fetchTicket(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	cl := meta.(*client.Client)

	collectionCursor := fmt.Sprintf("%s/collections", cl.Spec.Url)
	collectionParams := url.Values{}
	collectionParams.Set("limit", strconv.FormatInt(cl.Spec.Limit, 10))
	var collections []ticketing.Collection

	for {
		collectionRet, collectionParams, err := cl.Services.GetCollection(ctx, collectionCursor, collectionParams)
		if err != nil {
			return fmt.Errorf("tickets collection error: %s", err)
		}

		values := []string{"tensorflow", "airbyte", "kubernetes"}

		for _, collection := range collectionRet.Data {
			for _, value := range values {
				if *collection.Id == value {
					collections = append(collections, collection)
					break
				}
			}
		}

		// collections = append(collections, collectionRet.Data...)

		if collectionParams == nil {
			break
		}
	}

	for _, collection := range collections {
		key := fmt.Sprintf("ticketing-ticket-%s-%s-%s", cl.Spec.WorkspaceId, cl.Spec.IntegrationAccountId, *collection.Id)
		p := url.Values{}

		min, _ := time.Parse(time.RFC3339, cl.Spec.StartDate)
		cl.Logger().Info().Msg(fmt.Sprintf("backend %s", cl.Backend))
		if cl.Backend != nil {

			value, err := cl.Backend.GetKey(ctx, key)

			cl.Logger().Info().Msg(fmt.Sprintf("backend value %s", value))
			if err != nil {
				return fmt.Errorf("failed to retrieve state from backend: %w", err)
			}
			if value != "" {
				min, err = time.Parse(time.RFC3339, value)
				if err != nil {
					return fmt.Errorf("retrieved invalid state value: %q %w", value, err)
				}
			}
		}
		p.Set("created_after", min.Format(time.RFC3339))
		p.Set("raw", "true")
		p.Set("limit", strconv.FormatInt(cl.Spec.Limit, 10))
		cursor := fmt.Sprintf("%s/%s/tickets", cl.Spec.Url, *collection.Id)
		for {
			ret, p, err := cl.Services.GetTicket(ctx, cursor, p)
			if err != nil {
				cl.Logger().Err(err)
				return cl.Backend.Flush(ctx)
			}
			now := time.Now()
			for i := range ret.Data {
				ret.Data[i].CqCreatedAt = &now
				ret.Data[i].CqUpdatedAt = &now
				ret.Data[i].IntegrationAccountId = &cl.Spec.IntegrationAccountId
			}
			res <- ret.Data

			if p == nil {
				break
			}

		}

		err := cl.Backend.SetKey(ctx, key, time.Now().Format(time.RFC3339))

		if err != nil {
			return fmt.Errorf("failed to set state backend: %w", err)
		}
	}
	return cl.Backend.Flush(ctx)
}
