package cache

import (
	"os"
	"testing"
	"time"
)

type example struct {
	A string
	B string
	C []string
}

var storedValue = example{
	A: "AAAAA",
	B: "BBBBBB",
	C: []string{
		"11111",
		"22222",
		"33333",
	},
}

func TestNewCache(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
	}{
		{description: "exists cache dir", dir: os.TempDir(), file: "test1", expiredTime: 3 * time.Minute, expectErr: false},
		{description: "no exists cache dir", dir: "/unk", file: "test2", expiredTime: 0 * time.Minute, expectErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			_, err := New(tt.dir, tt.file, tt.expiredTime)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error %s", err.Error())
			}
		})
	}
}

func TestStore(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
	}{
		{description: "create cache file on temp dir", dir: os.TempDir(), file: "test1", expiredTime: 3 * time.Minute, expectErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cache, err := New(tt.dir, tt.file, tt.expiredTime)
			if err != nil {
				t.Errorf("cannot create cache instance. error %+v", err)
			}

			// remove cache file before test
			if err = cache.Clear(); err != nil {
				t.Errorf("unexpected error %+v", err)
			}

			err = cache.Store(&storedValue)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error %s", err.Error())
			}
		})
	}
}

func TestLoad(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
	}{
		{description: "create cache file on temp dir", dir: os.TempDir(), file: "test1", expiredTime: 3 * time.Minute, expectErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cache, err := New(tt.dir, tt.file, tt.expiredTime)
			if err != nil {
				t.Errorf("cannot create cache instance. error %+v", err)
			}

			// remove cache file before test
			if err = cache.Clear(); err != nil {
				t.Errorf("unexpected error %+v", err)
			}

			err = cache.Store(&storedValue)
			if err != nil {
				t.Errorf("cannot store data into cache. error %+v", err)
			}

			loadedValue := example{}
			err = cache.Load(&loadedValue)
			if tt.expectErr && err == nil {
				t.Errorf("expect error happens, but got response")
			}

			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error %s", err.Error())
			}
		})
	}
}

func TestExpired(t *testing.T) {
	tests := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
	}{
		{description: "create cache file on temp dir", dir: os.TempDir(), file: "test1", expiredTime: 3 * time.Minute, expectErr: false},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			cacher, err := New(tt.dir, tt.file, tt.expiredTime)
			if err != nil {
				t.Errorf("cannot create cache instance. error %+v", err)
			}
			cache := cacher.(*Cache)

			err = cache.Store(&storedValue)
			if err != nil {
				t.Errorf("cannot store data into cache. error %+v", err)
			}

			if age, err := cache.Age(); err != nil && 0 <= age {
				t.Errorf("unexpected cache expired or error %+v", err)
			}

			if cache.Expired() && !cache.NotExpired() {
				t.Errorf("unexpected cache expired")
			}
		})
	}
}
