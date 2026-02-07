package repository

import "github.com/kamil5b/go-nl-sql/internal/domain/model"

type EncryptorRepository interface {
	GenerateChecksum(metadata *model.DatabaseMetadata) (string, error)
	GenerateTenantID(dbUrl string) (string, error)
}
