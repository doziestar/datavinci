package encoder

import "encoding/base64"

// Base64Encoder is a helper struct for encoding and decoding Base64 strings.
type Base64Encoder struct{}

// NewBase64Encoder creates a new Base64Encoder.
//
// Usage:
//
//	encoder := NewBase64Encoder()
func NewBase64Encoder() *Base64Encoder {
	return &Base64Encoder{}
}

// Encode encodes a string to Base64.
//
// Parameters:
//   - data: The string to encode.
//
// Returns:
//   - The Base64 encoded string.
//
// Usage:
//
//	encoded := encoder.Encode("Hello, World!")
func (b *Base64Encoder) Encode(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

// Decode decodes a Base64 string.
//
// Parameters:
//   - encodedData: The Base64 encoded string to decode.
//
// Returns:
//   - The decoded string and an error if decoding fails.
//
// Usage:
//
//	decoded, err := encoder.Decode(encodedString)
func (b *Base64Encoder) Decode(encodedData string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
