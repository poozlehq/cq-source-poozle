package client

import (
	"fmt"

	"github.com/cloudquery/plugin-sdk/v4/state"
	"github.com/poozlehq/cq-ticketing/internal/ticketing"
	"github.com/rs/zerolog"
)

type Client struct {
	logger zerolog.Logger
	Spec   *Spec

	Services *ticketing.Client
	Backend  state.Client

	StartData string
}

func (c *Client) ID() string {
	return fmt.Sprintf("%s:%s", c.Spec.WorkspaceId, c.Spec.IntegrationAccountId)
}

func (c *Client) Logger() *zerolog.Logger {
	return &c.logger
}

func New(logger zerolog.Logger, spec Spec, services *ticketing.Client, bk state.Client) *Client {
	c := &Client{
		logger:   logger,
		Services: services,
		Spec:     &spec,
		Backend:  bk,
	}
	return c
}
