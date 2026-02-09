package ports

import "context"

type TaskQueuePort interface {
	EnqueueIngestionTask(ctx context.Context, tenantID string, dbURL string) error
}
