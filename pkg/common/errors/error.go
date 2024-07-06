// Package errors provides custom error types and error checking functions
// for common error scenarios in the DataVinci project.
package errors

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"strings"
)

// ErrorType represents the type of error.
type ErrorType uint

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeDatabaseConnection
	ErrorTypeTimeout
	ErrorTypePermission
	ErrorTypeQuery
	ErrorTypeExecution
	ErrorTypeTransaction
	ErrorTypeConfiguration
	ErrorTypeAPIConnection
	ErrorTypeUnsupported
	ErrorTypeFileConnection
	ErrorTypeNotFound
	ErrorTypeConnection
	ErrorTypeTransformation
	ErrorTypeEmptyPassword
	ErrorTypeInvalidCost
	ErrorTypeValidation
	ErrorTypeResourceExhausted
	ErrorTypeDataIntegrity
)

// ErrorMessages maps error types to human-readable messages.
var ErrorMessages = map[ErrorType]string{
	ErrorTypeUnknown:            "Unknown error occurred",
	ErrorTypeDatabaseConnection: "Database connection not established",
	ErrorTypeTimeout:            "Operation timed out",
	ErrorTypePermission:         "Permission denied",
	ErrorTypeQuery:              "Query execution failed",
	ErrorTypeExecution:          "Execution error",
	ErrorTypeTransaction:        "Transaction error",
	ErrorTypeConfiguration:      "Configuration error",
	ErrorTypeAPIConnection:      "API connection failed",
	ErrorTypeUnsupported:        "Unsupported operation",
	ErrorTypeFileConnection:     "File connection error",
	ErrorTypeNotFound:           "Resource not found",
	ErrorTypeConnection:         "Connection error",
	ErrorTypeTransformation:     "Data transformation error",
	ErrorTypeEmptyPassword:      "Empty password provided",
	ErrorTypeInvalidCost:        "Invalid bcrypt cost",
	ErrorTypeValidation:         "Validation error",
	ErrorTypeResourceExhausted:  "Resource exhausted",
	ErrorTypeDataIntegrity:      "Data integrity violation",
}

// GetErrorMessage returns the error message for a given ErrorType.
// If the ErrorType is not found in ErrorMessages, it returns a default message.
func GetErrorMessage(errType ErrorType) string {
	if msg, ok := ErrorMessages[errType]; ok {
		return msg
	}
	return fmt.Sprintf("Unknown error type (%d)", errType)
}

// Error represents a custom error with additional context.
type Error struct {
	Type    ErrorType         // The type of the error
	Message string            // A human-readable error message
	Err     error             // The underlying error, if any
	Cause   error             // The cause of the error, if any
	Stack   string            // The stack trace of the error
	Context map[string]string // Additional context for the error
}

// Error returns the error message.
// It implements the error interface.
func (e *Error) Error() string {
	message := e.Message
	if message == "" {
		message = GetErrorMessage(e.Type)
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", message, e.Err)
	}
	return message
}

// Unwrap returns the wrapped error.
// It allows Error to work with errors.Is and errors.As.
func (e *Error) Unwrap() error {
	return e.Err
}

// Is checks if the given error is of the given type.
// It allows Error to work with errors.Is and errors.As.
func (e *Error) Is(target error) bool {
	if target == nil {
		return false
	}
	t, ok := target.(*Error)
	if !ok {
		return errors.Is(e.Err, target)
	}
	return e.Type == t.Type
}

// NewError creates a new Error with the given type, message, and underlying error.
func NewError(errType ErrorType, message string, err error) *Error {
	return NewErrorWithContext(errType, message, err, nil)
}

// NewErrorWithContext creates a new Error with the given type, message, underlying error, and context.
func NewErrorWithContext(errType ErrorType, message string, err error, context map[string]string) *Error {
	if message == "" {
		message = GetErrorMessage(errType)
	}
	newError := &Error{
		Type:    errType,
		Message: message,
		Err:     err,
		Cause:   err,
		Context: context,
	}

	var sb strings.Builder
	for i := 1; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		sb.WriteString(fmt.Sprintf("%s:%d %s\n", file, line, fn.Name()))
	}
	newError.Stack = sb.String()

	return newError
}

// IsErrorType checks if the error is of the given type.
func IsErrorType(err error, errType ErrorType) bool {
	var customErr *Error
	if errors.As(err, &customErr) {
		return customErr.Type == errType
	}
	return false
}

// IsConnectionError checks if the given error is a connection error.
func IsConnectionError(err error) bool {
	if err == nil {
		return false
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypeDatabaseConnection || e.Type == ErrorTypeConnection || e.Type == ErrorTypeAPIConnection
	}
	var netErr *net.OpError
	var dnsErr *net.DNSError
	return errors.As(err, &netErr) || errors.As(err, &dnsErr) || isNetworkError(err)
}

// IsTimeoutError checks if the given error is a timeout error.
func IsTimeoutError(err error) bool {
	if err == nil {
		return false
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypeTimeout
	}
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return netErr.Timeout()
	}
	return errors.Is(err, context.DeadlineExceeded) ||
		strings.Contains(strings.ToLower(err.Error()), "timeout") ||
		strings.Contains(strings.ToLower(err.Error()), "deadline exceeded")
}

// IsPermissionError checks if the given error is a permission error.
func IsPermissionError(err error) bool {
	if err == nil {
		return false
	}
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypePermission
	}
	errLower := strings.ToLower(err.Error())
	return strings.Contains(errLower, "permission") ||
		strings.Contains(errLower, "access denied")
}

// isNetworkError checks if the error is a known network error.
func isNetworkError(err error) bool {
	return errors.Is(err, net.ErrClosed) || errors.Is(err, net.ErrWriteToConnected)
}

// AddErrorContext adds additional context to an existing Error.
func AddErrorContext(err *Error, key, value string) *Error {
	if err.Context == nil {
		err.Context = make(map[string]string)
	}
	err.Context[key] = value
	return err
}

// GetErrorContext retrieves the context from an Error.
func GetErrorContext(err error) map[string]string {
	var e *Error
	if errors.As(err, &e) {
		return e.Context
	}
	return nil
}

// WrapError wraps an existing error with a new Error type and message.
func WrapError(err error, errType ErrorType, message string) *Error {
	return NewError(errType, message, err)
}
