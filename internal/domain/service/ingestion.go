package service

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type IngestionService interface {
	VectorizeAndStore(ctx context.Context, metadata *model.DatabaseMetadata) error
}
