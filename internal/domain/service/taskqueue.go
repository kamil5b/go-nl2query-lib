package service

import "context"

type TaskQueue interface {
	EnqueueIngestionTask(ctx context.Context, tenantID string, dbURL string) error
}
