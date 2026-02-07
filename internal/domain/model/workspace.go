package model

import (
	"errors"
	"fmt"
	"time"
)

type WorkspaceStatus string

const (
	StatusInProgress WorkspaceStatus = "IN_PROGRESS"
	StatusDone       WorkspaceStatus = "DONE"
	StatusError      WorkspaceStatus = "ERROR"
)

type Workspace struct {
	TenantID  string
	DBURL     string
	Status    WorkspaceStatus
	Message   string
	Checksum  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewWorkspace(tenantID string, dbURL string) *Workspace {
	now := time.Now()
	return &Workspace{
		TenantID:  tenantID,
		DBURL:     dbURL,
		Status:    StatusInProgress,
		Message:   "Ingestion started",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (w *Workspace) SetInProgress() {
	w.Status = StatusInProgress
	w.Message = "Ingestion in progress"
	w.UpdatedAt = time.Now()
}

func (w *Workspace) SetDone() {
	w.Status = StatusDone
	w.Message = "Ingestion completed"
	w.UpdatedAt = time.Now()
}

func (w *Workspace) SetError(err error) {
	w.Status = StatusError
	w.Message = fmt.Sprintf("ERROR: %s", err.Error())
	w.UpdatedAt = time.Now()
}

func (w *Workspace) SetErrorWithMessage(message string) {
	w.Status = StatusError
	w.Message = fmt.Sprintf("ERROR: %s", message)
	w.UpdatedAt = time.Now()
}

func (w *Workspace) SetChecksum(checksum string) {
	w.Checksum = checksum
}

var (
	ErrWorkspaceNotFound   = errors.New("workspace not found")
	ErrInvalidDBURL        = errors.New("invalid database URL")
	ErrIngestionInProgress = errors.New("ingestion in progress")
	ErrDDLDMLDetected      = errors.New("DDL/DML operation detected")
	ErrSQLGenerationFailed = errors.New("failed to generate valid SQL after retries")
	ErrDatabaseUnreachable = errors.New("database unreachable")
	ErrNoMetadata          = errors.New("no metadata found for workspace")
)
