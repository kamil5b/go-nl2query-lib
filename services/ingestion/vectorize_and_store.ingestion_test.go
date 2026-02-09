package ingestion

import (
	"testing"

	"github.com/kamil5b/go-nl2query-lib/ports"
	ingestionTest "github.com/kamil5b/go-nl2query-lib/testsuites/ingestion"
)

func TestIngestionService_VectorizeAndStore(t *testing.T) {
	ingestionTest.UnitTestVectorizeAndStore(t, func(embedderAdapter ports.EmbedderPort, vectorStoreAdapter ports.VectorStorePort, statusAdapter ports.StatusPort) ports.IngestionService {
		return NewIngestionService(nil, embedderAdapter, vectorStoreAdapter, statusAdapter)
	}, metadataToTOON)
}
