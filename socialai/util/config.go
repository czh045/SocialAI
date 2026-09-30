package util

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ElasticsearchInfo struct {
	Address  string `yaml:"address"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type StorageInfo struct {
	Driver        string `yaml:"driver"`
	Directory     string `yaml:"directory"`
	PublicBaseURL string `yaml:"public_base_url"`
}

type GCSInfo struct {
	Bucket string `yaml:"bucket"`
}

type TokenInfo struct {
	Secret string `yaml:"secret"`
}

type ApplicationConfig struct {
	ElasticsearchConfig *ElasticsearchInfo `yaml:"elasticsearch"`
	StorageConfig       *StorageInfo       `yaml:"storage"`
	GCSConfig           *GCSInfo           `yaml:"gcs"`
	TokenConfig         *TokenInfo         `yaml:"token"`
}

// LoadApplicationConfig reads the YAML configuration used by local development
// and by the GAE deployment described in the course.
func LoadApplicationConfig(configPath string) (*ApplicationConfig, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config ApplicationConfig
	if err := yaml.Unmarshal(content, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
