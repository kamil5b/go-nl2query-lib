package service

import "context"

type DatabaseConnection interface {
	Connect(ctx context.Context, dbURL string) error
	Close() error
	Ping(ctx context.Context) error
	Execute(ctx context.Context, query string) ([]map[string]interface{}, error)
	ExecuteDryRun(ctx context.Context, query string) error
}
