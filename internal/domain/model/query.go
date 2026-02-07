package model

import (
	"time"
)

type Query struct {
	TenantID    string
	ResultQuery string
	ResultData  map[string]any
	Message     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
