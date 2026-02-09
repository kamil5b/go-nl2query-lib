package ports

import (
	"context"

	model "github.com/kamil5b/go-nl2query-lib/domains"
)

const (
	QueryServiceWarnDDLDMLDetected                 = "DDL or DML statement detected. Query won't be executed."
	QueryServiceWarnUseExistingClientDatabaseError = "Will using existing stored schema because connection to client database could not be established."
	QueryServiceWarnWontExecuteClientDatabaseError = "Query won't be executed because connection to client database could not be established."
)

type QueryService interface {
	PromptToQueryData(ctx context.Context, tenantID string, prompt string, withData bool) (result *model.Query, msg *string, err error)
}
