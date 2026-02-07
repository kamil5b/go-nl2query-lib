package service

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type VectorStore interface {
	Upsert(ctx context.Context, tenantID string, vectors []model.Vector) error
	Search(ctx context.Context, tenantID string, queryEmbedding []float32, limit int) ([]model.Vector, error)
	Delete(ctx context.Context, tenantID string) error
	Exists(ctx context.Context, tenantID string) (bool, error)
}
