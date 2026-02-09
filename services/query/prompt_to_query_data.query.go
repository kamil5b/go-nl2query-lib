package query

import (
	"context"
	"time"

	"github.com/kamil5b/go-nl2query-lib/domains"
	"github.com/kamil5b/go-nl2query-lib/ports"
)

func (s *QueryService) PromptToQueryData(ctx context.Context, tenantID string, prompt string, withData bool) (*domains.Query, error) {
	// Step 1: Check tenant status
	workspaceStatus, _, err := s.statusAdapter.GetStatus(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	if workspaceStatus == domains.StatusInProgress {
		return nil, ports.StatusInProgressError
	}

	// Step 2: Connect to internal database
	err = s.internalDatabaseAdapter.Connect(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Step 3: Get workspace information
	workspace, err := s.internalDatabaseAdapter.GetWorkspaceByTenantID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	// Step 4: Decrypt and connect to client database if withData is true
	var clientDBConnected bool

	if withData && workspace != nil {
		decryptedURL, decErr := s.encryptAdapter.Decrypt(workspace.EncryptedDBURL)
		if decErr == nil {
			// Try to connect to client database
			if connErr := s.clientDatabaseAdapter.Connect(ctx, decryptedURL); connErr == nil {
				clientDBConnected = true
			}
		}
	}

	// Step 5: Embed the prompt
	vector, err := s.embedderAdapter.Embed(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Step 6: Search for relevant vectors
	vectors, err := s.vectorStoreAdapter.Search(ctx, tenantID, vector, 10)
	if err != nil {
		return nil, err
	}

	// Step 7: Generate query in a loop until safe
	var query *string
	var lastError string

	for {
		var generationErr error
		if lastError == "" {
			query, generationErr = s.LLMAdapter.GenerateQuery(ctx, prompt, vectors)
		} else {
			query, generationErr = s.LLMAdapter.GenerateQuery(ctx, prompt, vectors, lastError)
		}

		if generationErr != nil {
			return nil, generationErr
		}

		// Validate safety
		isSafe, safeErr := s.queryValidatorAdapter.IsSafe(*query)
		if safeErr != nil || !isSafe {
			// Query is not safe, regenerate with error message
			if safeErr != nil {
				lastError = safeErr.Error()
			} else {
				lastError = "query validation failed"
			}
			continue
		}

		// Query is safe, use it
		break
	}

	// Step 8: Execute query if withData is true and client DB is connected
	if withData && clientDBConnected {
		for {
			// Check if query contains DDL/DML
			if s.queryValidatorAdapter.ContainsDDLDML(*query) {
				// Contains DDL/DML, don't execute, just return with query
				break
			}

			// Execute the query
			result, execErr := s.clientDatabaseAdapter.Execute(ctx, *query)
			if execErr == nil {
				// Execution successful, return with data
				return &domains.Query{
					TenantID:    tenantID,
					ResultQuery: query,
					ResultData:  convertToMapPtr(result),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}, nil
			}

			// Execution failed, try to regenerate query with error
			execErrMsg := execErr.Error()
			newQuery, genErr := s.LLMAdapter.GenerateQuery(ctx, prompt, vectors, execErrMsg)
			if genErr != nil {
				return nil, genErr
			}

			// Validate the new query
			isSafe, safeErr := s.queryValidatorAdapter.IsSafe(*newQuery)
			if safeErr != nil || !isSafe {
				// New query validation failed, regenerate again
				continue
			}

			// New query is safe, use it for next execution attempt
			query = newQuery
		}
	}

	// Return success with the generated query
	return &domains.Query{
		TenantID:    tenantID,
		ResultQuery: query,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

// convertToMapPtr converts []map[string]any to *map[string]any by taking the first element
func convertToMapPtr(data []map[string]any) *map[string]any {
	if len(data) == 0 {
		return nil
	}
	return &data[0]
}
