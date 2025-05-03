package config

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

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

	SessionAchievements struct  {
		StartSession       time.Time `mapstructure:"start_session"`
		EndSession         time.Time `mapstructure:"end_session"`
		SwapMultiplier     float32   `mapstructure:"swap_multiplier"`
		ReferralMultiplier float32   `mapstructure:"referral_multiplier"`
		Milestones         []int     `mapstructure:"milestones"`
		NFTs               []NFT     `mapstructure:"nfts"`
		Tokens             []Token   `mapstructure:"tokens"`
	}
}
type NFT struct {
	Multiplier      int    `mapstructure:"multiplier"`
	CollectionName  string `mapstructure:"collection_name"`
	Chain           string `mapstructure:"chain"`
	ContractAddress string `mapstructure:"contract_address"`
}

type Token struct {
	Multiplier      int    `mapstructure:"multiplier"`
	Name            string `mapstructure:"name"`
	MinTokens       int    `mapstructure:"min_tokens"`
	Chain           string `mapstructure:"chain"`
	ContractAddress string `mapstructure:"contract_address"`
}

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
	viper.SetDefault("sessionachievements.start_session", time.Now())
	viper.SetDefault("sessionachievements.end_session", time.Now().Add(24*time.Hour))
	viper.SetDefault("sessionachievements.swap_multiplier", 1.6)
	viper.SetDefault("sessionachievements.referral_multiplier", 1.5)
	viper.SetDefault("sessionachievements.milestones", []int{5000, 10000, 50000, 100000})
	viper.SetDefault("sessionachievements.nfts", []NFT{})
	viper.SetDefault("sessionachievements.tokens", []Token{})

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
