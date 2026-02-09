package domains

import (
	"strings"
	"testing"
)

func TestGoNL2QueryError_Error(t *testing.T) {
	tests := []struct {
		name     string
		error    GoNL2QueryError
		expected string
	}{
		{
			name: "error with message only",
			error: GoNL2QueryError{
				StatusCode: 400,
				Message:    "Invalid query",
			},
			expected: "Invalid query",
		},
		{
			name: "error with message and single additional info",
			error: GoNL2QueryError{
				StatusCode:          400,
				Message:             "Invalid query",
				AdditionalErrorInfo: []string{"Field 'name' is required"},
			},
			expected: "Invalid query: Field 'name' is required",
		},
		{
			name: "error with message and multiple additional info",
			error: GoNL2QueryError{
				StatusCode:          400,
				Message:             "Validation failed",
				AdditionalErrorInfo: []string{"Field 'name' is required", "Field 'email' must be valid", "Field 'age' must be positive"},
			},
			expected: "Validation failed: Field 'name' is required; Field 'email' must be valid; Field 'age' must be positive",
		},
		{
			name: "error with empty additional info slice",
			error: GoNL2QueryError{
				StatusCode:          500,
				Message:             "Internal server error",
				AdditionalErrorInfo: []string{},
			},
			expected: "Internal server error",
		},
		{
			name: "error with nil additional info",
			error: GoNL2QueryError{
				StatusCode:          500,
				Message:             "Internal server error",
				AdditionalErrorInfo: nil,
			},
			expected: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.error.Error()
			if got != tt.expected {
				t.Errorf("Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGoNL2QueryError_AddAdditionalErrorInfo(t *testing.T) {
	tests := []struct {
		name     string
		initial  GoNL2QueryError
		infos    []string
		expected []string
	}{
		{
			name: "add single info to empty error",
			initial: GoNL2QueryError{
				StatusCode: 400,
				Message:    "Test error",
			},
			infos:    []string{"Info 1"},
			expected: []string{"Info 1"},
		},
		{
			name: "add single info to error with existing info",
			initial: GoNL2QueryError{
				StatusCode:          400,
				Message:             "Test error",
				AdditionalErrorInfo: []string{"Existing info"},
			},
			infos:    []string{"New info"},
			expected: []string{"Existing info", "New info"},
		},
		{
			name: "add multiple infos sequentially",
			initial: GoNL2QueryError{
				StatusCode:          400,
				Message:             "Test error",
				AdditionalErrorInfo: []string{},
			},
			infos:    []string{"Info 1", "Info 2", "Info 3"},
			expected: []string{"Info 1", "Info 2", "Info 3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errPtr := &tt.initial
			for _, info := range tt.infos {
				errPtr = errPtr.AddAdditionalErrorInfo(info)
			}

			if len(errPtr.AdditionalErrorInfo) != len(tt.expected) {
				t.Errorf("AdditionalErrorInfo length = %d, want %d", len(errPtr.AdditionalErrorInfo), len(tt.expected))
			}

			for i, info := range errPtr.AdditionalErrorInfo {
				if info != tt.expected[i] {
					t.Errorf("AdditionalErrorInfo[%d] = %q, want %q", i, info, tt.expected[i])
				}
			}
		})
	}
}

func TestGoNL2QueryError_AddAdditionalErrorInfo_ReturnValue(t *testing.T) {
	err := &GoNL2QueryError{
		StatusCode: 400,
		Message:    "Test error",
	}

	result := err.AddAdditionalErrorInfo("Info 1")

	if result != err {
		t.Error("AddAdditionalErrorInfo should return the same pointer")
	}
}

func TestGoNL2QueryError_AddBatchAdditionalErrorInfo(t *testing.T) {
	tests := []struct {
		name     string
		initial  GoNL2QueryError
		infos    []string
		expected []string
	}{
		{
			name: "add empty batch to error",
			initial: GoNL2QueryError{
				StatusCode:          400,
				Message:             "Test error",
				AdditionalErrorInfo: []string{"Existing"},
			},
			infos:    []string{},
			expected: []string{"Existing"},
		},
		{
			name: "add batch to empty error",
			initial: GoNL2QueryError{
				StatusCode: 400,
				Message:    "Test error",
			},
			infos:    []string{"Info 1", "Info 2", "Info 3"},
			expected: []string{"Info 1", "Info 2", "Info 3"},
		},
		{
			name: "add batch to error with existing info",
			initial: GoNL2QueryError{
				StatusCode:          400,
				Message:             "Test error",
				AdditionalErrorInfo: []string{"Existing 1", "Existing 2"},
			},
			infos:    []string{"New 1", "New 2", "New 3"},
			expected: []string{"Existing 1", "Existing 2", "New 1", "New 2", "New 3"},
		},
		{
			name: "add single item batch",
			initial: GoNL2QueryError{
				StatusCode: 400,
				Message:    "Test error",
			},
			infos:    []string{"Single info"},
			expected: []string{"Single info"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errPtr := &tt.initial
			errPtr = errPtr.AddBatchAdditionalErrorInfo(tt.infos)

			if len(errPtr.AdditionalErrorInfo) != len(tt.expected) {
				t.Errorf("AdditionalErrorInfo length = %d, want %d", len(errPtr.AdditionalErrorInfo), len(tt.expected))
			}

			for i, info := range errPtr.AdditionalErrorInfo {
				if info != tt.expected[i] {
					t.Errorf("AdditionalErrorInfo[%d] = %q, want %q", i, info, tt.expected[i])
				}
			}
		})
	}
}

func TestGoNL2QueryError_AddBatchAdditionalErrorInfo_ReturnValue(t *testing.T) {
	err := &GoNL2QueryError{
		StatusCode: 400,
		Message:    "Test error",
	}

	result := err.AddBatchAdditionalErrorInfo([]string{"Info 1", "Info 2"})

	if result != err {
		t.Error("AddBatchAdditionalErrorInfo should return the same pointer")
	}
}

func TestGoNL2QueryError_MethodChaining(t *testing.T) {
	t.Run("chain AddAdditionalErrorInfo calls", func(t *testing.T) {
		err := &GoNL2QueryError{
			StatusCode: 400,
			Message:    "Test error",
		}

		result := err.
			AddAdditionalErrorInfo("Error 1").
			AddAdditionalErrorInfo("Error 2").
			AddAdditionalErrorInfo("Error 3")

		expectedInfo := []string{"Error 1", "Error 2", "Error 3"}
		if len(result.AdditionalErrorInfo) != len(expectedInfo) {
			t.Errorf("AdditionalErrorInfo length = %d, want %d", len(result.AdditionalErrorInfo), len(expectedInfo))
		}

		for i, info := range result.AdditionalErrorInfo {
			if info != expectedInfo[i] {
				t.Errorf("AdditionalErrorInfo[%d] = %q, want %q", i, info, expectedInfo[i])
			}
		}
	})

	t.Run("chain mixed Add methods", func(t *testing.T) {
		err := &GoNL2QueryError{
			StatusCode: 400,
			Message:    "Test error",
		}

		result := err.
			AddAdditionalErrorInfo("Error 1").
			AddBatchAdditionalErrorInfo([]string{"Error 2", "Error 3"}).
			AddAdditionalErrorInfo("Error 4")

		expectedInfo := []string{"Error 1", "Error 2", "Error 3", "Error 4"}
		if len(result.AdditionalErrorInfo) != len(expectedInfo) {
			t.Errorf("AdditionalErrorInfo length = %d, want %d", len(result.AdditionalErrorInfo), len(expectedInfo))
		}

		for i, info := range result.AdditionalErrorInfo {
			if info != expectedInfo[i] {
				t.Errorf("AdditionalErrorInfo[%d] = %q, want %q", i, info, expectedInfo[i])
			}
		}
	})
}

func TestGoNL2QueryError_ErrorInterface(t *testing.T) {
	t.Run("implements error interface", func(t *testing.T) {
		err := &GoNL2QueryError{
			StatusCode: 400,
			Message:    "Test error",
		}

		var e error = err
		if e == nil {
			t.Error("GoNL2QueryError should implement error interface")
		}

		if e.Error() != "Test error" {
			t.Errorf("Error() = %q, want %q", e.Error(), "Test error")
		}
	})
}

func TestGoNL2QueryError_EdgeCases(t *testing.T) {
	t.Run("error message with special characters", func(t *testing.T) {
		err := GoNL2QueryError{
			StatusCode:          400,
			Message:             "Error: Invalid input",
			AdditionalErrorInfo: []string{"Details: [1, 2, 3]", "Code: 'ABC-123'"},
		}

		result := err.Error()
		expectedSubstrings := []string{"Error: Invalid input", "Details: [1, 2, 3]", "Code: 'ABC-123'"}
		for _, substr := range expectedSubstrings {
			if !strings.Contains(result, substr) {
				t.Errorf("Error() should contain %q, got %q", substr, result)
			}
		}
	})

	t.Run("error with very long additional info", func(t *testing.T) {
		longInfo := strings.Repeat("x", 1000)
		err := GoNL2QueryError{
			StatusCode:          500,
			Message:             "Error",
			AdditionalErrorInfo: []string{longInfo},
		}

		result := err.Error()
		if !strings.Contains(result, longInfo) {
			t.Error("Error() should contain the long additional info")
		}
	})

	t.Run("add info to nil AdditionalErrorInfo", func(t *testing.T) {
		err := &GoNL2QueryError{
			StatusCode: 400,
			Message:    "Test",
		}

		result := err.AddAdditionalErrorInfo("Info")
		if len(result.AdditionalErrorInfo) != 1 {
			t.Errorf("AdditionalErrorInfo length = %d, want 1", len(result.AdditionalErrorInfo))
		}
	})
}

func TestGoNL2QueryError_StatusCode(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"Bad Request", 400},
		{"Unauthorized", 401},
		{"Forbidden", 403},
		{"Not Found", 404},
		{"Internal Server Error", 500},
		{"Service Unavailable", 503},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GoNL2QueryError{
				StatusCode: tt.statusCode,
				Message:    "Test",
			}

			if err.StatusCode != tt.statusCode {
				t.Errorf("StatusCode = %d, want %d", err.StatusCode, tt.statusCode)
			}
		})
	}
}
