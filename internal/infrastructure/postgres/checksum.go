package postgres

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
	"github.com/zeebo/blake3"
)

type ChecksumGenerator struct{}

func NewChecksumGenerator() *ChecksumGenerator {
	return &ChecksumGenerator{}
}

func (g *ChecksumGenerator) Generate(metadata *model.DatabaseMetadata) (string, error) {
	data, err := json.Marshal(metadata)
	if err != nil {
		return "", fmt.Errorf("marshal metadata: %w", err)
	}

	hash := blake3.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}
