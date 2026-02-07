package service

import "github.com/kamil5b/go-nl-sql/internal/domain/model"

type ChecksumGenerator interface {
	Generate(metadata *model.DatabaseMetadata) (string, error)
}
