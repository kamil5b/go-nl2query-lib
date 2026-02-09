package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

const (
	WorkspaceServiceWarnUseExistingClientDatabaseError = "Will using existing stored schema because connection to client database could not be established."
)

type WorkspaceService interface {
	GetByTenantID(ctx context.Context, tenantID string) (*model.Workspace, error)
	ListAll(ctx context.Context) ([]*model.Workspace, error)
	Delete(ctx context.Context, tenantID string) error
	SyncClientDatabase(ctx context.Context, dbUrl string) (result *model.DatabaseMetadata, msg *string, err error)
}
