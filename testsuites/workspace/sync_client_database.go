package workspace

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/go-nl2query-lib/ports"
	"github.com/kamil5b/go-nl2query-lib/testsuites/mocks"
	"github.com/stretchr/testify/require"
)

// @example usage:
//
//	func TestWorkspaceService_SyncClientDatabase(t *testing.T) {
//	    workspace.UnitTestSyncClientDatabase(t, New(embedderRepo, vectorStoreRepo, statusRepo))
//	}
func UnitTestSyncClientDatabase(
	t *testing.T,
	svcImp func(
		clientDatabaseRepo ports.ClientDatabaseRepository,
		internalDatabaseRepo ports.InternalDatabaseRepository,
		encryptRepo ports.EncryptRepository,
		hashRepo ports.HashRepository,
		taskQueueService ports.TaskQueueService,
	) ports.WorkspaceService,
) {
	var (
		mockClientDatabaseRepo   *mocks.MockClientDatabaseRepository
		mockInternalDatabaseRepo *mocks.MockInternalDatabaseRepository
		mockEncryptRepo          *mocks.MockEncryptRepository
		mockHashRepo             *mocks.MockHashRepository
		mockTaskQueueService     *mocks.MockTaskQueueService
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
				mockHashRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockEncryptRepo.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashRepo.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum2, nil)
				mockTaskQueueService.
					EXPECT().
					EnqueueIngestionTask(gomock.Any(), mockTenantID, mockString).
					Return(nil)

			},
			expectError: nil,
		},
		{
			name: "success create",
			prepareMock: func() {
				mockHashRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockEncryptRepo.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, nil) // No existing record
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashRepo.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum, nil)
				mockTaskQueueService.
					EXPECT().
					EnqueueIngestionTask(gomock.Any(), mockTenantID, mockString).
					Return(nil)

			},
			expectError: nil,
		},
		{
			name: "err enqueue ingestion task",
			prepareMock: func() {
				mockHashRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockEncryptRepo.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, nil) // No existing record
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashRepo.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum, nil)
				mockTaskQueueService.
					EXPECT().
					EnqueueIngestionTask(gomock.Any(), mockTenantID, mockString).
					Return(errors.New("err"))

			},
			expectError: errors.New("err"),
		},
		{
			name: "success with no changes",
			prepareMock: func() {
				mockHashRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockEncryptRepo.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashRepo.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum, nil) // Match existing checksum
			},
			expectError: nil,
		},
		{
			name: "error generate checksum",
			prepareMock: func() {
				mockHashRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockEncryptRepo.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockHashRepo.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return("", errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name: "success with warn",
			prepareMock: func() {
				mockHashRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockEncryptRepo.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockResult, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(errors.New("can't connect to client database")) // Simulate connection error
			},
			expectError: nil,
		},
		{
			name: "err executing internal DB",
			prepareMock: func() {
				mockHashRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockEncryptRepo.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, errors.New("err"))
			},
			expectError: errors.New("err"),
		},
		{
			name: "err connect internal DB",
			prepareMock: func() {
				mockHashRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return(mockTenantID)
				mockEncryptRepo.
					EXPECT().
					Encrypt(mockString).
					Return(mockEncryptedDBUrl)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(errors.New("err"))
			},
			expectError: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClientDatabaseRepo = mocks.NewMockClientDatabaseRepository(ctrl)
			mockInternalDatabaseRepo = mocks.NewMockInternalDatabaseRepository(ctrl)
			mockEncryptRepo = mocks.NewMockEncryptRepository(ctrl)
			mockHashRepo = mocks.NewMockHashRepository(ctrl)
			mockTaskQueueService = mocks.NewMockTaskQueueService(ctrl)

			svc := svcImp(
				mockClientDatabaseRepo,
				mockInternalDatabaseRepo,
				mockEncryptRepo,
				mockHashRepo,
				mockTaskQueueService,
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
