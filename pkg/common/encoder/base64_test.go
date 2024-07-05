package encoder_test

import (
	"pkg/common/encoder"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBase64Encoder(t *testing.T) {
	t.Run("NewBase64Encoder", func(t *testing.T) {
		enc := encoder.NewBase64Encoder()
		assert.NotNil(t, enc, "NewBase64Encoder should return a non-nil encoder")
	})

	t.Run("Encode", func(t *testing.T) {
		enc := encoder.NewBase64Encoder()
		testCases := []struct {
			input    string
			expected string
		}{
			{"Hello, World!", "SGVsbG8sIFdvcmxkIQ=="},
			{"", ""},
			{"1234567890", "MTIzNDU2Nzg5MA=="},
			{"!@#$%^&*()", "IUAjJCVeJiooKQ=="},
		}

		for _, tc := range testCases {
			t.Run(tc.input, func(t *testing.T) {
				result := enc.Encode(tc.input)
				assert.Equal(t, tc.expected, result, "Encode(%q) = %q; want %q", tc.input, result, tc.expected)
			})
		}
	})

	t.Run("Decode", func(t *testing.T) {
		enc := encoder.NewBase64Encoder()
		testCases := []struct {
			input    string
			expected string
			hasError bool
		}{
			{"SGVsbG8sIFdvcmxkIQ==", "Hello, World!", false},
			{"", "", false},
			{"MTIzNDU2Nzg5MA==", "1234567890", false},
			{"IUAjJCVeJiooKQ==", "!@#$%^&*()", false},
			{"Invalid Base64!", "", true},
		}

		for _, tc := range testCases {
			t.Run(tc.input, func(t *testing.T) {
				result, err := enc.Decode(tc.input)
				if tc.hasError {
					assert.Error(t, err, "Decode(%q) should return an error", tc.input)
				} else {
					assert.NoError(t, err, "Decode(%q) should not return an error", tc.input)
					assert.Equal(t, tc.expected, result, "Decode(%q) = %q; want %q", tc.input, result, tc.expected)
				}
			})
		}
	})

	t.Run("EncodeDecode", func(t *testing.T) {
		enc := encoder.NewBase64Encoder()
		testCases := []string{
			"Hello, World!",
			"",
			"1234567890",
			"!@#$%^&*()",
			"This is a longer string with spaces and punctuation.",
		}

		for _, tc := range testCases {
			t.Run(tc, func(t *testing.T) {
				encoded := enc.Encode(tc)
				decoded, err := enc.Decode(encoded)
				assert.NoError(t, err, "Decode(Encode(%q)) should not return an error", tc)
				assert.Equal(t, tc, decoded, "Decode(Encode(%q)) = %q; want %q", tc, decoded, tc)
			})
		}
	})
}
