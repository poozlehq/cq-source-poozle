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
	"github.com/poozlehq/cq-source-poozle/client"
	"github.com/poozlehq/cq-source-poozle/internal/payments"
)

func Charge() *schema.Table {
	return &schema.Table{
		Name:          "payments_charge",
		Resolver:      fetchCharge,
		Transform:     transformers.TransformWithStruct(&payments.Charge{}),
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
		},
	}
}

func fetchCharge(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	cl := meta.(*client.Client)

	key := fmt.Sprintf("payments-charge-%s-%s", cl.Spec.WorkspaceId, cl.Spec.IntegrationAccountId)
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

	cursor := fmt.Sprintf("%s/charges", cl.Spec.Url)
	for {
		ret, p, err := cl.Services.GetCharges(ctx, cursor, p)
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

	return cl.Backend.Flush(ctx)
}
