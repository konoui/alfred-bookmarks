package cache

import (
	"reflect"
	"testing"
)

func TestNewNilCache(t *testing.T) {
	tests := []struct {
		description string
		want        *NilCache
	}{
		{
			description: "exists cache dir",
			want:        &NilCache{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			got := NewNilCache()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %+v\ngot: %+v", tt.want, got)
			}
		})
	}
}

func TestNilCache(t *testing.T) {
	tests := []struct {
		description   string
		expiredResult bool
		loadResult    error
		storeResult   error
	}{
		{
			description:   "exists cache dir",
			expiredResult: true,
			loadResult:    nil,
			storeResult:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			c := NewNilCache()
			if got := c.Expired(); got != tt.expiredResult {
				t.Errorf("want: %+v\ngot: %+v", tt.expiredResult, got)
			}

			if got := c.Load(tt.storeResult); got != tt.loadResult {
				t.Errorf("want: %+v\ngot: %+v", tt.loadResult, got)
			}

			if got := c.Store(tt.storeResult); got != tt.storeResult {
				t.Errorf("want: %+v\ngot: %+v", tt.storeResult, got)
			}
		})
	}
}
