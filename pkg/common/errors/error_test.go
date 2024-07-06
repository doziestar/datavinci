package errors_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	. "pkg/common/errors"
)

func TestIsConnectionError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Connection error", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), true},
		{"Network error", net.ErrClosed, true},
		{"DNS error", &net.DNSError{Err: "no such host", Name: "example.com"}, true},
		{"Dial error", &net.OpError{Op: "dial", Err: fmt.Errorf("connection refused")}, true},
		{"Timeout error", NewError(ErrorTypeTimeout, "operation timed out", nil), false},
		{"Permission error", NewError(ErrorTypePermission, "access denied", nil), false},
		{"Other error", errors.New("some error"), false},
		{"Nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsConnectionError(tt.err)
			assert.Equal(t, tt.want, got, "IsConnectionError() = %v, want %v", got, tt.want)
		})
	}
}

func TestIsTimeoutError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Timeout error", NewError(ErrorTypeTimeout, "operation timed out", nil), true},
		{"Network timeout", &net.OpError{Op: "read", Err: os.ErrDeadlineExceeded}, true},
		{"Context deadline exceeded", context.DeadlineExceeded, true},
		{"I/O timeout string", errors.New("i/o timeout"), true},
		{"Connection error", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), false},
		{"Permission error", NewError(ErrorTypePermission, "access denied", nil), false},
		{"Other error", errors.New("some error"), false},
		{"Nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsTimeoutError(tt.err)
			assert.Equal(t, tt.want, got, "IsTimeoutError() = %v, want %v", got, tt.want)
		})
	}
}

func TestIsPermissionError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"Permission error", NewError(ErrorTypePermission, "access denied", nil), true},
		{"Permission string", errors.New("permission denied"), true},
		{"OS permission denied", os.ErrPermission, true},
		{"Access denied string", errors.New("access denied"), true},
		{"Connection error", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), false},
		{"Timeout error", NewError(ErrorTypeTimeout, "operation timed out", nil), false},
		{"Other error", errors.New("some error"), false},
		{"Nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsPermissionError(tt.err)
			assert.Equal(t, tt.want, got, "IsPermissionError() = %v, want %v", got, tt.want)
		})
	}
}

func TestErrorIs(t *testing.T) {
	baseErr := errors.New("base error")
	tests := []struct {
		name   string
		err    *Error
		target error
		want   bool
	}{
		{"Equal error", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), NewError(ErrorTypeDatabaseConnection, "connection failed", nil), true},
		{"Different error", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), NewError(ErrorTypePermission, "access denied", nil), false},
		{"Same type, different message", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), NewError(ErrorTypeDatabaseConnection, "connection error", nil), true},
		{"With wrapped error", NewError(ErrorTypeDatabaseConnection, "connection failed", baseErr), baseErr, true},
		{"Different wrapped error", NewError(ErrorTypeDatabaseConnection, "connection failed", baseErr), errors.New("other error"), false},
		{"Nil target", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errors.Is(tt.err, tt.target)
			assert.Equal(t, tt.want, got, "errors.Is() = %v, want %v", got, tt.want)
		})
	}
}

func TestErrorUnwrap(t *testing.T) {
	baseErr := errors.New("base error")
	tests := []struct {
		name string
		err  *Error
		want error
	}{
		{"Nil wrapped error", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), nil},
		{"With wrapped error", NewError(ErrorTypeDatabaseConnection, "connection failed", baseErr), baseErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Unwrap()
			assert.Equal(t, tt.want, got, "Error.Unwrap() = %v, want %v", got, tt.want)
		})
	}
}

func TestErrorError(t *testing.T) {
	tests := []struct {
		name    string
		err     *Error
		want    string
		wantErr bool
	}{
		{"Simple error", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), "connection failed", false},
		{"Error with type, no custom message", NewError(ErrorTypePermission, "", nil), "Permission denied", false},
		{"Error with wrapped error", NewError(ErrorTypeTimeout, "operation timed out", errors.New("underlying error")), "operation timed out: underlying error", false},
		{"Error with type, no custom message, with wrapped error", NewError(ErrorTypeTimeout, "", errors.New("underlying error")), "Operation timed out: underlying error", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			assert.Equal(t, tt.want, got, "Error.Error() = %v, want %v", got, tt.want)
		})
	}
}

func TestNewError(t *testing.T) {
	baseErr := errors.New("base error")
	tests := []struct {
		name    string
		errType ErrorType
		message string
		cause   error
	}{
		{"Connection error with custom message", ErrorTypeDatabaseConnection, "custom connection failed", nil},
		{"Timeout error with cause, no custom message", ErrorTypeTimeout, "", baseErr},
		{"Permission error with default message", ErrorTypePermission, "", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewError(tt.errType, tt.message, tt.cause)
			require.NotNil(t, err)
			assert.Equal(t, tt.errType, err.Type)
			if tt.message != "" {
				assert.Equal(t, tt.message, err.Message)
			} else {
				assert.Equal(t, GetErrorMessage(tt.errType), err.Message)
			}
			assert.Equal(t, tt.cause, err.Cause)
		})
	}
}

func TestErrorWithStack(t *testing.T) {
	err := NewError(ErrorTypeDatabaseConnection, "connection failed", nil)
	assert.NotEmpty(t, err.Stack, "Error stack should not be empty")
}

func TestErrorAs(t *testing.T) {
	var target *Error
	err := NewError(ErrorTypeDatabaseConnection, "connection failed", nil)

	assert.True(t, errors.As(err, &target))
	assert.Equal(t, err, target)

	var notTarget *net.OpError
	assert.False(t, errors.As(err, &notTarget))
}

func TestIsErrorType(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		errType  ErrorType
		expected bool
	}{
		{"Matching type", NewError(ErrorTypeDatabaseConnection, "connection failed", nil), ErrorTypeDatabaseConnection, true},
		{"Non-matching type", NewError(ErrorTypeTimeout, "operation timed out", nil), ErrorTypeDatabaseConnection, false},
		{"Non-custom error", errors.New("some error"), ErrorTypeDatabaseConnection, false},
		{"Nil error", nil, ErrorTypeDatabaseConnection, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsErrorType(tt.err, tt.errType)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetErrorContext(t *testing.T) {
	context := map[string]string{"server": "db01", "retry": "3"}
	err := NewErrorWithContext(ErrorTypeDatabaseConnection, "connection failed", nil, context)

	tests := []struct {
		name string
		err  error
		want map[string]string
	}{
		{"Get context from Error", err, context},
		{"Get context from non-Error", errors.New("some error"), nil},
		{"Get context from nil", nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetErrorContext(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWrapError(t *testing.T) {
	baseErr := errors.New("base error")

	tests := []struct {
		name    string
		err     error
		errType ErrorType
		message string
	}{
		{"Wrap base error", baseErr, ErrorTypeDatabaseConnection, "connection failed"},
		{"Wrap nil error", nil, ErrorTypeTimeout, "operation timed out"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wrapped := WrapError(tt.err, tt.errType, tt.message)
			assert.Equal(t, tt.errType, wrapped.Type)
			assert.Equal(t, tt.message, wrapped.Message)
			assert.Equal(t, tt.err, wrapped.Cause)
		})
	}
}

func TestErrorTypeString(t *testing.T) {
	tests := []struct {
		name     string
		errType  ErrorType
		expected string
	}{
		{"Unknown error", ErrorTypeUnknown, "Unknown error occurred"},
		{"Database connection error", ErrorTypeDatabaseConnection, "Database connection not established"},
		{"Timeout error", ErrorTypeTimeout, "Operation timed out"},
		{"Non-existent error type", ErrorType(9999), "Unknown error type (9999)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, GetErrorMessage(tt.errType))
		})
	}
}
func TestErrorStack(t *testing.T) {
	err := NewError(ErrorTypeDatabaseConnection, "connection failed", nil)

	assert.NotEmpty(t, err.Stack)
	assert.Contains(t, err.Stack, "errors_test.TestErrorStack")
	assert.Contains(t, err.Stack, "error_test.go")
}

func TestErrorString(t *testing.T) {
	tests := []struct {
		name     string
		err      *Error
		expected string
	}{
		{
			name:     "Error without underlying error",
			err:      NewError(ErrorTypeTimeout, "", nil),
			expected: "Operation timed out",
		},
		{
			name:     "Error with underlying error",
			err:      NewError(ErrorTypeDatabaseConnection, "", fmt.Errorf("connection refused")),
			expected: "Database connection not established: connection refused",
		},
		{
			name:     "Unknown error type",
			err:      NewError(ErrorType(9999), "", nil),
			expected: "Unknown error type (9999)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}
