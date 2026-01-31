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

	err := filepath.WalkDir("pages", func(path string, d os.DirEntry, err error) error {
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

		outPath := filepath.Join(*outputDir, p.Path)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		return os.WriteFile(outPath, []byte(html), 0644)
	})

	if err != nil {
		log.Fatal(err)
	}
}
