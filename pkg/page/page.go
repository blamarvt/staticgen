package page

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/blamarvt/staticgen/pkg/component"
	"github.com/blamarvt/staticgen/pkg/internal/xmlutil"
)

// Page represents a complete page with metadata and components
type Page struct {
	Title      string
	Path       string
	Components []*component.Instance
}

// LoadPage parses a page XML file into a Page with Component instances
func LoadPage(filepath string, registry *component.Registry) (*Page, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "reading page file")
	}

	// Parse the XML structure
	root, err := xmlutil.ParseXML(data)
	if err != nil {
		return nil, errors.Wrap(err, "parsing page XML")
	}

	// Verify root element is "page"
	if root.XMLName.Local != "page" {
		return nil, fmt.Errorf("root element must be 'page', got '%s'", root.XMLName.Local)
	}

	// Extract page metadata from attributes
	page := &Page{
		Components: []*component.Instance{},
	}

	if title, ok := root.GetAttr("title"); ok {
		page.Title = title
	}
	if path, ok := root.GetAttr("path"); ok {
		page.Path = path
	}

	// Parse child elements as components
	components, err := parseComponents(root.Content, registry)
	if err != nil {
		return nil, errors.Wrap(err, "parsing page components")
	}
	page.Components = components

	return page, nil
}

// parseComponents parses XML content into Component instances
func parseComponents(xmlContent []byte, registry *component.Registry) ([]*component.Instance, error) {
	// Wrap content in a root element for parsing
	wrappedXML := "<root>" + string(xmlContent) + "</root>"

	wrapper, err := xmlutil.ParseXML([]byte(wrappedXML))
	if err != nil {
		return nil, fmt.Errorf("failed to parse components: %w", err)
	}

	var components []*component.Instance
	for _, node := range wrapper.Children {
		comp, err := parseComponent(node, registry)
		if err != nil {
			return nil, err
		}
		if comp != nil {
			components = append(components, comp)
		}
	}

	return components, nil
}

// parseComponent converts an xmlNode to a Component instance
func parseComponent(node xmlutil.Node, registry *component.Registry) (*component.Instance, error) {
	// Extract component name from namespace (e.g., "component:titleBar" -> "titleBar")
	componentName := node.XMLName.Local

	// Skip empty text nodes
	if strings.TrimSpace(componentName) == "" {
		return nil, nil
	}

	// Check if this is a component (has the component namespace) or plain HTML
	isComponent := node.XMLName.Space != ""

	// If it's not a component, treat it as raw HTML
	if !isComponent {
		// Check if this element has a "slot" attribute
		if slotName, hasSlot := node.GetAttr("slot"); hasSlot {
			// This is a slot element - return only the inner content
			// Don't include the wrapper element itself
			return &component.Instance{
				DefinitionName: "__slot__", // Special marker
				Attributes:     map[string]string{"name": slotName},
				RawHTML:        string(node.Content), // Just the inner content
			}, nil
		}

		// Reconstruct the HTML for this node
		html, err := reconstructHTML(node)
		if err != nil {
			return nil, err
		}
		return &component.Instance{
			RawHTML: html,
		}, nil
	}

	// Create component instance
	comp := &component.Instance{
		DefinitionName: componentName,
		Attributes:     make(map[string]string),
		Children:       []*component.Instance{},
		Slots:          make(map[string]string),
	}

	// Extract attributes
	for _, attr := range node.Attrs {
		// Skip namespace declarations
		if attr.Name.Space == "xmlns" || attr.Name.Local == "hcmlns" {
			continue
		}
		comp.Attributes[attr.Name.Local] = attr.Value
	}

	// Parse nested components recursively
	if len(node.Content) > 0 {
		children, err := parseComponents(node.Content, registry)
		if err != nil {
			return nil, err
		}

		// Separate slots from regular children
		for _, child := range children {
			// Check if this child is a slot (marked with __slot__ DefinitionName)
			if child.DefinitionName == "__slot__" {
				slotName := child.Attributes["name"]
				comp.Slots[slotName] = child.RawHTML
				continue
			}
			// Regular child
			comp.Children = append(comp.Children, child)
		}
	}

	return comp, nil
}

// reconstructHTML rebuilds HTML from an XML node
func reconstructHTML(node xmlutil.Node) (string, error) {
	var html strings.Builder

	// Opening tag
	html.WriteString("<")
	html.WriteString(node.XMLName.Local)

	// Add attributes
	for _, attr := range node.Attrs {
		// Skip namespace declarations
		if attr.Name.Space == "xmlns" || attr.Name.Local == "hcmlns" {
			continue
		}
		html.WriteString(" ")
		html.WriteString(attr.Name.Local)
		html.WriteString("=\"")
		html.WriteString(attr.Value)
		html.WriteString("\"")
	}

	// Self-closing tag if no content
	if len(node.Content) == 0 {
		html.WriteString(" />")
		return html.String(), nil
	}

	html.WriteString(">")

	// Add content (which may contain nested elements)
	html.Write(node.Content)

	// Closing tag
	html.WriteString("</")
	html.WriteString(node.XMLName.Local)
	html.WriteString(">")

	return html.String(), nil
}
