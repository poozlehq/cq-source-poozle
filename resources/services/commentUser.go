package resources

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/apache/arrow/go/v13/arrow"
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	"github.com/poozlehq/cq-source-ticketing/client"
	"github.com/poozlehq/cq-source-ticketing/internal/ticketing"
)

func CommentUser() *schema.Table {
	return &schema.Table{
		Name:          "ticketing_user",
		Resolver:      fetchCommentUser,
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

func fetchCommentUser(ctx context.Context, meta schema.ClientMeta, parent *schema.Resource, res chan<- interface{}) error {
	cl := meta.(*client.Client)

	comment, ok := parent.Item.(ticketing.Comment)
	if !ok {
		return fmt.Errorf("parent.Item is not of type *ticketing.Collection, it is of type %T", parent.Item)
	}

	p := url.Values{}
	cursor := fmt.Sprintf("%s/engine/users/%s", cl.Spec.Url, *comment.CreatedById)

	ret, _, err := cl.Services.GetUser(ctx, cursor, p)
	if err != nil {
		return err
	}
	now := time.Now()
	ret.Data.CqCreatedAt = &now
	ret.Data.CqUpdatedAt = &now
	ret.Data.IntegrationAccountId = &cl.Spec.IntegrationAccountId

	res <- ret.Data

	return cl.Backend.Flush(ctx)
}
