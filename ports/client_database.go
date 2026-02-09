package ports

import "context"

type ClientDatabaseRepository interface {
	Connect(ctx context.Context, dbURL string) error
	Close() error
	Execute(ctx context.Context, query string) ([]map[string]any, error)
	ExecuteDryRun(ctx context.Context, query string) error
}
