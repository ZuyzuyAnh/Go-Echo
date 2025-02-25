package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type (
	Config struct {
		DB_DSN     string `mapstructure:"DB_DSN"`
		ServerHost string `mapstructure:"SERVER_HOST"`
		ServerPort int    `mapstructure:"SERVER_PORT"`
		JWTSecret  string `mapstructure:"JWT_SECRET"`
	}
)

func LoadConfig(filename string) (*Config, error) {
	viper.SetConfigFile(filename)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error loading .env file: %w", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	return &config, nil
}
