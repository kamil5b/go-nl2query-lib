package hash

import (
	"encoding/hex"

	"github.com/zeebo/blake3"
)

type Blake3HashGenerator struct{}

func New() *Blake3HashGenerator {
	return &Blake3HashGenerator{}
}

func (g *Blake3HashGenerator) Generate(input string) (string, error) {
	hash := blake3.Sum256([]byte(input))
	return hex.EncodeToString(hash[:]), nil
}
