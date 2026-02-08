package repository

import "context"

type LLMRepository interface {
	GenerateQuery(ctx context.Context, prompt string, context string) (string, error)
	CorrectQuery(ctx context.Context, prompt string, invalidSQL string, errors []string) (string, error)
}
