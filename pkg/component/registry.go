package component

import (
	"os"
	"path/filepath"
)

// Registry stores all loaded ComponentDefinitions
type Registry struct {
	definitions map[string]*Definition
}

func NewRegistry() *Registry {
	return &Registry{
		definitions: make(map[string]*Definition),
	}
}

// LoadAll loads all component definitions from a directory
func (r *Registry) LoadAll(dir string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}

		ext := filepath.Ext(path)
		if ext != ".xml" && ext != ".hcml" {
			return nil
		}

		def, err := LoadDefinition(path)
		if err != nil {
			return err
		}

		r.Register(def)
		return nil
	})
}

func (r *Registry) Register(def *Definition) {
	r.definitions[def.Name] = def
}

func (r *Registry) Get(name string) *Definition {
	return r.definitions[name]
}
