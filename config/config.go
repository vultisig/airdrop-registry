package config

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
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
	ReferralBot struct {
		APIKey      string `mapstructure:"api_key"`
		BaseAddress string `mapstructure:"base_address"`
	}
	Seasons           []AirdropSeason `mapstructure:"seasons"`
	VolumeTrackingAPI struct {
		AffiliateAddress   []string `mapstructure:"affiliate_address"`
		EtherscanAPIKey    string   `mapstructure:"etherscan_api_key"`
		EthplorerAPIKey    string   `mapstructure:"ethplorer_api_key"`
		TCMidgardBaseURL   string   `mapstructure:"tcmidgard_base_url"`
		MayaMidgardBaseURL string   `mapstructure:"mayamidgard_base_url"`
	}
}

type NFT struct {
	Token          `mapstructure:",squash"`
	CollectionName string `mapstructure:"collection_name" json:"collection_name"`
}

type Token struct {
	Multiplier      float64 `mapstructure:"multiplier" json:"multiplier"` //boosting multiplier
	Name            string  `mapstructure:"name" json:"name"`
	Chain           string  `mapstructure:"chain" json:"chain"`
	ContractAddress string  `mapstructure:"contract_address" json:"contract_address"`
}

type AirdropSeason struct {
	ID         uint        `mapstructure:"id" json:"id"`
	Start      time.Time   `mapstructure:"start" json:"start"`
	End        time.Time   `mapstructure:"end" json:"end"`
	Milestones []Milestone `mapstructure:"milestones" json:"milestones"` // list of vulti milestones
	NFTs       []NFT       `mapstructure:"nfts" json:"nfts"`             // list of boosting NFTs
	Tokens     []Token     `mapstructure:"tokens" json:"tokens"`         // list of boosting tokens
}
type Milestone struct {
	Minimum int `mapstructure:"minimum" json:"minimum"` // minimum amount of vulti to reach this milestone
	Prize   int `mapstructure:"prize" json:"prize"`     // prize for this milestone
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
	viper.SetDefault("vultiref.api_key", "")
	viper.SetDefault("vultiref.base_address", "")
	viper.SetDefault("season.swap_multiplier", 1.6)
	viper.SetDefault("season.referral_multiplier", 1.5)
	viper.SetDefault("season.milestones", []int{5000, 10000, 50000, 100000})
	viper.SetDefault("season.nfts", []NFT{})
	viper.SetDefault("season.tokens", []Token{})

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		if errors.As(err, &configFileNotFoundError) {
			log.Println("No config file found, using environment variables and defaults")
		}
	}

	err := viper.Unmarshal(&cfg, func(dc *mapstructure.DecoderConfig) {
		dc.DecodeHook = mapstructure.StringToTimeHookFunc(time.RFC3339)
	})
	if err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %w", err)
	}
	return &cfg, nil
}

func (cfg *Config) GetCurrentSeason() AirdropSeason {
	var currentSeason AirdropSeason
	for _, season := range cfg.Seasons {
		if time.Now().After(season.Start) && time.Now().Before(season.End) {
			currentSeason = season

		}
	}
	return currentSeason
}
