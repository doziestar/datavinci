// Package retry provides a configurable mechanism for retrying operations
// that may fail due to transient errors.
package retry

import (
	"context"
	"time"

	"datavinci/pkg/common/errors"
)

// Config represents the configuration for the retry mechanism.
type Config struct {
	// MaxAttempts is the maximum number of attempts to make before giving up.
	MaxAttempts int
	// InitialBackoff is the duration to wait before the first retry.
	InitialBackoff time.Duration
	// MaxBackoff is the maximum duration to wait between retries.
	MaxBackoff time.Duration
	// BackoffFactor is the factor by which to increase the backoff duration after each retry.
	BackoffFactor float64
	// RetryableErrors is a list of functions that determine if an error is retryable.
	RetryableErrors []func(error) bool
	// NonRetryableErrors is a list of functions that determine if an error should not be retried.
	NonRetryableErrors []func(error) bool
}

// DefaultConfig returns a default retry configuration.
// It sets up sensible defaults for retry attempts, backoff durations,
// and includes common retryable and non-retryable error checks.
func DefaultConfig() *Config {
	return &Config{
		MaxAttempts:    5,
		InitialBackoff: 100 * time.Millisecond,
		MaxBackoff:     10 * time.Second,
		BackoffFactor:  2,
		RetryableErrors: []func(error) bool{
			errors.IsConnectionError,
			errors.IsTimeoutError,
		},
		NonRetryableErrors: []func(error) bool{
			errors.IsPermissionError,
		},
	}
}

// Retry executes the given function with retries based on the provided configuration.
// It will retry the function until it succeeds, the maximum number of attempts is reached,
// or the context is canceled.
//
// The function respects the context's cancellation and will return early if the context
// is canceled.
//
// If the function returns a non-retryable error, Retry will return immediately without
// further attempts.
func Retry(ctx context.Context, fn func() error, config *Config) error {
	var err error
	attempt := 0
	backoff := config.InitialBackoff

	for attempt < config.MaxAttempts {
		err = fn()
		if err == nil {
			return nil
		}

		if isNonRetryableError(err, config.NonRetryableErrors) {
			return err
		}

		if !isRetryableError(err, config.RetryableErrors) {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(backoff):
			attempt++
			backoff = calculateNextBackoff(backoff, config)
		}
	}

	return errors.NewError(errors.ErrorTypeUnknown, "max retry attempts reached", err)
}

// isRetryableError checks if the given error is retryable based on the provided list of check functions.
func isRetryableError(err error, retryableErrors []func(error) bool) bool {
	for _, isRetryable := range retryableErrors {
		if isRetryable(err) {
			return true
		}
	}
	return false
}

// isNonRetryableError checks if the given error is non-retryable based on the provided list of check functions.
func isNonRetryableError(err error, nonRetryableErrors []func(error) bool) bool {
	for _, isNonRetryable := range nonRetryableErrors {
		if isNonRetryable(err) {
			return true
		}
	}
	return false
}

// calculateNextBackoff calculates the next backoff duration using exponential backoff.
// It ensures that the backoff duration does not exceed the maximum specified in the config.
func calculateNextBackoff(currentBackoff time.Duration, config *Config) time.Duration {
	nextBackoff := time.Duration(float64(currentBackoff) * config.BackoffFactor)
	if nextBackoff > config.MaxBackoff {
		return config.MaxBackoff
	}
	return nextBackoff
}
