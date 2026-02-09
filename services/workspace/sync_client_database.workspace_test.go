package workspace

import (
	"testing"

	"github.com/kamil5b/go-nl2query-lib/ports"
	workspaceTest "github.com/kamil5b/go-nl2query-lib/testsuites/workspace"
)

func TestWorkspaceService_PromptToWorkspaceData(t *testing.T) {
	workspaceTest.UnitTestSyncClientDatabase(t, func(
		statusAdapter ports.StatusPort,
		clientDatabaseAdapter ports.ClientDatabasePort,
		internalDatabaseAdapter ports.InternalDatabasePort,
		encryptAdapter ports.EncryptPort,
		hashAdapter ports.HashPort,
		taskQueueService ports.TaskQueuePort,
	) ports.WorkspaceService {
		return NewWorkspaceService(nil,
			statusAdapter,
			clientDatabaseAdapter,
			internalDatabaseAdapter,
			encryptAdapter,
			hashAdapter,
			taskQueueService,
		)
	})
}
