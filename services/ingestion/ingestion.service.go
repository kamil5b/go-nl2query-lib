package ingestion

import "github.com/kamil5b/go-nl2query-lib/ports"

type IngestionConfig struct{}

type IngestionService struct {
	Config *IngestionConfig

	embedderRepo    ports.EmbedderRepository
	vectorStoreRepo ports.VectorStoreRepository
	statusRepo      ports.StatusRepository
}

func NewIngestionService(
	config *IngestionConfig,

	embedderRepo ports.EmbedderRepository,
	vectorStoreRepo ports.VectorStoreRepository,
	statusRepo ports.StatusRepository,
) *IngestionService {
	return &IngestionService{
		Config: config,

		embedderRepo:    embedderRepo,
		vectorStoreRepo: vectorStoreRepo,
		statusRepo:      statusRepo,
	}
}
