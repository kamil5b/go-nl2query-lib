package repository

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type LLMRepository interface {
	GenerateQuery(ctx context.Context, prompt string, contexts []model.Vector, additionalPrompts ...string) (*string, error)
}
