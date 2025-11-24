package app

import (
	"fmt"
	"os"
	"time"

	"github.com/Vlad-Ali/Movies-service-back/internal/application/usecase/review/modelconfig"
	"github.com/Vlad-Ali/Movies-service-back/internal/infrastruture/postgresconfig"
	"gopkg.in/yaml.v3"
)

type Config struct {
	SecretKey      string                        `yaml:"secret_key"`
	Address        string                        `yaml:"address"`
	ReadTimeout    time.Duration                 `yaml:"read_timeout"`
	WriteTimeout   time.Duration                 `yaml:"write_timeout"`
	AllowedOrigins []string                      `yaml:"allowed_origins"`
	PostgresConfig postgresconfig.PostgresConfig `yaml:"postgres"`
	ModelConfig    modelconfig.ModelConfig       `yaml:"model"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}
	return &config, nil
}
