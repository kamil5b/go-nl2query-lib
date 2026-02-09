package domains

import (
	"time"
)

type Query struct {
	TenantID    string
	ResultQuery *string
	ResultData  *map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
