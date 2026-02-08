package service

import "context"

type TaskQueueService interface {
	EnqueueIngestionTask(ctx context.Context, tenantID string, dbURL string) error
}
