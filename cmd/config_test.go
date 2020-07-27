package cmd

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

// testConfig is the same value as .alfred-bookmarks
var testConfig = &Config{
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
}

func TestNewConfig(t *testing.T) {
	tests := []struct {
		description string
		want        *Config
	}{
		{
			description: "read config file except for firefox profile. firefox profile should be default value",
			want:        testConfig,
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

func Test_convertDefaultTTL(t *testing.T) {
	type args struct {
		hour int
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "return zero if pass misus value",
			args: args{
				hour: -1,
			},
			want: 0,
		},
		{
			name: "return 24 if pass zero value",
			args: args{
				hour: 0,
			},
			want: 24 * time.Hour,
		},
		{
			name: "return the value if pass non zero or non minus value",
			args: args{
				hour: 5,
			},
			want: 5 * time.Hour,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertDefaultTTL(tt.args.hour); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertDefaultTTL() = %v, want %v", got, tt.want)
			}
		})
	}
}
