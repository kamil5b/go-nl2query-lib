package repository

import "context"

type DatabaseRepository interface { // Can be for Internal or Client use. Can be SQL nor NoSQL
	Connect(ctx context.Context, dbURL string) error
	Close() error
	Ping(ctx context.Context) error
	Execute(ctx context.Context, query string) ([]map[string]any, error)
	ExecuteDryRun(ctx context.Context, query string) error
}
