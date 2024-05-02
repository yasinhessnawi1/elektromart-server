package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGenerateUUID by ensuring it returns a non zero value.
func TestGenerateUUID(t *testing.T) {
	// Generate a random UUID
	uuid := GenerateUUID()
	assert.NotZero(t, uuid, "GenerateUUID should return a non zero uint32")
}

func TestConvertStringToUint(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    uint32
		shouldError bool
	}{
		{"Valid Number", "12345", 12345, false},
		{"Zero Value", "0", 0, false},
		{"Negative Number", "-123", 0, true},
		{"Non numeric String", "abc123", 0, true},
		{"Overflow Number", "42958766587452", 0, true}, // This value exceeds uint32 max value
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ConvertStringToUint(test.input)
			if test.shouldError {
				assert.Zero(t, result, "Expected an error and returned of 0 for input", test.input)
			} else {
				assert.Equal(t, test.expected, result, "Expected correct conversion for input", test.input)
			}
		})
	}
}
