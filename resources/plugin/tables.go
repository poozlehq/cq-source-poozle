package plugin

import (
	"github.com/cloudquery/plugin-sdk/v4/schema"
	"github.com/cloudquery/plugin-sdk/v4/transformers"
	ticketing "github.com/poozlehq/cq-source-ticketing/resources/services/ticketing"
)

func getTables() []*schema.Table {
	tables := []*schema.Table{
		ticketing.Collection(),
		ticketing.Comment(),
		ticketing.Ticket(),
		ticketing.User(),
		ticketing.Team(),
		ticketing.Tag(),
	}

	if err := transformers.TransformTables(tables); err != nil {
		panic(err)
	}
	for _, t := range tables {
		schema.AddCqIDs(t)
	}

	return tables
}
