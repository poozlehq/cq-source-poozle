package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/poozlehq/cq-ticketing/client"
	"github.com/poozlehq/cq-ticketing/internal/ticketing"
)

func User() *schema.Table {
	return &schema.Table{
		Name:          "ticketing_user",
		Resolver:      fetchUser,
		Transform:     transformers.TransformWithStruct(&ticketing.User{}),
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
	}
}

func fetchUser(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	cl := meta.(*client.Client)

	collectionCursor := "/collections"
	collectionParams := url.Values{}
	collectionParams.Set("limit", strconv.FormatInt(cl.Spec.Limit, 10))
	var collections []ticketing.Collection

	for {
		collectionRet, collectionParams, err := cl.Services.GetCollection(ctx, collectionCursor, collectionParams)
		if err != nil {
			return err
		}

		collections = append(collections, collectionRet.Data...)
		if collectionParams == nil {
			break
		}
	}

	for _, collection := range collections {
		key := fmt.Sprintf("ticketing-user-%s-%s-%s", cl.Spec.WorkspaceId, cl.Spec.IntegrationAccountId, *collection.Id)
		p := url.Values{}
		p.Set("raw", "true")
		p.Set("limit", strconv.FormatInt(cl.Spec.Limit, 10))
		cursor := fmt.Sprintf("/%s/users", *collection.Id)
		for {
			ret, p, err := cl.Services.GetUser(ctx, cursor, p)
			cl.Logger().Info().Msg(fmt.Sprintf("params %s", p))

			if err != nil {
				return err
			}
			now := time.Now()
			for i := range ret.Data {
				ret.Data[i].CqCreatedAt = &now
				ret.Data[i].CqUpdatedAt = &now
				ret.Data[i].IntegrationAccountId = &cl.Spec.IntegrationAccountId
			}
			res <- ret.Data

			data, err := json.Marshal(ret.Data)
			if err != nil {
				cl.Logger().Error().Msg(fmt.Sprintf("Error marshaling data: %v", err))
			}
			cl.Logger().Info().Msg(fmt.Sprintf("response %s", data))
			if p == nil {
				break
			}
		}

		if err := cl.Backend.SetKey(ctx, key, time.Now().Format(time.RFC3339)); err != nil {
			return fmt.Errorf("failed to store state to backend: %w", err)
		}
	}

	return cl.Backend.Flush(ctx)
}
