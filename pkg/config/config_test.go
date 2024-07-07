package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestConfig(t *testing.T, content string) (cleanup func(), configPath string) {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "datavinci-config-test")
	require.NoError(t, err, "Failed to create temp directory")

	configPath = filepath.Join(tempDir, "config.yaml")
	err = os.WriteFile(configPath, []byte(content), 0644)
	require.NoError(t, err, "Failed to write temp config file")

	originalWD, err := os.Getwd()
	require.NoError(t, err, "Failed to get current working directory")

	err = os.Chdir(tempDir)
	require.NoError(t, err, "Failed to change working directory")

	viper.Reset()
	viper.SetConfigFile(configPath)

	cleanup = func() {
		os.Chdir(originalWD)
		os.RemoveAll(tempDir)
	}

	return cleanup, configPath
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name           string
		configContent  string
		expectedConfig *Config
		expectError    bool
		errorContains  string
	}{
		{
			name: "Valid full config",
			configContent: `
DatabaseURL: "postgres://user:pass@localhost:5432/datavinci"
AuthServiceAddress: "auth.datavinci.com:8080"
AuthzServiceAddress: "authz.datavinci.com:8081"
JWTSecret: "datavinci-secret-key"
`,
			expectedConfig: &Config{
				DatabaseURL:         "postgres://user:pass@localhost:5432/datavinci",
				AuthServiceAddress:  "auth.datavinci.com:8080",
				AuthzServiceAddress: "authz.datavinci.com:8081",
				JWTSecret:           "datavinci-secret-key",
			},
			expectError: false,
		},
		{
			name:           "Empty config file",
			configContent:  "",
			expectedConfig: &Config{},
			expectError:    false,
		},
		{
			name: "Partial config",
			configContent: `
DatabaseURL: "postgres://user:pass@localhost:5432/datavinci"
AuthServiceAddress: "auth.datavinci.com:8080"
`,
			expectedConfig: &Config{
				DatabaseURL:        "postgres://user:pass@localhost:5432/datavinci",
				AuthServiceAddress: "auth.datavinci.com:8080",
			},
			expectError: false,
		},
		{
			name: "Invalid YAML",
			configContent: `
DatabaseURL: "postgres://user:pass@localhost:5432/datavinci"
AuthServiceAddress: "auth.datavinci.com:8080"
  IndentedLine
`,
			expectError:   true,
			errorContains: "invalid YAML content",
		},
		{
			name: "Extra fields in config",
			configContent: `
DatabaseURL: "postgres://user:pass@localhost:5432/datavinci"
AuthServiceAddress: "auth.datavinci.com:8080"
ExtraField: "extra value"
`,
			expectedConfig: &Config{
				DatabaseURL:        "postgres://user:pass@localhost:5432/datavinci",
				AuthServiceAddress: "auth.datavinci.com:8080",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup, _ := setupTestConfig(t, tt.configContent)
			defer cleanup()

			cfg, err := Load()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					// assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedConfig, cfg)
			}
		})
	}
}

func TestLoadConfigNotFound(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "datavinci-config-not-found")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	originalWD, err := os.Getwd()
	require.NoError(t, err)
	defer os.Chdir(originalWD)

	err = os.Chdir(tempDir)
	require.NoError(t, err)

	viper.Reset()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	cfg, err := Load()
	assert.NoError(t, err, "Expected no error when config file is not found")
	assert.Equal(t, &Config{}, cfg, "Expected empty config when file is not found")
}

func TestLoadEnvironmentVariables(t *testing.T) {
	cleanup, _ := setupTestConfig(t, "")
	defer cleanup()

	// Set environment variables
	os.Setenv("DATABASEURL", "postgres://env:pass@localhost:5432/envdb")
	os.Setenv("AUTHSERVICEADDRESS", "envauth.datavinci.com:8080")
	defer func() {
		os.Unsetenv("DATABASEURL")
		os.Unsetenv("AUTHSERVICEADDRESS")
	}()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, &Config{
		DatabaseURL:        "postgres://env:pass@localhost:5432/envdb",
		AuthServiceAddress: "envauth.datavinci.com:8080",
	}, cfg)
}

func TestLoadConfigPrecedence(t *testing.T) {
	configContent := `
DatabaseURL: "postgres://file:pass@localhost:5432/filedb"
AuthServiceAddress: "fileauth.datavinci.com:8080"
`
	cleanup, _ := setupTestConfig(t, configContent)
	defer cleanup()

	// Set environment variables
	os.Setenv("DATABASEURL", "postgres://env:pass@localhost:5432/envdb")
	defer os.Unsetenv("DATABASEURL")

	viper.AutomaticEnv()

	cfg, err := Load()
	assert.NoError(t, err)
	assert.Equal(t, &Config{
		DatabaseURL:        "postgres://env:pass@localhost:5432/envdb", // From env
		AuthServiceAddress: "fileauth.datavinci.com:8080",              // From file
	}, cfg)
}
