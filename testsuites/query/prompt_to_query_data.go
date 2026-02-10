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
		queryErrorLimit int,
		executionErrorLimit int,
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
	mockWorkspace := &domains.Workspace{
		TenantID:       mockTenantID,
		EncryptedDBURL: mockEncryptedDBUrl,
	}
	mockQueryResultErrSyntax := "```sql SELECT * FROM table; ```"
	mockQueryResultErr := "SELECT * FROM table;"
	mockQueryResult := "SELECT * FROM tables;"

	mockQueryErrorLimit := 2
	mockExecutionErrorLimit := 2
	mockVector := []float32{0.1, 0.2, 0.3}
	mockVectorEntity := []domains.Vector{
		{
			TenantID:  mockTenantID,
			Embedding: mockVector,
			Content:   mockContent,
		},
	}
	constToWarn := func(msg string) *string {
		return &msg
	}

	dataResult := map[string]any{
		"result": "success",
	}

	tests := []struct {
		name             string
		prepareMock      func()
		isReturningQuery *string
		isReturningData  map[string]any
		withData         bool
		warnMessage      *string
		expectError      error
	}{
		{
			name:             "success with data and full route with loops",
			withData:         true,
			isReturningQuery: &mockQueryResult,
			isReturningData:  dataResult,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
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
					Search(gomock.Any(), mockTenantID, mockVector, 10).
					Return(mockVectorEntity, nil)

				// Outer loop iteration 0
				// Inner loop iteration 0 - syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// Inner loop iteration 1 - still syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// After inner loop - generate with accumulated errors
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErr, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErr).
					Return(true, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					ContainsDDLDML(mockQueryResultErr).
					Return(false)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), mockQueryResultErr).
					Return(nil, errors.New("execution error"))

				// Outer loop iteration 1
				// Inner loop iteration 0 - syntax error with execution error from previous attempt
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErr, "execution error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// Inner loop iteration 1 - still syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// After inner loop - generate with accumulated errors
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(true, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					ContainsDDLDML(mockQueryResult).
					Return(false)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), mockQueryResult).
					Return(dataResult, nil)
			},
			expectError: nil,
		},
		{
			name:             "success with warn because exceeding execution limit",
			withData:         true,
			isReturningQuery: &mockQueryResultErrSyntax,
			warnMessage:      constToWarn(ports.QueryServiceWarnQueryGeneratedUnsafe),
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
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
					Search(gomock.Any(), mockTenantID, mockVector, 10).
					Return(mockVectorEntity, nil)

				// Outer loop iteration 0
				// Inner loop iteration 0 - syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// Inner loop iteration 1 - still syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// After inner loop - generate with accumulated errors
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErr, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErr).
					Return(true, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					ContainsDDLDML(mockQueryResultErr).
					Return(false)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), mockQueryResultErr).
					Return(nil, errors.New("execution error"))

				// Outer loop iteration 1
				// Inner loop iteration 0 - syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErr, "execution error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// Inner loop iteration 1 - still syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// After inner loop - generate with accumulated errors (this will be the final query returned)
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))
			},
			expectError: nil,
		},
		{
			name:             "success with warn because query keep error",
			withData:         true,
			isReturningQuery: &mockQueryResultErrSyntax,
			warnMessage:      constToWarn(ports.QueryServiceWarnQueryGeneratedUnsafe),
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
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
					Search(gomock.Any(), mockTenantID, mockVector, 10).
					Return(mockVectorEntity, nil)

				// Outer loop iteration 0
				// Inner loop iteration 0 - syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// Inner loop iteration 1 - still syntax error
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))

				// After inner loop - generate with accumulated errors
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity, mockQueryResultErrSyntax, "syntax error").
					Return(&mockQueryResultErrSyntax, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResultErrSyntax).
					Return(false, errors.New("syntax error"))
			},
			expectError: nil,
		},
		{
			name:             "success with data and full route with loops but final DDL/DML warning so no data execution",
			withData:         true,
			isReturningQuery: &mockQueryResult,
			warnMessage:      constToWarn(ports.QueryServiceWarnDDLDMLDetected),
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
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
					Search(gomock.Any(), mockTenantID, mockVector, 10).
					Return(mockVectorEntity, nil)
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(true, nil)
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(true, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					ContainsDDLDML(mockQueryResult).
					Return(true)
			},
			expectError: nil,
		},
		{
			name:             "success but fail to connect client database",
			withData:         true,
			isReturningQuery: &mockQueryResult,
			warnMessage:      constToWarn(ports.QueryServiceWarnWontExecuteClientDatabaseError),
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
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
					Search(gomock.Any(), mockTenantID, mockVector, 10).
					Return(mockVectorEntity, nil)
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(true, nil)
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(true, nil)
			},
			expectError: nil,
		},
		{
			name:             "success with data and full route with loops but expected no data",
			withData:         false,
			isReturningQuery: &mockQueryResult,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
				mockEmbedderAdapter.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreAdapter.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, 10).
					Return(mockVectorEntity, nil)

				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(true, nil)
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(&mockQueryResult, nil)
				mockQueryValidatorAdapter.
					EXPECT().
					IsSafe(mockQueryResult).
					Return(true, nil)
			},
			expectError: nil,
		},
		{
			name:     "error status in progress",
			withData: false,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusInProgress, nil, nil)
			},
			expectError: ports.StatusInProgressError,
		},
		{
			name:     "err get status",
			withData: false,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusError, nil, errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name:     "err connect internal DB",
			withData: false,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name:     "err executing internal DB get workspace",
			withData: false,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name:     "err decrypt database URL",
			withData: true,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
				mockEncryptAdapter.
					EXPECT().
					Decrypt(mockEncryptedDBUrl).
					Return("", errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name:     "err embed prompt",
			withData: false,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
				mockEmbedderAdapter.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(nil, errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name:     "err vector store search",
			withData: false,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
				mockEmbedderAdapter.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreAdapter.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, 10).
					Return(nil, errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name:     "err generate query initial",
			withData: false,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockTenantID).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
				mockEmbedderAdapter.
					EXPECT().
					Embed(gomock.Any(), mockString).
					Return(mockVector, nil)
				mockVectorStoreAdapter.
					EXPECT().
					Search(gomock.Any(), mockTenantID, mockVector, 10).
					Return(mockVectorEntity, nil)
				mockLLMAdapter.
					EXPECT().
					GenerateQuery(gomock.Any(), mockString, mockVectorEntity).
					Return(nil, errors.New("err"))
			},
			expectError: errors.New("err"),
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
				mockQueryErrorLimit,
				mockExecutionErrorLimit,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			res, msg, err := svc.PromptToQueryData(context.Background(), mockTenantID, mockString, tt.withData)

			if tt.expectError != nil {
				require.Error(t, err)
				require.Equal(t, err, tt.expectError)
				require.Equal(t, msg, tt.warnMessage)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.isReturningQuery, res.ResultQuery)
				require.Equal(t, tt.isReturningData, res.ResultData)
			}
		})
	}
}
