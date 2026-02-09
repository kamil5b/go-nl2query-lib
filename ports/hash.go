package ports

import model "github.com/kamil5b/go-nl2query-lib/domains"

type HashRepository interface {
	GenerateChecksum(metadata *model.DatabaseMetadata) (string, error)
	GenerateTenantID(dbUrl string) string
}
