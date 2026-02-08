package workspace

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/kamil5b/go-nl-sql/internal/domain/repository"
	"github.com/kamil5b/go-nl-sql/internal/domain/service"
	"github.com/kamil5b/go-nl-sql/mocks"
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
		clientDatabaseRepo repository.DatabaseRepository,
		internalDatabaseRepo repository.DatabaseRepository,
		encryptorRepo repository.EncryptorRepository,
		taskQueueService service.TaskQueueService,
	) service.WorkspaceService,
) {
	var (
		mockClientDatabaseRepo   *mocks.MockDatabaseRepository
		mockInternalDatabaseRepo *mocks.MockDatabaseRepository
		mockEncryptorRepo        *mocks.MockEncryptorRepository
		mockTaskQueueService     *mocks.MockTaskQueueService
	)

	mockString := "mocked string URL"
	mockTenantID := "tenant_123"
	mockChecksum := "checksum_abc"
	mockChecksum2 := "checksum_def"

	mockResult := []map[string]any{
		{"tenant_id": mockTenantID, "checksum": mockChecksum},
	}

	tests := []struct {
		name        string
		clientDBURL string
		prepareMock func()
		expectError error
	}{
		{
			name:        "success update",
			clientDBURL: mockString,
			prepareMock: func() {
				mockEncryptorRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return([]byte(mockTenantID), nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(mockResult, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockEncryptorRepo.
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
			name:        "success create",
			clientDBURL: mockString,
			prepareMock: func() {
				mockEncryptorRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return([]byte(mockTenantID), nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(nil, nil) // No existing record
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockEncryptorRepo.
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
			name:        "success with no changes",
			clientDBURL: mockString,
			prepareMock: func() {
				mockEncryptorRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return([]byte(mockTenantID), nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(mockResult, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(nil)
				mockClientDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(gomock.Any(), nil)
				mockEncryptorRepo.
					EXPECT().
					GenerateChecksum(gomock.Any()).
					Return(mockChecksum, nil) // Match existing checksum
			},
			expectError: nil,
		},
		{
			name:        "success with warn",
			clientDBURL: mockString,
			prepareMock: func() {
				mockEncryptorRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return([]byte(mockTenantID), nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(mockResult, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(errors.New("can't connect to client database")) // Simulate connection error\
			},
			expectError: nil,
		},
		{
			name:        "success with warn",
			clientDBURL: mockString,
			prepareMock: func() {
				mockEncryptorRepo.
					EXPECT().
					GenerateTenantID(mockString).
					Return([]byte(mockTenantID), nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), gomock.Any()).
					Return(nil)
				mockInternalDatabaseRepo.
					EXPECT().
					Execute(gomock.Any(), gomock.Any()).
					Return(mockResult, nil)
				mockClientDatabaseRepo.
					EXPECT().
					Connect(gomock.Any(), mockString).
					Return(errors.New("can't connect to client database")) // Simulate connection error\
			},
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockClientDatabaseRepo = mocks.NewMockDatabaseRepository(ctrl)
			mockInternalDatabaseRepo = mocks.NewMockDatabaseRepository(ctrl)
			mockEncryptorRepo = mocks.NewMockEncryptorRepository(ctrl)
			mockTaskQueueService = mocks.NewMockTaskQueueService(ctrl)

			svc := svcImp(
				mockClientDatabaseRepo,
				mockInternalDatabaseRepo,
				mockEncryptorRepo,
				mockTaskQueueService,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			_, err := svc.SyncClientDatabase(context.Background(), tt.clientDBURL)

			if tt.expectError != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
