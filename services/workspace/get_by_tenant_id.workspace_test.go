package workspace

import (
	"testing"

	"github.com/kamil5b/go-nl2query-lib/ports"
	workspaceTest "github.com/kamil5b/go-nl2query-lib/testsuites/workspace"
)

func TestWorkspaceService_GetByTenantID(t *testing.T) {
	workspaceTest.UnitTestGetByTenantID(t, func(
		internalDatabaseAdapter ports.InternalDatabasePort,
	) ports.WorkspaceService {
		return NewWorkspaceService(nil,
			nil,
			nil,
			internalDatabaseAdapter,
			nil,
			nil,
			nil,
		)
	})
}
