package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_FILENAME = "config.yaml"
)

type KdbConfig struct {
	Storage   KdbStorageConfig   `yaml:"storage"`
	Embedding KdbEmbeddingConfig `yaml:"embedding"`
}

type KdbStorageConfig struct {
	Provider string `yaml:"provider"`
	URL      string `yaml:"url"`
}

type KdbEmbeddingConfig struct {
	Provider string                  `yaml:"provider"`
	URL      string                  `yaml:"url"`
	Model    KdbEmbeddingModelConfig `yaml:"model"`
}

type KdbEmbeddingModelConfig struct {
	Name       string `yaml:"name"`
	Dimensions int    `yaml:"dimensions"`
}

func LoadConfig(filename string) (*KdbConfig, error) {
	// Expand the ~ in the filename if present
	if filename[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("error getting home directory: %w", err)
		}
		filename = filepath.Join(home, filename[2:])
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var config KdbConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *KdbConfig, filename string) error {
	// Expand the ~ in the filename if present
	if filename[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("error getting home directory: %w", err)
		}
		filename = filepath.Join(home, filename[2:])
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("error marshaling config: %w", err)
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing config file: %w", err)
	}

	return nil
}

func UpdateConfig(filename string, updates map[string]interface{}) error {
	config, err := LoadConfig(filename)
	if err != nil {
		return err
	}

	// Update the config with the new values
	for key, value := range updates {
		switch key {
		case "storage.provider":
			config.Storage.Provider = value.(string)
		case "storage.url":
			config.Storage.URL = value.(string)
		case "embedding.provider":
			config.Embedding.Provider = value.(string)
		case "embedding.url":
			config.Embedding.URL = value.(string)
		case "embedding.model.name":
			config.Embedding.Model.Name = value.(string)
		case "embedding.model.dimensions":
			config.Embedding.Model.Dimensions = value.(int)
		default:
			return fmt.Errorf("unknown config key: %s", key)
		}
	}

	return SaveConfig(config, filename)
}

func NewDefaultConfig() *KdbConfig {
	return &KdbConfig{
		Storage: KdbStorageConfig{
			Provider: "DuckDb",
			URL:      "~/.kdb/knowledgebase.ddb",
		},
		Embedding: KdbEmbeddingConfig{
			Provider: "Ollama",
			URL:      "http://localhost:11434",
			Model: KdbEmbeddingModelConfig{
				Name:       "nomic-embed-text",
				Dimensions: 768,
			},
		},
	}
}
