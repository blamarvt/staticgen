package component

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Instance is an instantiation of a definition with specific attribute values
type Instance struct {
	DefinitionName string            // e.g., "titleBar"
	Attributes     map[string]string // e.g., {"icon": "fa-user", "text": "Add User"}
	Children       []*Instance       // Nested components
	Slots          map[string]string // Named slots for injecting markup content
	RawHTML        string            // Raw HTML content (for non-component elements)
}

// Render generates HTML by applying attributes to the definition's template
func (c *Instance) Render(registry *Registry) (string, error) {
	// If this is a raw HTML instance, just return the HTML directly
	if c.RawHTML != "" {
		return c.RawHTML, nil
	}

	def := registry.Get(c.DefinitionName)
	if def == nil {
		return "", fmt.Errorf("component definition not found: %s", c.DefinitionName)
	}

	// Render all children first
	var childrenHTML strings.Builder
	for _, child := range c.Children {
		childHTML, err := child.Render(registry)
		if err != nil {
			return "", fmt.Errorf("failed to render child component: %w", err)
		}
		childrenHTML.WriteString(childHTML)
	}

	// Create template data with attributes and children
	templateData := make(map[string]interface{})

	// Add all attributes with capitalized first letter for Go template convention
	for key, value := range c.Attributes {
		// Capitalize first letter: "icon" -> "Icon"
		capitalizedKey := strings.ToUpper(key[:1]) + key[1:]
		templateData[capitalizedKey] = value
	}

	// Add rendered children
	templateData["Children"] = childrenHTML.String()

	// Create template with custom functions
	funcMap := template.FuncMap{
		"Slot": func(name string) string {
			if content, exists := c.Slots[name]; exists {
				return content
			}
			return "" // Return empty string if slot not defined
		},
	}

	// Parse and execute the template with custom functions
	tmpl, err := template.New(c.DefinitionName).Funcs(funcMap).Parse(def.Template)
	if err != nil {
		return "", fmt.Errorf("failed to parse template for %s: %w", c.DefinitionName, err)
	}

	var output bytes.Buffer
	if err := tmpl.Execute(&output, templateData); err != nil {
		return "", fmt.Errorf("failed to execute template for %s: %w", c.DefinitionName, err)
	}

	return output.String(), nil
}
