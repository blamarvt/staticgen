package page

import (
	"strings"

	"codeberg.org/derat/htmlpretty"
	"golang.org/x/net/html"

	"github.com/blamarvt/staticgen/pkg/component"
	"github.com/blamarvt/staticgen/pkg/vars"
)

// Generate creates the final HTML from a page
func Generate(p *Page, registry *component.Registry, variables *vars.Store) (string, error) {
	var hb strings.Builder

	// Render each component
	for _, comp := range p.Components {
		rendered, err := comp.Render(registry, variables)
		if err != nil {
			return "", err
		}
		hb.WriteString(rendered)
	}

	parsed, err := html.Parse(strings.NewReader(hb.String()))
	if err != nil {
		return "", err
	}

	outputBuffer := &strings.Builder{}

	err = htmlpretty.Print(outputBuffer, parsed, "\t", 120)
	if err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}
