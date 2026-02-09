package model

import (
	"errors"
	"time"
)

type WorkspaceStatus string

const (
	StatusInProgress WorkspaceStatus = "IN_PROGRESS"
	StatusDone       WorkspaceStatus = "DONE"
	StatusError      WorkspaceStatus = "ERROR"
)

type Workspace struct {
	TenantID       string
	EncryptedDBURL string
	Status         WorkspaceStatus
	Checksum       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

var (
	ErrWorkspaceNotFound     = errors.New("workspace not found")
	ErrInvalidDBURL          = errors.New("invalid database URL")
	ErrIngestionInProgress   = errors.New("ingestion in progress")
	ErrDDLDMLDetected        = errors.New("DDL/DML operation detected")
	ErrQueryGenerationFailed = errors.New("failed to generate valid Query after retries")
	ErrDatabaseUnreachable   = errors.New("database unreachable")
	ErrNoMetadata            = errors.New("no metadata found for workspace")
)
