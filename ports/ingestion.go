package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

type IngestionService interface {
	VectorizeAndStore(ctx context.Context, metadata *model.DatabaseMetadata) error
}
