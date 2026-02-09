package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

type ClientDatabasePort interface {
	Connect(ctx context.Context, dbURL string) error
	Close() error
	Execute(ctx context.Context, query string) ([]map[string]any, error)
	GetDatabaseMetadata(ctx context.Context) (*model.DatabaseMetadata, error)
	ExecuteDryRun(ctx context.Context, query string) error
}
