package query

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/go-nl-sql/internal/domain/model"
	"github.com/kamil5b/go-nl-sql/internal/domain/repository"
	"github.com/kamil5b/go-nl-sql/internal/domain/service"
	"github.com/kamil5b/go-nl-sql/mocks"
	"github.com/stretchr/testify/require"
)

// @example usage:
//
//	func TestQueryService_PromptToQueryData(t *testing.T) {
//	    workspace.UnitTestPromptToQueryData(t, New())
//	}
func UnitTestPromptToQueryData(
	t *testing.T,
	svcImp func(
		clientDatabaseRepo repository.ClientDatabaseRepository,
		internalDatabaseRepo repository.InternalDatabaseRepository,
		encryptRepo repository.EncryptRepository,
		embedderRepo repository.EmbedderRepository,
		vectorStoreRepo repository.VectorStoreRepository,
		LLMRepo repository.LLMRepository,
		QueryValidator repository.QueryValidatorRepository,
	) service.QueryService,
) {
	var (
		mockClientDatabaseRepo   *mocks.MockClientDatabaseRepository
		mockInternalDatabaseRepo *mocks.MockInternalDatabaseRepository
		mockEncryptRepo          *mocks.MockEncryptRepository
		mockEmbedderRepo         *mocks.MockEmbedderRepository
		mockVectorStoreRepo      *mocks.MockVectorStoreRepository
		mockLLMRepo              *mocks.MockLLMRepository
		mockQueryValidatorRepo   *mocks.MockQueryValidatorRepository
	)

	mockString := "mocked string prompt"
	mockTenantID := "tenant_123"
	mockEncryptedDBUrl := "encrypted_tenant_123"
	mockURL := "https://database.url/tenant_123"
	mockContent := "mocked relevant data"
	mockResult := []map[string]any{
		{"tenant_id": mockTenantID, "encryptedDBUrl": mockEncryptedDBUrl},
	}
	mockQueryResultErrSyntax := "```sql SELECT * FROM table; ```"
	mockQueryResult := "SELECT * FROM table;"

	mockVector := []float32{0.1, 0.2, 0.3}
	mockVentorEntity := []model.Vector{
		{
			TenantID:  mockTenantID,
			Embedding: mockVector,
			Content:   mockContent,
		},
	}

	tests := []struct {
		name        string
		prepareMock func()
		withData    bool
		expectError error
	}{
		{
			name:     "success with data and loops",
			withData: true,
			prepareMock: func() {
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return("mocked database schema", nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEncryptRepo.
					EXPECT().
					Decrypt(mockEncryptedDBUrl).
					Return(mockURL, nil)
				mockEmbedderRepo.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreRepo.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// Start loop of generating query until safe
				mockLLMRepo.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity).
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorRepo.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(errors.New("syntax error"))
				mockLLMRepo.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error")).
					Return(&mockQueryResult, nil)
				mockQueryValidatorRepo.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(nil)
				mockQueryValidatorRepo.
					EXPECT().
					ContainsDDLDML(mockQueryResult).
					Return(false)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), mockQueryResult).
					Return(gomock.Any(), nil)
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClientDatabaseRepo = mocks.NewMockClientDatabaseRepository(ctrl)
			mockInternalDatabaseRepo = mocks.NewMockInternalDatabaseRepository(ctrl)
			mockEncryptRepo = mocks.NewMockEncryptRepository(ctrl)
			mockEmbedderRepo = mocks.NewMockEmbedderRepository(ctrl)
			mockVectorStoreRepo = mocks.NewMockVectorStoreRepository(ctrl)
			mockLLMRepo = mocks.NewMockLLMRepository(ctrl)
			mockQueryValidatorRepo = mocks.NewMockQueryValidatorRepository(ctrl)

			svc := svcImp(
				mockClientDatabaseRepo,
				mockInternalDatabaseRepo,
				mockEncryptRepo,
				mockEmbedderRepo,
				mockVectorStoreRepo,
				mockLLMRepo,
				mockQueryValidatorRepo,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			_, err := svc.PromptToQueryData(context.Background(), mockTenantID, mockString, tt.withData)

			if tt.expectError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
