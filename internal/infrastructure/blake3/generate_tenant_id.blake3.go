package hash

import (
	"encoding/hex"
	"fmt"

	"github.com/zeebo/blake3"
)

func (g *Blake3HashGenerator) GenerateTenantID(input string) string {
	hash := blake3.Sum256([]byte(input))
	hashResult := hex.EncodeToString(hash[:8])
	return fmt.Sprintf("tenant_%s", hashResult)
}
