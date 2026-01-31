package vars

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of .staticgen.yml
type Config struct {
	Variables map[string]string `yaml:"variables"`
}

// LoadConfig reads the .staticgen.yml file and returns a Config
func LoadConfig(filepath string) (*Config, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadFromConfig populates the Store with variables from a Config
func (s *Store) LoadFromConfig(config *Config) {
	if config == nil || config.Variables == nil {
		return
	}

	for key, value := range config.Variables {
		s.Set(key, value)
	}
}
