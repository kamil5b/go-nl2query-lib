package query

import (
	"context"
	"errors"

	"github.com/kamil5b/go-nl2query-lib/domains"
)

func (s *QueryService) PromptToQueryData(ctx context.Context, tenantID string, prompt string, withData bool) (*domains.Query, error) {
	return nil, errors.New("not implemented")
}
