package workspace

import (
	"context"

	"github.com/kamil5b/go-nl2query-lib/domains"
	"github.com/kamil5b/go-nl2query-lib/ports"
)

func (ws *WorkspaceService) Delete(ctx context.Context, tenantID string) error {
	// Step 1: Check status
	status, _, err := ws.statusAdapter.GetStatus(ctx, tenantID)
	if err != nil {
		return err
	}

	if status == domains.StatusInProgress {
		return ports.StatusInProgressError
	}

	// Step 2: Delete workspace from internal database
	if err := ws.internalDatabaseAdapter.DeleteWorkspaceByTenantID(ctx, tenantID); err != nil {
		return err
	}

	return nil
}
