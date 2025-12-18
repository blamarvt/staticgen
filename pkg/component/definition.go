package component

import (
	"os"
	"strings"

	"github.com/blamarvt/staticgen/pkg/internal/xmlutil"
)

// Definition is the template/function that defines what a component is
type Definition struct {
	Name      string
	Namespace string
	Template  string // The Go template text with {{ .Var }} placeholders
	// Schema info for validation
	RequiredAttrs []string
	OptionalAttrs []string
}

// LoadDefinition reads a component definition file
func LoadDefinition(filepath string) (*Definition, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	// Parse the XML structure
	root, err := xmlutil.ParseXML(data)
	if err != nil {
		return nil, err
	}

	// Component name is the root element's local name
	def := &Definition{
		Name:          root.XMLName.Local,
		Namespace:     root.GetNamespace(),
		RequiredAttrs: []string{},
		OptionalAttrs: []string{},
	}

	// Extract the inner content as the template
	// This is the HTML/template content inside the component definition
	def.Template = string(root.Content)

	// Clean up the template (remove extra whitespace at start/end)
	def.Template = strings.TrimSpace(def.Template)

	return def, nil
}
