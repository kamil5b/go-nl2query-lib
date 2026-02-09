package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

type QueryService interface {
	PromptToQueryData(ctx context.Context, tenantID string, prompt string, withData bool) (*model.Query, error)
}
