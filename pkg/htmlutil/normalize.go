package htmlutil

import (
	"bytes"
	"strings"

	"github.com/yosssi/gohtml"
	"golang.org/x/net/html"
)

// MustNormalize normalizes HTML content by parsing and re-rendering it
// This removes extra whitespace, formats consistently, and validates structure
func MustNormalize(input string) string {
	// Parse the HTML
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		panic(err)
	}

	// Render back to normalized HTML
	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		panic(err)
	}

	return gohtml.Format(buf.String())
}
