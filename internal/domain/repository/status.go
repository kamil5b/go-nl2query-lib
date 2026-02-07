package repository

import (
	"context"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
)

type StatusRepository interface {
	SetInProgress(ctx context.Context, tenantID string) error
	SetDone(ctx context.Context, tenantID string) error
	SetError(ctx context.Context, tenantID string, message string) error
	SetWarn(ctx context.Context, tenantID string, message string) error
	GetStatus(ctx context.Context, tenantID string) (model.WorkspaceStatus, string, error)
	Clear(ctx context.Context, tenantID string) error
}
