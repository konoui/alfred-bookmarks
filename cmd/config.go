package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/konoui/alfred-bookmarks/pkg/bookmarker"
	"github.com/spf13/viper"
)

const (
	firefoxDefaultProfileName = "default"
	chromeDefaultProfileName  = "default"
)

var (
	firefoxDefaultProfilePath = os.ExpandEnv("${HOME}/Library/Application Support/Firefox/Profiles")
	chromeDefaultProfilePath  = os.ExpandEnv("${HOME}/Library/Application Support/Google/Chrome")
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
	Enable      bool   `mapstructure:"enable"`
	ProfileName string `mapstructure:"profile_name,omitempty"`
	ProfilePath string `mapstructure:"profile_path,omitempty"`
}

// Chrome Configuration
type Chrome struct {
	Enable      bool   `mapstructure:"enable"`
	ProfileName string `mapstructure:"profile_name,omitempty"`
	ProfilePath string `mapstructure:"profile_path,omitempty"`
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
	viper.SetDefault("firefox.profile_name", firefoxDefaultProfileName)
	viper.SetDefault("firefox.profile_path", firefoxDefaultProfilePath)
	viper.SetDefault("chrome.profile_name", chromeDefaultProfileName)
	viper.SetDefault("chrome.profile_path", chromeDefaultProfilePath)
	defer c.resolvePath()
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
	_, firefoxErr := bookmarker.GetFirefoxBookmarkFile(firefoxDefaultProfilePath, firefoxDefaultProfileName)
	_, chromeErr := bookmarker.GetChromeBookmarkFile(chromeDefaultProfilePath, chromeDefaultProfileName)
	_, safariErr := bookmarker.GetSafariBookmarkFile()
	if firefoxErr != nil && chromeErr != nil && safariErr != nil {
		return c, errors.New("found no available bookmarks on your computer")
	}

	for _, err := range []error{firefoxErr, chromeErr, safariErr} {
		if err != nil {
			awf.Logf(err.Error() + "\n")
		}
	}

	c.RemoveDuplicate = true
	if firefoxErr == nil {
		c.Firefox.Enable = true
		c.Firefox.ProfileName = firefoxDefaultProfileName
		c.Firefox.ProfilePath = firefoxDefaultProfilePath
	}
	if chromeErr == nil {
		c.Chrome.Enable = true
		c.Chrome.ProfileName = chromeDefaultProfileName
		c.Chrome.ProfilePath = chromeDefaultProfilePath
	}
	if safariErr == nil {
		c.Safari.Enable = true
	}

	return c, nil
}

func (c *Config) resolvePath() {
	c.Firefox.ProfilePath = os.ExpandEnv(c.Firefox.ProfilePath)
	c.Chrome.ProfilePath = os.ExpandEnv(c.Chrome.ProfilePath)
}

func convertDefaultTTL(hour int) time.Duration {
	if hour == 0 {
		hour = 24
	} else if hour < 0 {
		hour = 0
	}
	return time.Duration(hour) * time.Hour
}
