package hash

import (
	"encoding/hex"
	"encoding/json"

	"github.com/kamil5b/go-nl-sql/internal/domain/model"
	"github.com/zeebo/blake3"
)

func (g *Blake3HashGenerator) GenerateChecksum(metadata *model.DatabaseMetadata) (string, error) {
	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		return "", err
	}
	hash := blake3.Sum256(jsonMetadata)
	hashResult := hex.EncodeToString(hash[:])
	return hashResult, nil
}
