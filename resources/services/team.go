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
	"github.com/poozlehq/cq-source-ticketing/client"
	"github.com/poozlehq/cq-source-ticketing/internal/ticketing"
)

func Team() *schema.Table {
	return &schema.Table{
		Name:          "ticketing_team",
		Resolver:      fetchTeam,
		Transform:     transformers.TransformWithStruct(&ticketing.Team{}),
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

func fetchTeam(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	cl := meta.(*client.Client)

	key := fmt.Sprintf("ticketing-team-%s-%s", cl.Spec.WorkspaceId, cl.Spec.IntegrationAccountId)
	p := url.Values{}
	p.Set("raw", "true")
	p.Set("limit", strconv.FormatInt(cl.Spec.Limit, 10))
	cursor := "/teams"
	for {
		ret, p, err := cl.Services.GetTeam(ctx, cursor, p)
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

	return cl.Backend.Flush(ctx)
}
