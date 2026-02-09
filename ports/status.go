package ports

import (
	"context"
)

type StatusPort interface {
	SetInProgress(ctx context.Context, tenantID string) error
	SetDone(ctx context.Context, tenantID string) error
	SetError(ctx context.Context, tenantID string, message string) error
	SetWarn(ctx context.Context, tenantID string, message string) error
	GetStatus(ctx context.Context, tenantID string) (*string, error)
	Clear(ctx context.Context, tenantID string) error
}
