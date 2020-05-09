package cmd

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		description string
		want        *Config
	}{
		{
			description: "read config file except for firefox profile, that is default value",
			want: &Config{
				RemoveDuplicate: true,
				// disable cache
				MaxCacheAge: -1,
				Firefox: Firefox{
					Enable:  true,
					Profile: "default",
				},
				Chrome: Chrome{
					Enable:  true,
					Profile: "Default",
				},
			},
		},
	}
	for _, tt := range tests {
		c, err := newConfig()
		if err != nil {
			t.Fatal(err)
		}

		if !cmp.Equal(c, tt.want) {
			t.Errorf("want: \n%+v, got: \n%+v", tt.want, c)
		}
	}
}
