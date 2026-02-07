package service

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type MetadataExtractor interface {
	Extract(ctx context.Context, dbURL string) (*model.DatabaseMetadata, error)
}
