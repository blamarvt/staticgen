package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/blamarvt/staticgen/pkg/component"
	"github.com/blamarvt/staticgen/pkg/page"
)

func main() {
	outputDir := flag.String("output", "dist", "output directory for generated pages")
	flag.Parse()

	registry := component.NewRegistry()
	if err := registry.LoadAll("templates/components"); err != nil {
		log.Fatal(err)
	}

	pagesDir := "pages"
	err := filepath.WalkDir(pagesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".hcml" {
			return err
		}

		p, err := page.LoadPage(path, registry)
		if err != nil {
			return err
		}

		html, err := page.Generate(p, registry)
		if err != nil {
			return err
		}

		// Use the path attribute from the page, or derive from file location
		var outPath string
		if p.Path != "" {
			outPath = filepath.Join(*outputDir, filepath.FromSlash(p.Path))
		} else {
			// Derive output path from source file location relative to pages dir
			relPath, err := filepath.Rel(pagesDir, path)
			if err != nil {
				return err
			}
			// Replace .hcml extension with .html
			outPath = filepath.Join(*outputDir, relPath[:len(relPath)-len(filepath.Ext(relPath))]+".html")
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(outPath, []byte(html), 0644)
	})

	if err != nil {
		log.Fatal(err)
	}
}
