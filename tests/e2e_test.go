package tests

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/blamarvt/staticgen/pkg/component"
	"github.com/blamarvt/staticgen/pkg/htmlutil"
	"github.com/blamarvt/staticgen/pkg/page"
)

func TestEndToEndSimplePage(t *testing.T) {
	// Create registry and load component definitions
	registry := component.NewRegistry()
	err := registry.LoadAll("fixtures/components")
	require.NoError(t, err, "Failed to load component definitions")

	greetingDef := registry.Get("greeting")
	require.NotNil(t, greetingDef, "Greeting component definition should be loaded")

	// Load the simple page
	p, err := page.LoadPage("fixtures/pages/simple.hcml", registry)
	require.NoError(t, err, "Failed to load page")

	assert.Equal(t, "Test Page", p.Title, "Page title mismatch")
	assert.Equal(t, "/test.html", p.Path, "Page path mismatch")

	// Generate HTML
	html, err := page.Generate(p, registry)
	require.NoError(t, err, "Failed to generate HTML")

	assert.Equal(
		t,
		htmlutil.MustNormalize(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>Test Page</title>
				</head>
				<body>
					<div class="greeting">
						<h1>Hello, World!</h1>
						<p>Welcome to our test!</p>
					</div>
				</body>
			</html>
		`),
		htmlutil.MustNormalize(html),
	)
}

func TestEndToEndMixedContent(t *testing.T) {
	// Create registry and load component definitions
	registry := component.NewRegistry()
	err := registry.LoadAll("fixtures/components")
	require.NoError(t, err, "Failed to load component definitions")

	// Load the mixed content page
	p, err := page.LoadPage("fixtures/pages/mixed.hcml", registry)
	require.NoError(t, err, "Failed to load page")

	assert.Equal(t, "/mixed.html", p.Path, "Page path mismatch")

	// Generate HTML
	html, err := page.Generate(p, registry)
	require.NoError(t, err, "Failed to generate HTML")

	assert.Equal(
		t,
		htmlutil.MustNormalize(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>Mixed Content Test</title>
				</head>
				<body>
					<div class="container">
						<h2>Container with Mixed Content</h2>
						<div class="plain-html">
							<p>This is plain HTML inside a component</p>
							<span>No template needed!</span>
						</div>
						<div class="greeting">
							<h1>Hello, Alice!</h1>
							<p>Components work too!</p>
						</div>
						<section>
							<h3>More plain HTML</h3>
							<ul>
								<li>Item 1</li>
								<li>Item 2</li>
							</ul>
						</section>
					</div>
				</body>
			</html>
		`),
		htmlutil.MustNormalize(html),
	)
}

func TestEndToEndNestedComponents(t *testing.T) {
	// Create registry and load component definitions
	registry := component.NewRegistry()
	err := registry.LoadAll("fixtures/components")
	require.NoError(t, err, "Failed to load component definitions")

	// Load the nested page
	p, err := page.LoadPage("fixtures/pages/nested.hcml", registry)
	require.NoError(t, err, "Failed to load page")

	assert.Equal(t, "Nested Test", p.Title, "Page title mismatch")

	// Generate HTML
	html, err := page.Generate(p, registry)
	require.NoError(t, err, "Failed to generate HTML")

	assert.Equal(
		t,
		htmlutil.MustNormalize(`
			<!DOCTYPE html>
			<html>
				<head>
					<title>Nested Test</title>
				</head>
				<body>
					<div class="container">
						<h2>Main Container</h2>
						<div class="greeting">
							<h1>Hello, User!</h1>
							<p>This is nested!</p>
						</div>
					</div>
				</body>
			</html>
		`),
		htmlutil.MustNormalize(html),
	)
}

func TestEndToEndComponentNotFound(t *testing.T) {
	// Create registry with limited components
	registry := component.NewRegistry()

	// Load a page that references a non-existent component
	p, err := page.LoadPage("fixtures/pages/simple.hcml", registry)
	require.NoError(t, err, "Failed to load page")

	// Try to generate HTML without loading definitions
	_, err = page.Generate(p, registry)
	assert.Error(t, err, "Expected error when component definition not found")
	assert.Contains(t, err.Error(), "component definition not found", "Error message should mention missing component")
}

func TestEndToEndAllFixtures(t *testing.T) {
	// Create registry and load all component definitions
	registry := component.NewRegistry()
	err := registry.LoadAll("fixtures/components")
	require.NoError(t, err, "Failed to load component definitions")

	// Test all page fixtures
	pages := []string{
		"fixtures/pages/simple.hcml",
		"fixtures/pages/nested.hcml",
		"fixtures/pages/mixed.hcml",
	}

	for _, pagePath := range pages {
		t.Run(filepath.Base(pagePath), func(t *testing.T) {
			p, err := page.LoadPage(pagePath, registry)
			require.NoError(t, err, "Failed to load page")

			html, err := page.Generate(p, registry)
			require.NoError(t, err, "Failed to generate HTML")

			// Basic sanity checks
			assert.Contains(t, html, "<!DOCTYPE html>", "HTML should contain DOCTYPE")
			assert.Contains(t, html, "<html>", "HTML should contain opening html tag")
			assert.Contains(t, html, "</html>", "HTML should contain closing html tag")
			assert.Contains(t, html, "<title>", "HTML should contain opening title tag")
			assert.Contains(t, html, "</title>", "HTML should contain closing title tag")
			assert.Contains(t, html, "<body>", "HTML should contain opening body tag")
			assert.Contains(t, html, "</body>", "HTML should contain closing body tag")
		})
	}
}
