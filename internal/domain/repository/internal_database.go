package repository

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type InternalDatabaseRepository interface {
	Connect(ctx context.Context, dbURL string) error
	Close() error
	GetWorkspaceByTenantID(ctx context.Context, tenantID string) (*model.Workspace, error)
	UpsertWorkspace(ctx context.Context, workspace *model.Workspace) error
}
