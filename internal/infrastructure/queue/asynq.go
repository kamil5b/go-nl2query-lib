package queue

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
)

const TypeIngestion = "ingestion:metadata"

type AsynqTaskQueue struct {
	client *asynq.Client
}

func NewAsynqTaskQueue(rc *redis.Client) *AsynqTaskQueue {
	opts := asynq.RedisClientOpt{
		Addr: rc.Options().Addr,
	}
	return &AsynqTaskQueue{
		client: asynq.NewClient(opts),
	}
}

func (q *AsynqTaskQueue) EnqueueIngestionTask(ctx context.Context, tenantID, dbURL string) error {
	payload := map[string]string{
		"tenant_id": tenantID,
		"db_url":    dbURL,
	}

	data, _ := json.Marshal(payload)
	task := asynq.NewTask(TypeIngestion, data)

	_, err := q.client.EnqueueContext(ctx, task)
	return err
}

func (q *AsynqTaskQueue) Start(ctx context.Context) error {
	return nil
}

func (q *AsynqTaskQueue) Stop() error {
	if q.client != nil {
		return q.client.Close()
	}
	return nil
}
