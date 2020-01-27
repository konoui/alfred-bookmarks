package cmd

import (
	"github.com/spf13/viper"
)

// Config configuration which browser bookmark read
type Config struct {
	Firefox         Firefox `mapstructure:"firefox"`
	Chrome          Chrome  `mapstructure:"chrome"`
	RemoveDuplicate bool    `mapstructure:"remove_duplicate"`
	MaxCacheAge     int     `mapstructure:"cache_age_hours"`
}

// Firefox Configuration
type Firefox struct {
	Enable  bool   `mapstructure:"enable"`
	Profile string `mapstructure:"profile,omitempty"`
}

// Chrome Configuration
type Chrome struct {
	Enable  bool   `mapstructure:"enable"`
	Profile string `mapstructure:"profile,omitempty"`
}

// NewConfig return alfred bookmark configuration
func newConfig() (*Config, error) {
	var c Config
	viper.SetConfigType("yaml")
	viper.SetConfigName(".alfred-bookmarks")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/")

	// Set Default Value overwritten with config file
	viper.SetDefault("firefox.profile", "default")
	viper.SetDefault("chrome.profile", "default")
	if err := viper.ReadInConfig(); err != nil {
		return &Config{}, err
	}

	if err := viper.Unmarshal(&c); err != nil {
		return &Config{}, err
	}

	return &c, nil
}
