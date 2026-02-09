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
		statusAdapter ports.StatusPort,
		clientDatabaseAdapter ports.ClientDatabasePort,
		internalDatabaseAdapter ports.InternalDatabasePort,
		encryptAdapter ports.EncryptPort,
		embedderAdapter ports.EmbedderPort,
		vectorStoreAdapter ports.VectorStorePort,
		LLMAdapter ports.LLMPort,
		QueryValidator ports.QueryValidatorPort,
	) ports.QueryService,
) {
	var (
		mockStatusAdapter           *mocks.MockStatusPort
		mockClientDatabaseAdapter   *mocks.MockClientDatabasePort
		mockInternalDatabaseAdapter *mocks.MockInternalDatabasePort
		mockEncryptAdapter          *mocks.MockEncryptPort
		mockEmbedderAdapter         *mocks.MockEmbedderPort
		mockVectorStoreAdapter      *mocks.MockVectorStorePort
		mockLLMAdapter              *mocks.MockLLMPort
		mockQueryValidatorAdapter   *mocks.MockQueryValidatorPort
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
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEncryptAdapter.
					EXPECT().
					Decrypt(mockEncryptedDBUrl).
					Return(mockURL, nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockURL).
					Return(nil)
				mockEmbedderAdapter.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreAdapter.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// === Start loop of generating query until safe ===
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity).
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(errors.New("syntax error"))
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error").Error()).
					Return(&mockQueryResultErr, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErr).
					Return(nil)
				// End loop
				mockQueryValidatorAdapter.
					EXPECT().
					ContainsDDLDML(mockQueryResultErr).
					Return(false)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), mockQueryResultErr).
					Return(nil, errors.New("execution error"))
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("execution error").Error()).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(nil)
				mockQueryValidatorAdapter.
					EXPECT().
					ContainsDDLDML(mockQueryResult).
					Return(false)
				mockClientDatabaseAdapter.
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
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEncryptAdapter.
					EXPECT().
					Decrypt(mockEncryptedDBUrl).
					Return(mockURL, nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockURL).
					Return(nil)
				mockEmbedderAdapter.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreAdapter.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// === Start loop of generating query until safe ===
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity).
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(errors.New("syntax error"))
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error").Error()).
					Return(&mockQueryResultErr, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErr).
					Return(nil)
				// End loop
				mockQueryValidatorAdapter.
					EXPECT().
					ContainsDDLDML(mockQueryResultErr).
					Return(false)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), mockQueryResultErr).
					Return(nil, errors.New("execution error"))
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("execution error").Error()).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(nil)
				mockQueryValidatorAdapter.
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
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEncryptAdapter.
					EXPECT().
					Decrypt(mockEncryptedDBUrl).
					Return(mockURL, nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockURL).
					Return(errors.New("connection error"))
				mockEmbedderAdapter.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreAdapter.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// === Start loop of generating query until safe ===
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity).
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(errors.New("syntax error"))
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error").Error()).
					Return(&mockQueryResultErr, nil)
				mockQueryValidatorAdapter.
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
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockEmbedderAdapter.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreAdapter.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, gomock.Any()).
					Return(mockVentorEntity, nil)
				// === Start loop of generating query until safe ===
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity).
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(errors.New("syntax error"))
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVentorEntity, errors.New("syntax error").Error()).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
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

			mockStatusAdapter = mocks.NewMockStatusPort(ctrl)
			mockClientDatabaseAdapter = mocks.NewMockClientDatabasePort(ctrl)
			mockInternalDatabaseAdapter = mocks.NewMockInternalDatabasePort(ctrl)
			mockEncryptAdapter = mocks.NewMockEncryptPort(ctrl)
			mockEmbedderAdapter = mocks.NewMockEmbedderPort(ctrl)
			mockVectorStoreAdapter = mocks.NewMockVectorStorePort(ctrl)
			mockLLMAdapter = mocks.NewMockLLMPort(ctrl)
			mockQueryValidatorAdapter = mocks.NewMockQueryValidatorPort(ctrl)

			svc := svcImp(
				mockStatusAdapter,
				mockClientDatabaseAdapter,
				mockInternalDatabaseAdapter,
				mockEncryptAdapter,
				mockEmbedderAdapter,
				mockVectorStoreAdapter,
				mockLLMAdapter,
				mockQueryValidatorAdapter,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			_, err := svc.PromptToQueryData(context.Background(), mockTenantID, mockString, tt.withData)

			if tt.expectError != nil {
				require.Error(t, err)
				require.Equal(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
