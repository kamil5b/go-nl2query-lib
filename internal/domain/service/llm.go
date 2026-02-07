package service

import "context"

type LLMService interface {
	GenerateSQL(ctx context.Context, prompt string, context string) (string, error)
	CorrectSQL(ctx context.Context, prompt string, invalidSQL string, error string) (string, error)
}
