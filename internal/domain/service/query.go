package service

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type QueryService interface {
	PromptToQueryData(ctx context.Context, tenantID string, prompt string, withData bool) (*model.Query, error)
}
