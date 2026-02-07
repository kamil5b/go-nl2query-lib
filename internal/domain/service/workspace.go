package service

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type WorkspaceService interface {
	Save(ctx context.Context, workspace *model.Workspace) error
	GetByTenantID(ctx context.Context, tenantID string) (*model.Workspace, error)
	ListAll(ctx context.Context) ([]*model.Workspace, error)
	Update(ctx context.Context, workspace *model.Workspace) error
	Delete(ctx context.Context, tenantID string) error
}
