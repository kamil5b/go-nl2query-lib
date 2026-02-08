package repository

import "github.com/kamil5b/go-nl-sql/internal/domain/model"

type HashRepository interface {
	GenerateChecksum(metadata *model.DatabaseMetadata) (string, error)
	GenerateTenantID(dbUrl string) string
}
