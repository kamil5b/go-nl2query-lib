package workspace

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

func (ws *WorkspaceService) ListAll(ctx context.Context) ([]*model.Workspace, error) {
	return ws.internalDatabaseAdapter.ListAllWorkspaces(ctx)
}
