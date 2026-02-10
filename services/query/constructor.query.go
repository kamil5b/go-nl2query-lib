package query

import "github.com/kamil5b/go-nl2query-lib/ports"

type QueryConfig struct {
	ExecutionRetryLimit int
	QueryFixAttempts    int
}

type QueryService struct {
	Config *QueryConfig

	statusAdapter           ports.StatusPort
	clientDatabaseAdapter   ports.ClientDatabasePort
	internalDatabaseAdapter ports.InternalDatabasePort
	encryptAdapter          ports.EncryptPort
	embedderAdapter         ports.EmbedderPort
	vectorStoreAdapter      ports.VectorStorePort
	LLMAdapter              ports.LLMPort
	queryValidatorAdapter   ports.QueryValidatorPort
}

func NewQueryService(
	config *QueryConfig,

	statusAdapter ports.StatusPort,
	clientDatabaseAdapter ports.ClientDatabasePort,
	internalDatabaseAdapter ports.InternalDatabasePort,
	encryptAdapter ports.EncryptPort,
	embedderAdapter ports.EmbedderPort,
	vectorStoreAdapter ports.VectorStorePort,
	LLMAdapter ports.LLMPort,
	queryValidatorAdapter ports.QueryValidatorPort,
) *QueryService {
	return &QueryService{
		Config: config,

		embedderAdapter:         embedderAdapter,
		vectorStoreAdapter:      vectorStoreAdapter,
		statusAdapter:           statusAdapter,
		clientDatabaseAdapter:   clientDatabaseAdapter,
		internalDatabaseAdapter: internalDatabaseAdapter,
		encryptAdapter:          encryptAdapter,
		LLMAdapter:              LLMAdapter,
		queryValidatorAdapter:   queryValidatorAdapter,
	}
}
