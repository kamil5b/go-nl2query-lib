package workspace

import (
	"github.com/kamil5b/go-nl2query-lib/ports"
)

type WorkspaceConfig struct{}

type WorkspaceService struct {
	Config                  *WorkspaceConfig
	statusAdapter           ports.StatusPort
	clientDatabaseAdapter   ports.ClientDatabasePort
	internalDatabaseAdapter ports.InternalDatabasePort
	encryptAdapter          ports.EncryptPort
	hashAdapter             ports.HashPort
	taskQueueService        ports.TaskQueuePort
}

func NewWorkspaceService(
	config *WorkspaceConfig,
	statusAdapter ports.StatusPort,
	clientDatabaseAdapter ports.ClientDatabasePort,
	internalDatabaseAdapter ports.InternalDatabasePort,
	encryptAdapter ports.EncryptPort,
	hashAdapter ports.HashPort,
	taskQueueService ports.TaskQueuePort,
) *WorkspaceService {
	return &WorkspaceService{
		Config:                  config,
		statusAdapter:           statusAdapter,
		clientDatabaseAdapter:   clientDatabaseAdapter,
		internalDatabaseAdapter: internalDatabaseAdapter,
		encryptAdapter:          encryptAdapter,
		hashAdapter:             hashAdapter,
		taskQueueService:        taskQueueService,
	}
}
