package page

import (
	"strings"

	"github.com/blamarvt/staticgen/pkg/component"
)

// Generate creates the final HTML from a page
func Generate(p *Page, registry *component.Registry) (string, error) {
	var html strings.Builder

	// Render each component
	for _, comp := range p.Components {
		rendered, err := comp.Render(registry)
		if err != nil {
			return "", err
		}
		html.WriteString(rendered)
	}

	return html.String(), nil
}
