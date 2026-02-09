package ingestion

import "github.com/kamil5b/go-nl2query-lib/ports"

type IngestionConfig struct{}

type IngestionService struct {
	Config *IngestionConfig

	embedderAdapter    ports.EmbedderPort
	vectorStoreAdapter ports.VectorStorePort
	statusAdapter      ports.StatusPort
}

func NewIngestionService(
	config *IngestionConfig,

	embedderAdapter ports.EmbedderPort,
	vectorStoreAdapter ports.VectorStorePort,
	statusAdapter ports.StatusPort,
) *IngestionService {
	return &IngestionService{
		Config: config,

		embedderAdapter:    embedderAdapter,
		vectorStoreAdapter: vectorStoreAdapter,
		statusAdapter:      statusAdapter,
	}
}
