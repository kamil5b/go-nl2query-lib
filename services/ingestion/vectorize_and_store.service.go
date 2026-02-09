package ingestion

import (
	"context"
	"errors"

	"github.com/kamil5b/go-nl2query-lib/domains"
)

func (s *IngestionService) VectorizeAndStore(ctx context.Context, metadata *domains.DatabaseMetadata) error {
	return errors.New("not implemented")
}
