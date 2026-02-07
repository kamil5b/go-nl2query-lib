package repository

import (
	"context"
	"fmt"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
	"github.com/redis/go-redis/v9"
)

type RedisStatusRepository struct {
	client *redis.Client
}

func NewRedisStatusRepository(client *redis.Client) *RedisStatusRepository {
	return &RedisStatusRepository{client: client}
}

func (r *RedisStatusRepository) SetInProgress(ctx context.Context, tenantID string) error {
	key := fmt.Sprintf("status:%s", tenantID)
	value := fmt.Sprintf("%s:Ingestion in progress", model.StatusInProgress)
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisStatusRepository) SetDone(ctx context.Context, tenantID string) error {
	key := fmt.Sprintf("status:%s", tenantID)
	return r.client.Del(ctx, key).Err()
}

func (r *RedisStatusRepository) SetError(ctx context.Context, tenantID string, message string) error {
	key := fmt.Sprintf("status:%s", tenantID)
	value := fmt.Sprintf("%s:%s", model.StatusError, message)
	return r.client.Set(ctx, key, value, 0).Err()
}

func (r *RedisStatusRepository) GetStatus(ctx context.Context, tenantID string) (model.WorkspaceStatus, string, error) {
	key := fmt.Sprintf("status:%s", tenantID)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return model.StatusDone, "", nil
	}
	if err != nil {
		return "", "", err
	}

	var status model.WorkspaceStatus
	var message string
	_, err = fmt.Sscanf(data, "%s:%s", &status, &message)
	if err != nil {
		status = model.WorkspaceStatus(data)
		message = ""
	}

	return status, message, nil
}

func (r *RedisStatusRepository) Clear(ctx context.Context, tenantID string) error {
	key := fmt.Sprintf("status:%s", tenantID)
	return r.client.Del(ctx, key).Err()
}
