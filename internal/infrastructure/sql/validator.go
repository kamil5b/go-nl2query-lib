package sql

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrSQLGenerationFailed = errors.New("failed to generate valid SQL after retries")
)

type Validator struct {
	ddlDMLPattern *regexp.Regexp
}

func NewValidator() *Validator {
	ddlDMLPattern := regexp.MustCompile(`(?i)\b(INSERT|UPDATE|DELETE|DROP|CREATE|ALTER|TRUNCATE|REPLACE|GRANT|REVOKE)\b`)
	return &Validator{
		ddlDMLPattern: ddlDMLPattern,
	}
}

func (v *Validator) IsSafe(sql string) (bool, error) {
	return !v.ContainsDDLDML(sql), nil
}

func (v *Validator) ContainsDDLDML(sql string) bool {
	sql = removeComments(sql)
	sql = removeStringLiterals(sql)
	return v.ddlDMLPattern.MatchString(sql)
}

func (v *Validator) ValidateSyntax(sql string) error {
	if sql == "" {
		return ErrSQLGenerationFailed
	}

	normalized := strings.ToUpper(strings.TrimSpace(sql))
	if !strings.HasPrefix(normalized, "SELECT") && !strings.HasPrefix(normalized, "WITH") {
		return ErrSQLGenerationFailed
	}

	if !isBalanced(sql) {
		return ErrSQLGenerationFailed
	}

	return nil
}

func removeComments(sql string) string {
	sql = regexp.MustCompile(`--[^\n]*`).ReplaceAllString(sql, "")
	sql = regexp.MustCompile(`/\*[\s\S]*?\*/`).ReplaceAllString(sql, "")
	return sql
}

func removeStringLiterals(sql string) string {
	sql = regexp.MustCompile(`'(?:''|[^'])*'`).ReplaceAllString(sql, "''")
	sql = regexp.MustCompile(`"(?:""|[^"])*"`).ReplaceAllString(sql, `""`)
	return sql
}

func isBalanced(sql string) bool {
	count := 0
	inString := false
	stringChar := rune(0)

	for _, char := range sql {
		if !inString {
			if char == '\'' || char == '"' {
				inString = true
				stringChar = char
			} else if char == '(' {
				count++
			} else if char == ')' {
				count--
			}
		} else {
			if char == stringChar {
				inString = false
			}
		}

		if count < 0 {
			return false
		}
	}

	return count == 0 && !inString
}
