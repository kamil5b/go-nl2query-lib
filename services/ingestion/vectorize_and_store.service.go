package ingestion

import (
	"context"

	"github.com/kamil5b/go-nl2query-lib/domains"
	toon "github.com/toon-format/toon-go"
)

func (s *IngestionService) VectorizeAndStore(ctx context.Context, metadata *domains.DatabaseMetadata) error {
	// Set status to in progress
	if err := s.statusRepo.SetInProgress(ctx, metadata.TenantID); err != nil {
		return err
	}

	// Prepare content strings from tables for embedding
	contents := metadataToTOON(metadata)

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

func metadataToTOON(meta *domains.DatabaseMetadata) []string {
	if meta == nil {
		return nil
	}

	type Row struct {
		TenantID string `json:"tenant_id"`
		Table    string `json:"table"`
		Column   string `json:"column"`
		Type     string `json:"type"`
		Nullable bool   `json:"nullable"`
		Primary  bool   `json:"primary_key"`
		Foreign  bool   `json:"foreign_key"`
		Comment  string `json:"comment"`
	}

	var docs []string

	for _, t := range meta.Tables {
		for _, c := range t.Columns {

			row := Row{
				TenantID: meta.TenantID,
				Table:    t.Name,
				Column:   c.Name,
				Type:     c.Type,
				Nullable: c.Nullable,
				Primary:  c.IsPrimaryKey,
				Foreign:  c.IsForeignKey,
				Comment:  c.Comments,
			}

			encoded, err := toon.MarshalString([]Row{row}) // one column per embedding unit
			if err == nil {
				docs = append(docs, encoded)
			}
		}
	}

	return docs
}
