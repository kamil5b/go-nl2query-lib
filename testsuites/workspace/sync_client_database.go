package workspace

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
//	func TestWorkspaceService_SyncClientDatabase(t *testing.T) {
//	    workspace.UnitTestSyncClientDatabase(t, New(embedderAdapter, vectorStoreAdapter, statusAdapter))
//	}
func UnitTestSyncClientDatabase(
	t *testing.T,
	svcImp func(
		statusAdapter ports.StatusPort,
		clientDatabaseAdapter ports.ClientDatabasePort,
		internalDatabaseAdapter ports.InternalDatabasePort,
		encryptAdapter ports.EncryptPort,
		hashAdapter ports.HashPort,
		taskQueueService ports.TaskQueuePort,
	) ports.WorkspaceService,
) {
	var (
		mockStatusAdapter           *mocks.MockStatusPort
		mockClientDatabaseAdapter   *mocks.MockClientDatabasePort
		mockInternalDatabaseAdapter *mocks.MockInternalDatabasePort
		mockEncryptAdapter          *mocks.MockEncryptPort
		mockHashAdapter             *mocks.MockHashPort
		mockTaskQueuePort           *mocks.MockTaskQueuePort
	)

	mockString := "mocked string URL"
	mockTenantID := "tenant_123"
	mockEncryptedDBUrl := "encrypted_tenant_123"
	mockChecksum := "checksum_abc"
	mockChecksum2 := "checksum_def"

	mockResult := []map[string]any{
		{"tenant_id": mockTenantID, "checksum": mockChecksum},
	}

	tests := []struct {
		name        string
		prepareMock func()
		expectError error
	}{
		{
			name: "success update",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockEncryptAdapter.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashAdapter.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum2, nil)
				mockTaskQueuePort.
					EXPECT().
					EnqueueIngestionTask(gomock.Any(), mockTenantID, mockString).
					Return(nil)

			},
			expectError: nil,
		},
		{
			name: "success create",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockEncryptAdapter.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, nil) // No existing record
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashAdapter.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum, nil)
				mockTaskQueuePort.
					EXPECT().
					EnqueueIngestionTask(gomock.Any(), mockTenantID, mockString).
					Return(nil)

			},
			expectError: nil,
		},
		{
			name: "err enqueue ingestion task",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockEncryptAdapter.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, nil) // No existing record
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashAdapter.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum, nil)
				mockTaskQueuePort.
					EXPECT().
					EnqueueIngestionTask(gomock.Any(), mockTenantID, mockString).
					Return(errors.New("err"))

			},
			expectError: errors.New("err"),
		},
		{
			name: "success with no changes",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockEncryptAdapter.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashAdapter.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum, nil) // Match existing checksum
			},
			expectError: nil,
		},
		{
			name: "error generate checksum",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockEncryptAdapter.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashAdapter.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return("", errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name: "success with warn",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockEncryptAdapter.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockClientDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(errors.New("can't connect to client database")) // Simulate connection error
			},
			expectError: nil,
		},
		{
			name: "err executing internal DB",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockEncryptAdapter.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name: "err connect internal DB",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockEncryptAdapter.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseAdapter.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name: "err get status",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusError, nil, errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name: "err status in progress",
			prepareMock: func() {
				mockHashAdapter.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusInProgress, nil, nil)
			},
			expectError: ports.StatusInProgressError,
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
			mockHashAdapter = mocks.NewMockHashPort(ctrl)
			mockTaskQueuePort = mocks.NewMockTaskQueuePort(ctrl)

			svc := svcImp(
				mockStatusAdapter,
				mockClientDatabaseAdapter,
				mockInternalDatabaseAdapter,
				mockEncryptAdapter,
				mockHashAdapter,
				mockTaskQueuePort,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			_, err := svc.SyncClientDatabase(context.Background(), mockString)

			if tt.expectError != nil {
				require.Error(t, err)
				require.Equal(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
