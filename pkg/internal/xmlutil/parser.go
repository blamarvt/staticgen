package xmlutil

import (
	"encoding/xml"
	"fmt"
)

// Node represents a generic XML element for parsing
type Node struct {
	XMLName  xml.Name
	Attrs    []xml.Attr `xml:",any,attr"`
	Content  []byte     `xml:",innerxml"`
	Children []Node     `xml:",any"`
}

// ParseXML parses XML data into a Node structure
func ParseXML(data []byte) (*Node, error) {
	var root Node
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("failed to parse XML: %w", err)
	}
	return &root, nil
}

// GetAttr retrieves an attribute value by name
func (n *Node) GetAttr(name string) (string, bool) {
	for _, attr := range n.Attrs {
		if attr.Name.Local == name {
			return attr.Value, true
		}
	}
	return "", false
}

// GetNamespace retrieves the namespace URI
func (n *Node) GetNamespace() string {
	return n.XMLName.Space
}
