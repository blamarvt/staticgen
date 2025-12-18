package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/blamarvt/staticgen/pkg/component"
	"github.com/blamarvt/staticgen/pkg/page"
)

func main() {
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

		return os.WriteFile("output/"+p.Path, []byte(html), 0644)
	})

	if err != nil {
		log.Fatal(err)
	}
}
