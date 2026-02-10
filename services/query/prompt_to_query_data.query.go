package query

import (
	"context"
	"errors"
	"time"

	"github.com/kamil5b/go-nl2query-lib/domains"
	"github.com/kamil5b/go-nl2query-lib/ports"
)

func (s *QueryService) PromptToQueryData(ctx context.Context, tenantID string, prompt string, withData bool) (*domains.Query, *string, error) {
	// Step 1: Check tenant status
	var warn *string
	workspaceStatus, _, err := s.statusAdapter.GetStatus(ctx, tenantID)
	if err != nil {
		return nil, nil, err
	}
	if workspaceStatus == domains.StatusInProgress {
		return nil, nil, ports.StatusInProgressError
	}

	// Step 2: Connect to internal database
	err = s.internalDatabaseAdapter.Connect(ctx, tenantID)
	if err != nil {
		return nil, nil, err
	}

	// Step 3: Get workspace information
	workspace, err := s.internalDatabaseAdapter.GetWorkspaceByTenantID(ctx, tenantID)
	if err != nil {
		return nil, nil, err
	}

	// Step 4: Decrypt and connect to client database if withData is true
	var clientDBConnected bool

	if withData && workspace != nil {
		decryptedURL, decErr := s.encryptAdapter.Decrypt(workspace.EncryptedDBURL)
		if decErr != nil {
			return nil, nil, decErr
		}
		if connErr := s.clientDatabaseAdapter.Connect(ctx, decryptedURL); connErr == nil {
			clientDBConnected = true
		} else {
			warnMsg := ports.QueryServiceWarnWontExecuteClientDatabaseError
			warn = &warnMsg
		}
	}

	// Step 5: Embed the prompt
	vector, err := s.embedderAdapter.Embed(ctx, prompt)
	if err != nil {
		return nil, nil, err
	}

	// Step 6: Search for relevant vectors
	vectors, err := s.vectorStoreAdapter.Search(ctx, tenantID, vector, 10)
	if err != nil {
		return nil, nil, err
	}

	// Step 7: Generate query with nested retry loops for syntax and execution errors
	var query *string
	additionalArgs := []string{}
	// Outer loop: for execution errors
	for executionIdx := 0; executionIdx < s.Config.ExecutionRetryLimit+1; executionIdx++ {

		// Inner loop: for syntax/validation errors
		for syntaxIdx := 0; syntaxIdx < s.Config.QueryFixAttempts+1; syntaxIdx++ {
			// Generate query
			var genErr error
			query, genErr = s.LLMAdapter.GenerateQuery(ctx, prompt, vectors, additionalArgs...)
			if genErr != nil {
				return nil, nil, genErr
			}

			// Validate safety
			isSafe, safeErr := s.queryValidatorAdapter.IsSafe(*query)
			if safeErr != nil && isSafe {
				break
			}
			if safeErr == nil {
				safeErr = errors.New("query deemed unsafe by validator")
			}
			additionalArgs = []string{*query, safeErr.Error()}
		}
		query, err := s.LLMAdapter.GenerateQuery(ctx, prompt, vectors, additionalArgs...)
		if err != nil {
			return nil, nil, err
		}
		// Validate safety
		isSafe, safeErr := s.queryValidatorAdapter.IsSafe(*query)
		if !isSafe && safeErr == nil {
			safeErr = errors.New("query deemed unsafe by validator")
		}
		if safeErr != nil || !isSafe || !withData || !clientDBConnected || query == nil {
			if safeErr != nil || !isSafe {
				warnMsg := ports.QueryServiceWarnQueryGeneratedUnsafe
				warn = &warnMsg
			}
			if !clientDBConnected {
				warnMsg := ports.QueryServiceWarnWontExecuteClientDatabaseError
				warn = &warnMsg
			}
			return &domains.Query{
				TenantID:    tenantID,
				ResultQuery: query,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, warn, nil
		}

		// Step 8: Check for DDL/DML and execute if applicable
		// Check if query contains DDL/DML
		if s.queryValidatorAdapter.ContainsDDLDML(*query) {
			// Contains DDL/DML, don't execute, return with warn
			warnMsg := ports.QueryServiceWarnDDLDMLDetected
			warn = &warnMsg
			return &domains.Query{
				TenantID:    tenantID,
				ResultQuery: query,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, warn, nil
		}

		// Execute the query
		result, execErr := s.clientDatabaseAdapter.Execute(ctx, *query)
		if execErr == nil {
			// Execution successful, return with data
			return &domains.Query{
				TenantID:    tenantID,
				ResultQuery: query,
				ResultData:  result,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			}, warn, nil
		}

		// Execution failed, prepare error for next outer loop iteration
		additionalArgs = []string{*query, execErr.Error()}
		// Continue outer loop to retry query generation with execution error

	}
	warnMsg := ports.QueryServiceWarnQueryGeneratedUnsafe
	warn = &warnMsg

	// Return success with the generated query
	return &domains.Query{
		TenantID:    tenantID,
		ResultQuery: query,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, warn, nil
}
