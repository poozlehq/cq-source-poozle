package plugin

import (
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	resources "github.com/poozlehq/cq-source-ticketing/resources/services"
)

func getTables() []*schema.Table {
	tables := []*schema.Table{
		resources.Collection(),
		resources.Comment(),
		resources.Ticket(),
		resources.User(),
		resources.Team(),
		resources.Tag(),
	}

	if err := transformers.TransformTables(tables); err != nil {
		panic(err)
	}
	for _, t := range tables {
		schema.AddCqIDs(t)
	}

	return tables
}
