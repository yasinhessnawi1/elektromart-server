package config

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestLoadConfig tests loading of environment variables from a .env file
// It creates a temporary .env file, writes a variable to it, and then loads the file
// The test checks if the variable is loaded correctly
func TestLoadConfig(t *testing.T) {
	// Create a temporary .env file
	file, err := os.Create(".env")
	if err != nil {
		t.Fatalf("Failed to create temporary .env file: %v", err)
	}
	defer os.Remove(".env")

	_, err = file.WriteString("TEST_VAR=success\n")
	if err != nil {
		t.Fatalf("Failed to write to temporary .env file: %v", err)
	}
	file.Close()

	LoadConfig()

	value, exists := os.LookupEnv("TEST_VAR")
	assert.True(t, exists, "Environment variable TEST_VAR should exist")
	assert.Equal(t, "success", value, "Environment variable TEST_VAR should be 'success'")
}

// TestGetConfig checks if the correct value is returned for an environment variable
// It sets an environment variable, retrieves it using GetConfig, and checks if the value is correct
// It also tests the retrieval of a non-existent variable
func TestGetConfig(t *testing.T) {
	// Set environment variable
	os.Setenv("TEST_VAR", "value")
	defer os.Unsetenv("TEST_VAR")

	// Retrieve the value using GetConfig
	value := GetConfig("TEST_VAR")
	assert.Equal(t, "value", value, "GetConfig should return the correct value for TEST_VAR")

	// Test retrieval of a non existent variable
	missingValue := GetConfig("MISSING_VAR")
	assert.Emptyf(t, missingValue, "GetConvig should return an empty string for missing variavles")
}
