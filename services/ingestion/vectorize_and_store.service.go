package ingestion

import (
	"context"

	"github.com/kamil5b/go-nl2query-lib/domains"
)

func (s *IngestionService) VectorizeAndStore(ctx context.Context, metadata *domains.DatabaseMetadata) error {
	// Set status to in progress
	if err := s.statusRepo.SetInProgress(ctx, metadata.TenantID); err != nil {
		return err
	}

	// Prepare content strings from tables for embedding
	var contents []string
	for _, table := range metadata.Tables {
		contents = append(contents, table.Name+" table")
	}

	// Embed the content batch
	embeddings, err := s.embedderRepo.EmbedBatch(ctx, contents)
	if err != nil {
		// Set error status and return
		_ = s.statusRepo.SetError(ctx, metadata.TenantID, err.Error())
		return err
	}

	// Create vector entities from embeddings
	vectors := make([]domains.Vector, len(embeddings))
	for i, embedding := range embeddings {
		vectors[i] = domains.Vector{
			TenantID:  metadata.TenantID,
			Embedding: embedding,
			Content:   contents[i],
		}
	}

	// Upsert vectors to the store
	if err := s.vectorStoreRepo.Upsert(ctx, metadata.TenantID, vectors); err != nil {
		// Set error status and return
		_ = s.statusRepo.SetError(ctx, metadata.TenantID, err.Error())
		return err
	}

	// Set status to done
	if err := s.statusRepo.SetDone(ctx, metadata.TenantID); err != nil {
		return err
	}

	return nil
}
