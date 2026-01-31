package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blamarvt/staticgen/pkg/component"
	"github.com/blamarvt/staticgen/pkg/page"
	"github.com/blamarvt/staticgen/pkg/vars"
)

func TestVariableFunctionality(t *testing.T) {
	// Create registry and load component definitions
	registry := component.NewRegistry()
	err := registry.LoadAll("fixtures/components")
	require.NoError(t, err, "Failed to load component definitions")

	// Load the simple page
	p, err := page.LoadPage("fixtures/pages/simple.hcml", registry)
	require.NoError(t, err, "Failed to load page")

	// Create a variables store and set some variables
	variables := vars.NewStore()
	variables.Set("siteName", "My Awesome Site")
	variables.Set("year", "2026")
	variables.Set("author", "Test User")

	// Generate HTML with variables
	html, err := page.Generate(p, registry, variables)
	require.NoError(t, err, "Failed to generate HTML")

	// Just verify it generates without error
	// Variables can be used in templates with {{ Var "variableName" }}
	assert.NotEmpty(t, html)
}

func TestVariableStore(t *testing.T) {
	store := vars.NewStore()

	// Test Set and Get
	store.Set("key1", "value1")
	value, ok := store.Get("key1")
	assert.True(t, ok, "Key should exist")
	assert.Equal(t, "value1", value)

	// Test non-existent key
	_, ok = store.Get("nonexistent")
	assert.False(t, ok, "Key should not exist")

	// Test GetOrDefault
	assert.Equal(t, "value1", store.GetOrDefault("key1", "default"))
	assert.Equal(t, "default", store.GetOrDefault("nonexistent", "default"))
}

func TestVariablesInTemplate(t *testing.T) {
	// Create registry and load component definitions (including footer)
	registry := component.NewRegistry()
	err := registry.LoadAll("fixtures/components")
	require.NoError(t, err, "Failed to load component definitions")

	// Load the variables page
	p, err := page.LoadPage("fixtures/pages/variables.hcml", registry)
	require.NoError(t, err, "Failed to load page")

	// Create a variables store and set some variables
	variables := vars.NewStore()
	variables.Set("siteName", "My Awesome Site")
	variables.Set("year", "2026")
	variables.Set("author", "John Doe")

	// Generate HTML with variables
	html, err := page.Generate(p, registry, variables)
	require.NoError(t, err, "Failed to generate HTML")

	// Verify the variables were rendered in the template
	assert.Contains(t, html, "2026", "Year variable should be rendered")
	assert.Contains(t, html, "My Awesome Site", "Site name variable should be rendered")
	assert.Contains(t, html, "John Doe", "Author variable should be rendered")
	assert.Contains(t, html, "© 2026 My Awesome Site", "Footer with variables should be rendered correctly")
}

func TestLoadConfigIntegration(t *testing.T) {
	// Load config from the root .staticgen.yml file
	config, err := vars.LoadConfig("../.staticgen.yml")
	require.NoError(t, err, "Failed to load config file")
	require.NotNil(t, config)

	// Create a store and load from config
	variables := vars.NewStore()
	variables.LoadFromConfig(config)

	// Verify the variable from config is available
	siteName, ok := variables.Get("siteName")
	assert.True(t, ok, "siteName should be loaded from config")
	assert.Equal(t, "My Static Site", siteName)

	// Load registry and page
	registry := component.NewRegistry()
	err = registry.LoadAll("fixtures/components")
	require.NoError(t, err, "Failed to load component definitions")

	p, err := page.LoadPage("fixtures/pages/variables.hcml", registry)
	require.NoError(t, err, "Failed to load page")

	// Add additional variables
	variables.Set("year", "2026")
	variables.Set("author", "Config User")

	// Generate HTML with variables from config
	html, err := page.Generate(p, registry, variables)
	require.NoError(t, err, "Failed to generate HTML")

	// Verify the config variable is in the output
	assert.Contains(t, html, "My Static Site", "Config siteName should be rendered")
}
