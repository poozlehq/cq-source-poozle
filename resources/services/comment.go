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

func Comment() *schema.Table {
	return &schema.Table{
		Name:      "ticketing_comment",
		Resolver:  fetchComment,
		Transform: transformers.TransformWithStruct(&ticketing.Comment{}),
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
				Name:           "updated_at",
				Type:           arrow.FixedWidthTypes.Timestamp_us,
				Resolver:       schema.PathResolver("UpdatedAt"),
				IncrementalKey: true,
			},
		},
		IsIncremental: true,
	}
}

func fetchComment(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	cl := meta.(*client.Client)

	ticket, ok := parent.Item.(ticketing.Ticket)
	if !ok {
		return fmt.Errorf("parent.Item is not of type *ticketing.Collection, it is of type %T", parent.Item)
	}

	p := url.Values{}

	p.Set("raw", "true")
	p.Set("limit", strconv.FormatInt(cl.Spec.Limit, 10))

	cursor := fmt.Sprintf("%s/%s/tickets/%s/comments", cl.Spec.Url, *ticket.CollectionId, *ticket.Id)
	for {
		ret, p, err := cl.Services.GetComment(ctx, cursor, p)
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

		if p == nil {
			break
		}
	}

	return nil
}
