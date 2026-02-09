package workspace

import (
	"testing"

	"github.com/kamil5b/go-nl2query-lib/ports"
	workspaceTest "github.com/kamil5b/go-nl2query-lib/testsuites/workspace"
)

func TestWorkspaceService_Delete(t *testing.T) {
	workspaceTest.UnitTestDelete(t, func(
		statusAdapter ports.StatusPort,
		internalDatabaseAdapter ports.InternalDatabasePort,
	) ports.WorkspaceService {
		return NewWorkspaceService(nil,
			statusAdapter,
			nil,
			internalDatabaseAdapter,
			nil,
			nil,
			nil,
		)
	})
}
