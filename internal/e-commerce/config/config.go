package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type CleanerConfig struct {
	Textures map[string]struct {
		Cleaner string `yaml:"cleaner"`
	} `yaml:"textures"`
}

func LoadCleanerConfig(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg CleanerConfig
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, err
	}

	cleanerMap := make(map[string]string)
	for k, v := range cfg.Textures {
		cleanerMap[k] = v.Cleaner
	}

	return cleanerMap, nil
}
