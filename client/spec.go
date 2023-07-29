package client

import "fmt"

type Spec struct {
	ApiKey               string `json:"api_key"`
	WorkspaceId          string `json:"workspace_id"`
	IntegrationAccountId string `json:"integration_account_id"`
	StartDate            string `json:"start_date"`
	Uid                  string `json:"uid"`
	Url                  string `json:"url"`

	Timeout    int64 `json:"timeout_secs,omitempty"`
	MaxRetries int64 `json:"max_retries,omitempty"`
	PageSize   int64 `json:"page_size,omitempty"`
	Limit      int64 `json:"limit,omitempty"`
}

func (s *Spec) Validate() error {
	if s.ApiKey == "" && len(s.ApiKey) == 0 {
		return fmt.Errorf("missing personal access token or app auth in configuration")
	}

	if len(s.WorkspaceId) == 0 && len(s.IntegrationAccountId) == 0 {
		return fmt.Errorf("missing workspace or integration account ID in configuration")
	}

	return nil
}

func (s *Spec) SetDefaults() {
	if s.Timeout < 1 {
		s.Timeout = 60
	}
	if s.MaxRetries < 1 {
		s.MaxRetries = 3
	}
	if s.Limit < 1 {
		s.Limit = 100
	}
}
