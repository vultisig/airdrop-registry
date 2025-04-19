package config

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

// Config is a struct that holds the configuration for the application
type Config struct {
	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}
	MySQL struct {
		Database string `mapstructure:"database"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
	}
	Worker struct {
		StartID int64 `mapstructure:"start_id"`
		// we will have 2x concurrency workers (for active position and balance)
		Concurrency int64 `mapstructure:"concurrency"`
	}
	OpenSea struct {
		APIKey string `mapstructure:"api_key"`
	}
	Vultiref struct {
		APIKey      string `mapstructure:"api_key"`
		BaseAddress string `mapstructure:"base_address"`
	}
}

// LoadConfig loads application configuration from a YAML file, environment variables, and default values.
// It returns a populated Config struct or an error if unmarshaling fails. If the config file is missing, environment variables and defaults are used.
func LoadConfig() (*Config, error) {
	var cfg Config
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("mysql.database", "airdrop")
	viper.SetDefault("mysql.user", "root")
	viper.SetDefault("mysql.password", "password")
	viper.SetDefault("mysql.host", "localhost")
	viper.SetDefault("mysql.port", 3301)
	viper.SetDefault("worker.start_id", 0)
	viper.SetDefault("worker.concurrency", 10)

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("No config file found, using environment variables and defaults")
		}
	}

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}
	return &cfg, nil
}
