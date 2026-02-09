package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

var (
	StatusInProgressError = model.GoNL2QueryError{
		StatusCode: 409,
		Message:    "Workspace ingestion is in progress",
	}
	StatusErrorError = model.GoNL2QueryError{
		StatusCode: 500,
		Message:    "Workspace ingestion has failed",
	}
)

type StatusPort interface {
	SetInProgress(ctx context.Context, tenantID string) error
	SetDone(ctx context.Context, tenantID string) error
	SetError(ctx context.Context, tenantID string, message string) error
	SetWarn(ctx context.Context, tenantID string, message string) error
	GetStatus(ctx context.Context, tenantID string) (model.WorkspaceStatus, *string, error)
	Clear(ctx context.Context, tenantID string) error
}
