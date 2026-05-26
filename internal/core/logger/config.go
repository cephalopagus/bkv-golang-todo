package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	LogLevel  string `envconfig:"LEVEL" required:"true"`
	LogFolder string `envconfig:"FOLDER" required:"true"`
}

func NewConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("LOGGER", &cfg); err != nil {
		return Config{}, fmt.Errorf("Process envconfig: %w", err)
	}
	return cfg, nil
}
func NewConfigMust() Config {
	cfg, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("Get Logger config: %w", err))
	}
	return cfg
}
