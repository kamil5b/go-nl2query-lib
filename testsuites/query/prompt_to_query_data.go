package query

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/go-nl2query-lib/domains"
	"github.com/kamil5b/go-nl2query-lib/ports"
	"github.com/kamil5b/go-nl2query-lib/testsuites/mocks"
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
		statusRepo ports.StatusRepository,
		clientDatabaseRepo ports.ClientDatabaseRepository,
		internalDatabaseRepo ports.InternalDatabaseRepository,
		encryptRepo ports.EncryptRepository,
		embedderRepo ports.EmbedderRepository,
		vectorStoreRepo ports.VectorStoreRepository,
		LLMRepo ports.LLMRepository,
		QueryValidator ports.QueryValidatorRepository,
	) ports.QueryService,
) {
	var (
		mockStatusRepo           *mocks.MockStatusRepository
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
	mockQueryResultErr := "SELECT * FROM table;"
	mockQueryResult := "SELECT * FROM tables;"

	mockVector := []float32{0.1, 0.2, 0.3}
	mockVentorEntity := []domains.Vector{
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
			name:     "success with data and full route with loops",
			withData: true,
			prepareMock: func() {
				mockStatusRepo.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(nil, nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEncryptRepo.
					EXPECT().
					Decrypt(mockEncryptedDBUrl).
					Return(mockURL, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockURL).
					Return(nil)
				mockEmbedderRepo.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreRepo.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// === Start loop of generating query until safe ===
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
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error").Error()).
					Return(&mockQueryResultErr, nil)
				mockQueryValidatorRepo.
					EXPECT().
					IsSafe(mockQueryResultErr).
					Return(nil)
				// End loop
				mockQueryValidatorRepo.
					EXPECT().
					ContainsDDLDML(mockQueryResultErr).
					Return(false)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), mockQueryResultErr).
					Return(nil, errors.New("execution error"))
				mockLLMRepo.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("execution error").Error()).
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
					Return(nil, errors.New("execution error"))
			},
			expectError: nil,
		},
		{
			name:     "success with data and full route with loops but final DDL/DML warning so no data execution",
			withData: true,
			prepareMock: func() {
				mockStatusRepo.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(nil, nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEncryptRepo.
					EXPECT().
					Decrypt(mockEncryptedDBUrl).
					Return(mockURL, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockURL).
					Return(nil)
				mockEmbedderRepo.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreRepo.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// === Start loop of generating query until safe ===
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
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error").Error()).
					Return(&mockQueryResultErr, nil)
				mockQueryValidatorRepo.
					EXPECT().
					IsSafe(mockQueryResultErr).
					Return(nil)
				// End loop
				mockQueryValidatorRepo.
					EXPECT().
					ContainsDDLDML(mockQueryResultErr).
					Return(false)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), mockQueryResultErr).
					Return(nil, errors.New("execution error"))
				mockLLMRepo.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("execution error").Error()).
					Return(&mockQueryResult, nil)
				mockQueryValidatorRepo.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(nil)
				mockQueryValidatorRepo.
					EXPECT().
					ContainsDDLDML(mockQueryResult).
					Return(true) // This time contains DDL/DML, Warn will be returned but 200
			},
			expectError: nil,
		},
		{
			name:     "success with data and full route with loops but fail to connect client database",
			withData: true,
			prepareMock: func() {
				mockStatusRepo.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(nil, nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEncryptRepo.
					EXPECT().
					Decrypt(mockEncryptedDBUrl).
					Return(mockURL, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockURL).
					Return(errors.New("connection error"))
				mockEmbedderRepo.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreRepo.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// === Start loop of generating query until safe ===
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
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error").Error()).
					Return(&mockQueryResultErr, nil)
				mockQueryValidatorRepo.
					EXPECT().
					IsSafe(mockQueryResultErr).
					Return(nil)
				// End loop
			},
			expectError: nil,
		},
		{
			name:     "success with data and full route with loops but expected no data",
			withData: false,
			prepareMock: func() {
				mockStatusRepo.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(nil, nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEmbedderRepo.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreRepo.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// === Start loop of generating query until safe ===
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
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error").Error()).
					Return(&mockQueryResult, nil)
				mockQueryValidatorRepo.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(nil)
				// End loop
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStatusRepo = mocks.NewMockStatusRepository(ctrl)
			mockClientDatabaseRepo = mocks.NewMockClientDatabaseRepository(ctrl)
			mockInternalDatabaseRepo = mocks.NewMockInternalDatabaseRepository(ctrl)
			mockEncryptRepo = mocks.NewMockEncryptRepository(ctrl)
			mockEmbedderRepo = mocks.NewMockEmbedderRepository(ctrl)
			mockVectorStoreRepo = mocks.NewMockVectorStoreRepository(ctrl)
			mockLLMRepo = mocks.NewMockLLMRepository(ctrl)
			mockQueryValidatorRepo = mocks.NewMockQueryValidatorRepository(ctrl)

			svc := svcImp(
				mockStatusRepo,
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
