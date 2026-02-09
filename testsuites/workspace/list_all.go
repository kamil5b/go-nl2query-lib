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
//	func TestWorkspaceService_ListAll(t *testing.T) {
//	    workspace.UnitTestListAll(t, NewWorkspaceService(config, statusAdapter, clientDatabaseAdapter, internalDatabaseAdapter, encryptAdapter, hashAdapter, taskQueueService))
//	}
func UnitTestListAll(
	t *testing.T,
	svcImp func(
		internalDatabaseAdapter ports.InternalDatabasePort,
	) ports.WorkspaceService,
) {
	var (
		mockInternalDatabaseAdapter *mocks.MockInternalDatabasePort
	)

	mockWorkspace1 := &domains.Workspace{
		TenantID: "tenant_123",
		Checksum: "checksum_abc",
	}

	mockWorkspace2 := &domains.Workspace{
		TenantID: "tenant_456",
		Checksum: "checksum_def",
	}

	mockWorkspaceList := []*domains.Workspace{mockWorkspace1, mockWorkspace2}

	tests := []struct {
		name        string
		prepareMock func()
		expectError error
		expectData  []*domains.Workspace
	}{
		{
			name: "success list all workspaces",
			prepareMock: func() {
				mockInternalDatabaseAdapter.
					EXPECT().
					ListAllWorkspaces(gomock.Any()).
					Return(mockWorkspaceList, nil)
			},
			expectError: nil,
			expectData:  mockWorkspaceList,
		},
		{
			name: "success empty list",
			prepareMock: func() {
				mockInternalDatabaseAdapter.
					EXPECT().
					ListAllWorkspaces(gomock.Any()).
					Return([]*domains.Workspace{}, nil)
			},
			expectError: nil,
			expectData:  []*domains.Workspace{},
		},
		{
			name: "error listing workspaces",
			prepareMock: func() {
				mockInternalDatabaseAdapter.
					EXPECT().
					ListAllWorkspaces(gomock.Any()).
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

			result, err := svc.ListAll(context.Background())

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
