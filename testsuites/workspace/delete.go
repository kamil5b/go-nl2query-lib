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
//	func TestWorkspaceService_Delete(t *testing.T) {
//	    workspace.UnitTestDelete(t, NewWorkspaceService(config, statusAdapter, clientDatabaseAdapter, internalDatabaseAdapter, encryptAdapter, hashAdapter, taskQueueService))
//	}
func UnitTestDelete(
	t *testing.T,
	svcImp func(
		statusAdapter ports.StatusPort,
		internalDatabaseAdapter ports.InternalDatabasePort,
	) ports.WorkspaceService,
) {
	var (
		mockStatusAdapter           *mocks.MockStatusPort
		mockInternalDatabaseAdapter *mocks.MockInternalDatabasePort
	)

	mockTenantID := "tenant_123"

	tests := []struct {
		name        string
		tenantID    string
		prepareMock func()
		expectError error
	}{
		{
			name:     "success delete workspace",
			tenantID: mockTenantID,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					DeleteWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(nil)
			},
			expectError: nil,
		},
		{
			name:     "error deleting workspace",
			tenantID: mockTenantID,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, nil)
				mockInternalDatabaseAdapter.
					EXPECT().
					DeleteWorkspaceByTenantID(gomock.Any(), mockTenantID).
					Return(errors.New("database error"))
			},
			expectError: errors.New("database error"),
		},
		{
			name:     "error status",
			tenantID: mockTenantID,
			prepareMock: func() {
				mockStatusAdapter.
					EXPECT().
					GetStatus(gomock.Any(), mockTenantID).
					Return(domains.StatusDone, nil, errors.New("status error"))
			},
			expectError: errors.New("status error"),
		},
		{
			name:     "error deleting workspace because status is in progress",
			tenantID: mockTenantID,
			prepareMock: func() {
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
			mockInternalDatabaseAdapter = mocks.NewMockInternalDatabasePort(ctrl)

			svc := svcImp(
				mockStatusAdapter,
				mockInternalDatabaseAdapter,
			)

			if tt.prepareMock != nil {
				tt.prepareMock()
			}

			err := svc.Delete(context.Background(), tt.tenantID)

			if tt.expectError != nil {
				require.Error(t, err)
				require.Equal(t, err, tt.expectError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
