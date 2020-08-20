package cmd

import (
	"errors"
	"time"

	"github.com/konoui/alfred-bookmarks/pkg/bookmarker"
	"github.com/spf13/viper"
)

const (
	firefoxDefaultProfile = "default"
	chromeDefaultProfile  = "default"
)

// Config configuration which browser bookmark read
type Config struct {
	Firefox         Firefox `mapstructure:"firefox"`
	Chrome          Chrome  `mapstructure:"chrome"`
	Safari          Safari  `mapstructure:"safari"`
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

// Safari Configuration
type Safari struct {
	Enable bool `mapstructure:"enable"`
}

// NewConfig return alfred bookmark configuration
func newConfig() (c *Config, err error) {
	c = new(Config)
	viper.SetConfigType("yaml")
	viper.SetConfigName(".alfred-bookmarks")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/")

	// Set default value overwritten with config file
	viper.SetDefault("firefox.profile", firefoxDefaultProfile)
	viper.SetDefault("chrome.profile", chromeDefaultProfile)
	if err = viper.ReadInConfig(); err != nil {
		// Try to continue using available bookmarks if config file does not exist
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return
		}
		return availableConfig()
	}

	if err = viper.Unmarshal(c); err != nil {
		return
	}
	return
}

func availableConfig() (*Config, error) {
	c := new(Config)
	_, firefoxErr := bookmarker.GetFirefoxBookmarkFile(firefoxDefaultProfile)
	_, chromeErr := bookmarker.GetChromeBookmarkFile(chromeDefaultProfile)
	_, safariErr := bookmarker.GetSafariBookmarkFile()
	if firefoxErr != nil && chromeErr != nil && safariErr != nil {
		return c, errors.New("found no available bookmarks on your computer")
	}

	c.RemoveDuplicate = true
	if firefoxErr == nil {
		c.Firefox.Enable = true
		c.Firefox.Profile = firefoxDefaultProfile
	}
	if chromeErr == nil {
		c.Chrome.Enable = true
		c.Chrome.Profile = chromeDefaultProfile
	}
	if safariErr == nil {
		c.Safari.Enable = true
	}

	return c, nil
}

func convertDefaultTTL(hour int) time.Duration {
	if hour == 0 {
		hour = 24
	} else if hour < 0 {
		hour = 0
	}
	return time.Duration(hour) * time.Hour
}
