package service

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type QueryService interface {
	PromptToQuery(ctx context.Context, prompt string) (string, error)
	QueryToData(ctx context.Context, query string) (map[string]any, error)
	PromptToData(ctx context.Context, prompt string) (map[string]any, error)
	PromptToQueryData(ctx context.Context, prompt string) (*model.Query, error)
}
