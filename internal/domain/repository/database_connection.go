package repository

import "context"

type ClientDatabaseRepository interface {
	Connect(ctx context.Context, dbURL string) error
	Close() error
	Ping(ctx context.Context) error
	Execute(ctx context.Context, query string) ([]map[string]any, error)
	ExecuteDryRun(ctx context.Context, query string) error
}
