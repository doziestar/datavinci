// Package errors provide custom error types and error checking functions
// for common error scenarios in the DataVinci project.
package errors

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// ErrorType represents the type of error.
type ErrorType uint

// ErrorMessages maps error types to human-readable messages.
var ErrorMessages = map[ErrorType]string{
	ErrorTypeUnknown:    "unknown",
	ErrorTypeConnection: "connection",
	ErrorTypeTimeout:    "timeout",
	ErrorTypePermission: "permission",
}

const (
	// ErrorTypeUnknown represents an unknown error.
	ErrorTypeUnknown ErrorType = iota
	// ErrorTypeConnection represents a connection error.
	ErrorTypeConnection
	// ErrorTypeTimeout represents a timeout error.
	ErrorTypeTimeout
	// ErrorTypePermission represents a permission error.
	ErrorTypePermission
)

// Error represents a custom error with additional context.
type Error struct {
	Type    ErrorType // The type of the error
	Message string    // A human-readable error message
	Err     error     // The underlying error, if any
}

// Error returns the error message.
// It implements the error interface.
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error.
// It allows Error to work with errors.Is and errors.As.
func (e *Error) Unwrap() error {
	return e.Err
}

// NewError creates a new Error with the given type, message, and underlying error.
func NewError(errType ErrorType, message string, err error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// IsConnectionError checks if the given error is a connection error.
// It returns true for custom Error types with ErrorTypeConnection,
// or for standard library network errors.
func IsConnectionError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypeConnection
	}
	return isNetworkError(err)
}

// IsTimeoutError checks if the given error is a timeout error.
// It returns true for custom Error types with ErrorTypeTimeout,
// or for errors that contain "timeout" in their message.
func IsTimeoutError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypeTimeout
	}
	return isNetworkError(err) && strings.Contains(strings.ToLower(err.Error()), ErrorMessages[ErrorTypeTimeout])
}

// IsPermissionError checks if the given error is a permission error.
// It returns true for custom Error types with ErrorTypePermission,
// or for errors that contain "permission" in their message.
func IsPermissionError(err error) bool {
	var e *Error
	if errors.As(err, &e) {
		return e.Type == ErrorTypePermission
	}
	return strings.Contains(strings.ToLower(err.Error()), ErrorMessages[ErrorTypePermission])
}

// isNetworkError checks if the error is a known network error.
func isNetworkError(err error) bool {
	return errors.Is(err, net.ErrClosed) || errors.Is(err, net.ErrWriteToConnected)
}