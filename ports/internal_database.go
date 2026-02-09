package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

type InternalDatabasePort interface {
	Connect(ctx context.Context, dbURL string) error
	Close() error
	ListAllWorkspaces(ctx context.Context) ([]*model.Workspace, error)
	DeleteWorkspaceByTenantID(ctx context.Context, tenantID string) error
	GetWorkspaceByTenantID(ctx context.Context, tenantID string) (*model.Workspace, error)
	UpsertWorkspace(ctx context.Context, workspace *model.Workspace) error
}
