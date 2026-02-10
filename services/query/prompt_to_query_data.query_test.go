package query

import (
	"testing"

	"github.com/kamil5b/go-nl2query-lib/ports"
	queryTest "github.com/kamil5b/go-nl2query-lib/testsuites/query"
)

func TestQueryService_PromptToQueryData(t *testing.T) {
	queryTest.UnitTestPromptToQueryData(t, func(
		statusAdapter ports.StatusPort,
		clientDatabaseAdapter ports.ClientDatabasePort,
		internalDatabaseAdapter ports.InternalDatabasePort,
		encryptAdapter ports.EncryptPort,
		embedderAdapter ports.EmbedderPort,
		vectorStoreAdapter ports.VectorStorePort,
		LLMAdapter ports.LLMPort,
		QueryValidator ports.QueryValidatorPort,
		queryErrorLimit int,
		executionErrorLimit int,
	) ports.QueryService {
		return NewQueryService(
			&QueryConfig{
				ExecutionRetryLimit: executionErrorLimit,
				QueryFixAttempts:    queryErrorLimit,
			},
			statusAdapter,
			clientDatabaseAdapter,
			internalDatabaseAdapter,
			encryptAdapter,
			embedderAdapter,
			vectorStoreAdapter,
			LLMAdapter,
			QueryValidator,
		)
	})
}
