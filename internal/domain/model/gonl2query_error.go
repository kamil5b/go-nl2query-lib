package model

import "strings"

type GoNL2QueryError struct {
	StatusCode          int      `json:"statusCode"`
	Message             string   `json:"message"`
	AdditionalErrorInfo []string `json:"additionalErrorInfo,omitempty"`
}

func (g GoNL2QueryError) Error() string {
	if len(g.AdditionalErrorInfo) > 0 {
		additionalInfo := ": " + strings.Join(g.AdditionalErrorInfo, "; ")
		return g.Message + additionalInfo
	}
	return g.Message
}

func (g *GoNL2QueryError) AddAdditionalErrorInfo(info string) *GoNL2QueryError {
	g.AdditionalErrorInfo = append(g.AdditionalErrorInfo, info)
	return g
}

func (g *GoNL2QueryError) AddBatchAdditionalErrorInfo(info []string) *GoNL2QueryError {
	g.AdditionalErrorInfo = append(g.AdditionalErrorInfo, info...)
	return g
}
