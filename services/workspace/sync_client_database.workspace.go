package workspace

import (
	"context"

	"github.com/kamil5b/go-nl2query-lib/domains"
	model "github.com/kamil5b/go-nl2query-lib/domains"
	"github.com/kamil5b/go-nl2query-lib/ports"
)

func (ws *WorkspaceService) SyncClientDatabase(ctx context.Context, dbUrl string) (*model.DatabaseMetadata, error) {
	// Step 1: Generate tenant ID from database URL
	tenantID := ws.hashAdapter.GenerateTenantID(dbUrl)

	// Step 2: Check status
	status, _, err := ws.statusAdapter.GetStatus(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	if status == domains.StatusInProgress {
		return nil, ports.StatusInProgressError
	}

	// Step 3: Encrypt the database URL
	encryptedDBUrl := ws.encryptAdapter.Encrypt(dbUrl)

	// Step 4: Connect to internal database
	if err := ws.internalDatabaseAdapter.Connect(ctx, encryptedDBUrl); err != nil {
		return nil, err
	}

	// Step 5: Check if workspace already exists
	existingWorkspace, err := ws.internalDatabaseAdapter.GetWorkspaceByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	var existingChecksum string
	if existingWorkspace != nil {
		existingChecksum = existingWorkspace.Checksum
	}

	// Step 6: Connect to client database (treat connection error as warning)
	if err := ws.clientDatabaseAdapter.Connect(ctx, dbUrl); err != nil {
		// Connection error is treated as a warning, return nil
		return nil, nil
	}

	// Step 7: Execute query to get database metadata
	metadata, err := ws.clientDatabaseAdapter.GetDatabaseMetadata(ctx)
	if err != nil {
		return nil, err
	}

	// Step 9: Generate checksum for the database
	newChecksum, err := ws.hashAdapter.GenerateChecksum(metadata)
	if err != nil {
		return nil, err
	}

	metadata.Checksum = newChecksum

	// Step 10: Check if checksum has changed
	if existingWorkspace != nil && newChecksum == existingChecksum {
		// No changes detected, return success without enqueueing task
		return nil, nil
	}

	// Step 11: Enqueue ingestion task if checksum changed or workspace is new
	if err := ws.taskQueueService.EnqueueIngestionTask(ctx, tenantID, dbUrl); err != nil {
		return nil, err
	}

	return metadata, nil
}
