package vars

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".staticgen.yml")

	configContent := `variables:
  siteName: "Test Site"
  author: "Test Author"
  year: "2026"
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err)

	// Load the config
	config, err := LoadConfig(configPath)
	require.NoError(t, err)
	require.NotNil(t, config)

	// Verify the variables were loaded
	assert.Equal(t, "Test Site", config.Variables["siteName"])
	assert.Equal(t, "Test Author", config.Variables["author"])
	assert.Equal(t, "2026", config.Variables["year"])
	assert.Len(t, config.Variables, 3)
}

func TestLoadConfigFileNotFound(t *testing.T) {
	_, err := LoadConfig("nonexistent.yml")
	assert.Error(t, err)
	assert.True(t, os.IsNotExist(err))
}

func TestLoadConfigInvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, ".staticgen.yml")

	invalidContent := `variables:
  - this is invalid
    yaml: [structure
`
	err := os.WriteFile(configPath, []byte(invalidContent), 0644)
	require.NoError(t, err)

	_, err = LoadConfig(configPath)
	assert.Error(t, err)
}

func TestStoreLoadFromConfig(t *testing.T) {
	store := NewStore()

	config := &Config{
		Variables: map[string]string{
			"key1": "value1",
			"key2": "value2",
			"key3": "value3",
		},
	}

	store.LoadFromConfig(config)

	// Verify all variables were loaded
	assert.Equal(t, "value1", store.GetOrDefault("key1", ""))
	assert.Equal(t, "value2", store.GetOrDefault("key2", ""))
	assert.Equal(t, "value3", store.GetOrDefault("key3", ""))
}

func TestStoreLoadFromConfigNil(t *testing.T) {
	store := NewStore()
	store.Set("existing", "value")

	// Loading nil config should not crash or clear existing values
	store.LoadFromConfig(nil)
	assert.Equal(t, "value", store.GetOrDefault("existing", ""))

	// Loading config with nil Variables should not crash
	store.LoadFromConfig(&Config{Variables: nil})
	assert.Equal(t, "value", store.GetOrDefault("existing", ""))
}
