package ingestion

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/go-nl-sql/internal/domain/model"
	"github.com/kamil5b/go-nl-sql/internal/domain/service"
	"github.com/kamil5b/go-nl-sql/mocks"
	"github.com/stretchr/testify/require"
)

// @example usage:
//
//	func TestIngestionService_VectorizeAndStore(t *testing.T) {
//	    ingestion.UnitTestVectorizeAndStore(t, New(embedderRepo, vectorStoreRepo, statusRepo))
//	}
func UnitTestVectorizeAndStore(
	t *testing.T,
	svcImp func(
		embedderRepo *mocks.MockEmbedderRepository,
		vectorStoreRepo *mocks.MockVectorStoreRepository,
		statusRepo *mocks.MockStatusRepository,
	) service.IngestionService,
) {
	var (
		mockEmbedderRepo    *mocks.MockEmbedderRepository
		mockVectorStoreRepo *mocks.MockVectorStoreRepository
		mockStatusRepo      *mocks.MockStatusRepository
	)

	reqMetaData := &model.DatabaseMetadata{
		TenantID: "tenant_abc",
	}

	mockVector := []float32{0.1, 0.2, 0.3}

	tests := []struct {
		name        string
		metadata    *model.DatabaseMetadata
		prepareMock func()
		expectError error
	}{
		{
			name:     "success",
			metadata: reqMetaData,
			prepareMock: func() {
				mockStatusRepo.EXPECT().
					SetInProgress(gomock.Any(), reqMetaData.TenantID).
					Return(nil)

				mockEmbedderRepo.EXPECT().
					Embed(gomock.Any(), gomock.Any()).
					Return(mockVector, nil)

				mockVectorStoreRepo.EXPECT().
					Upsert(gomock.Any(), reqMetaData.TenantID, mockVector).
					Return(nil)

				mockStatusRepo.EXPECT().
					SetDone(gomock.Any(), reqMetaData.TenantID).
					Return(nil)
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockEmbedderRepo = mocks.NewMockEmbedderRepository(ctrl)
			mockVectorStoreRepo = mocks.NewMockVectorStoreRepository(ctrl)
			mockStatusRepo = mocks.NewMockStatusRepository(ctrl)

			svc := svcImp(
				mockEmbedderRepo,
				mockVectorStoreRepo,
				mockStatusRepo,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			err := svc.VectorizeAndStore(context.Background(), tt.metadata)

			if tt.expectError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
