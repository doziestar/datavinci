package retry

import (
	"context"
	"testing"
	"time"

	"datavinci/pkg/common/errors"
)

func TestRetry(t *testing.T) {
	tests := []struct {
		name          string
		fn            func() error
		config        *Config
		expectedCalls int
		expectedError bool
	}{
		{
			name: "Success on first try",
			fn: func() error {
				return nil
			},
			config:        DefaultConfig(),
			expectedCalls: 1,
			expectedError: false,
		},
		{
			name: "Success after retries",
			fn: (func() func() error {
				count := 0
				return func() error {
					if count < 2 {
						count++
						return errors.NewError(errors.ErrorTypeConnection, "connection failed", nil)
					}
					return nil
				}
			})(),
			config:        DefaultConfig(),
			expectedCalls: 3,
			expectedError: false,
		},
		{
			name: "Max retries reached",
			fn: func() error {
				return errors.NewError(errors.ErrorTypeConnection, "connection failed", nil)
			},
			config:        DefaultConfig(),
			expectedCalls: 5,
			expectedError: true,
		},
		{
			name: "Non-retryable error",
			fn: func() error {
				return errors.NewError(errors.ErrorTypePermission, "permission denied", nil)
			},
			config:        DefaultConfig(),
			expectedCalls: 1,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			wrappedFn := func() error {
				calls++
				return tt.fn()
			}

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := Retry(ctx, wrappedFn, tt.config)

			if (err != nil) != tt.expectedError {
				t.Errorf("Retry() error = %v, expectedError %v", err, tt.expectedError)
			}

			if calls != tt.expectedCalls {
				t.Errorf("Retry() calls = %v, expectedCalls %v", calls, tt.expectedCalls)
			}
		})
	}
}
