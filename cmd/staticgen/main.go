package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/blamarvt/staticgen/pkg/component"
	"github.com/blamarvt/staticgen/pkg/page"
)

func main() {
	outputDir := flag.String("output", "dist", "output directory for generated pages")
	flag.Parse()

	registry := component.NewRegistry()
	if err := registry.LoadAll("templates"); err != nil {
		log.Fatal(errors.Wrap(err, "loading components"))
	}

	pagesDir := "pages"
	err := filepath.WalkDir(pagesDir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || filepath.Ext(path) != ".hcml" {
			return errors.Wrap(err, "walking pages directory")
		}

		p, err := page.LoadPage(path, registry)
		if err != nil {
			return errors.Wrap(err, "loading page "+path)
		}

		html, err := page.Generate(p, registry)
		if err != nil {
			return errors.Wrap(err, "generating page "+path)
		}

		// Use the path attribute from the page, or derive from file location
		var outPath string
		if p.Path != "" {
			outPath = filepath.Join(*outputDir, filepath.FromSlash(p.Path))
		} else {
			// Derive output path from source file location relative to pages dir
			relPath, err := filepath.Rel(pagesDir, path)
			if err != nil {
				return errors.Wrap(err, "getting relative path for "+path)
			}
			// Replace .hcml extension with .html
			outPath = filepath.Join(*outputDir, relPath[:len(relPath)-len(filepath.Ext(relPath))]+".html")
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return errors.Wrap(err, "creating directories for "+outPath)
		}

		return os.WriteFile(outPath, []byte(html), 0644)
	})

	if err != nil {
		log.Fatal(errors.Wrap(err, "walking pages directory"))
	}
}
