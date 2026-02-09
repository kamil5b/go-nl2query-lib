package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

type LLMRepository interface {
	GenerateQuery(ctx context.Context, prompt string, contexts []model.Vector, additionalPrompts ...string) (*string, error)
}
