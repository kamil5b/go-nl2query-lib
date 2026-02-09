package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

type WorkspaceService interface {
	GetByTenantID(ctx context.Context, tenantID string) (*model.Workspace, error)
	ListAll(ctx context.Context) ([]*model.Workspace, error)
	Delete(ctx context.Context, tenantID string) error
	SyncClientDatabase(ctx context.Context, dbUrl string) (*model.DatabaseMetadata, error)
}
