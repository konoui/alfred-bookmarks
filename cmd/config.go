package cmd

import (
	"time"

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
func newConfig() (c *Config, err error) {
	c = new(Config)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".alfred-bookmarks")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/")

	// Set default value overwritten with config file
	viper.SetDefault("firefox.profile", "default")
	viper.SetDefault("chrome.profile", "default")
	if err = viper.ReadInConfig(); err != nil {
		return
	}

	if err = viper.Unmarshal(c); err != nil {
		return
	}
	return
}

func convertDefaultTTL(hour int) time.Duration {
	if hour == 0 {
		hour = 24
	} else if hour < 0 {
		hour = 0
	}
	return time.Duration(hour) * time.Hour
}
