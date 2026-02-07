package repository

import "context"

type LLMRepository interface {
	GenerateSQL(ctx context.Context, prompt string, context string) (string, error)
	CorrectSQL(ctx context.Context, prompt string, invalidSQL string, error []string) (string, error)
}
