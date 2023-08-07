package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type AgentConfig struct {
	Port   string `env:"PORT" env-default:"8080"`
	APIKey string
}

func ReadAgentConfig() (*AgentConfig, error) {
	var config AgentConfig

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
