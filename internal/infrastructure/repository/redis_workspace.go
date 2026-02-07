package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
	"github.com/redis/go-redis/v9"
)

type RedisWorkspaceRepository struct {
	client *redis.Client
}

func NewRedisWorkspaceRepository(client *redis.Client) *RedisWorkspaceRepository {
	return &RedisWorkspaceRepository{client: client}
}

func (r *RedisWorkspaceRepository) Save(ctx context.Context, workspace *model.Workspace) error {
	key := fmt.Sprintf("workspace:%s", workspace.TenantID)
	data, err := json.Marshal(workspace)
	if err != nil {
		return fmt.Errorf("marshal workspace: %w", err)
	}
	return r.client.Set(ctx, key, data, 0).Err()
}

func (r *RedisWorkspaceRepository) GetByTenantID(ctx context.Context, tenantID string) (*model.Workspace, error) {
	key := fmt.Sprintf("workspace:%s", tenantID)
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, model.ErrWorkspaceNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get workspace: %w", err)
	}

	var workspace model.Workspace
	if err := json.Unmarshal([]byte(data), &workspace); err != nil {
		return nil, fmt.Errorf("unmarshal workspace: %w", err)
	}
	return &workspace, nil
}

func (r *RedisWorkspaceRepository) ListAll(ctx context.Context) ([]*model.Workspace, error) {
	keys, err := r.client.Keys(ctx, "workspace:*").Result()
	if err != nil {
		return nil, fmt.Errorf("list workspace keys: %w", err)
	}

	workspaces := make([]*model.Workspace, 0, len(keys))
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		var workspace model.Workspace
		if err := json.Unmarshal([]byte(data), &workspace); err != nil {
			continue
		}
		workspaces = append(workspaces, &workspace)
	}
	return workspaces, nil
}

func (r *RedisWorkspaceRepository) Update(ctx context.Context, workspace *model.Workspace) error {
	return r.Save(ctx, workspace)
}

func (r *RedisWorkspaceRepository) Delete(ctx context.Context, tenantID string) error {
	key := fmt.Sprintf("workspace:%s", tenantID)
	return r.client.Del(ctx, key).Err()
}
