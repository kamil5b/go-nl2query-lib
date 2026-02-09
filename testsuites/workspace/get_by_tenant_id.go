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
//	func TestWorkspaceService_GetByTenantID(t *testing.T) {
//	    workspace.UnitTestGetByTenantID(t, NewWorkspaceService(config, statusAdapter, clientDatabaseAdapter, internalDatabaseAdapter, encryptAdapter, hashAdapter, taskQueueService))
//	}
func UnitTestGetByTenantID(
	t *testing.T,
	svcImp func(
		internalDatabaseAdapter ports.InternalDatabasePort,
	) ports.WorkspaceService,
) {
	var (
		mockInternalDatabaseAdapter *mocks.MockInternalDatabasePort
	)

	mockTenantID := "tenant_123"
	mockChecksum := "checksum_abc"

	mockWorkspace := &domains.Workspace{
		TenantID: mockTenantID,
		Checksum: mockChecksum,
	}

	tests := []struct {
		name        string
		tenantID    string
		prepareMock func()
		expectError error
		expectData  *domains.Workspace
	}{
		{
			name:     "success get workspace",
			tenantID: mockTenantID,
			prepareMock: func() {
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(mockWorkspace, nil)
			},
			expectError: nil,
			expectData:  mockWorkspace,
		},
		{
			name:     "success workspace not found",
			tenantID: mockTenantID,
			prepareMock: func() {
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, nil)
			},
			expectError: nil,
			expectData:  nil,
		},
		{
			name:     "error getting workspace",
			tenantID: mockTenantID,
			prepareMock: func() {
				mockInternalDatabaseAdapter.
					EXPECT().
					GetWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil, errors.New("database error"))
			},
			expectError: errors.New("database error"),
			expectData:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockInternalDatabaseAdapter = mocks.NewMockInternalDatabasePort(ctrl)

			svc := svcImp(
				mockInternalDatabaseAdapter,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			result, err := svc.GetByTenantID(context.Background(), tt.tenantID)

			if tt.expectError != nil {
				require.Error(t, err)
				require.Equal(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.expectData, result)
			}
		})
	}
}
