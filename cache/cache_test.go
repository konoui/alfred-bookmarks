package cache

import (
	"os"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cases := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
	}{
		{description: "exists cache dir", dir: os.TempDir(), file: "test1", expiredTime: 3 * time.Minute, expectErr: false},
		{description: "no exists cache dir", dir: "/unk", file: "test2", expiredTime: 0 * time.Minute, expectErr: true},
	}

	for _, c := range cases {
		_, err := NewCache(c.dir, c.file, c.expiredTime)
		if c.expectErr && err == nil {
			t.Errorf("%s: expect error happens, but got response", c.description)
		}

		if !c.expectErr && err != nil {
			t.Errorf("%s: unexpected error %s", c.description, err.Error())
		}
	}
}

func TestStore(t *testing.T) {
	cases := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
	}{
		{description: "create cache file on temp dir", dir: os.TempDir(), file: "test1", expiredTime: 3 * time.Minute, expectErr: false},
	}

	for _, c := range cases {
		type example struct {
			A string
			B string
			C []string
		}
		storedValue := example{
			A: "AAAAA",
			B: "BBBBBB",
			C: []string{
				"11111",
				"22222",
				"33333",
			},
		}
		cache, err := NewCache(c.dir, c.file, c.expiredTime)
		if err != nil {
			t.Errorf("cannot create cache instance. error %+v", err)
		}

		// remove cache file before test
		if err = cache.Clear(); err != nil {
			t.Errorf("unexpected error %+v", err)
		}

		err = cache.Store(&storedValue)
		if c.expectErr && err == nil {
			t.Errorf("%s: expect error happens, but got response", c.description)
		}

		if !c.expectErr && err != nil {
			t.Errorf("%s: unexpected error %s", c.description, err.Error())
		}
	}
}

func TestLoad(t *testing.T) {
	cases := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
	}{
		{description: "create cache file on temp dir", dir: os.TempDir(), file: "test1", expiredTime: 3 * time.Minute, expectErr: false},
	}

	for _, c := range cases {
		type example struct {
			A string
			B string
			C []string
		}
		storedValue := example{
			A: "AAAAA",
			B: "BBBBBB",
			C: []string{
				"11111",
				"22222",
				"33333",
			},
		}
		cache, err := NewCache(c.dir, c.file, c.expiredTime)
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
		if c.expectErr && err == nil {
			t.Errorf("%s: expect error happens, but got response", c.description)
		}

		if !c.expectErr && err != nil {
			t.Errorf("%s: unexpected error %s", c.description, err.Error())
		}
	}
}

func TestExpired(t *testing.T) {
	cases := []struct {
		description string
		dir         string
		file        string
		expiredTime time.Duration
		expectErr   bool
	}{
		{description: "create cache file on temp dir", dir: os.TempDir(), file: "test1", expiredTime: 3 * time.Minute, expectErr: false},
	}

	for _, c := range cases {
		type example struct {
			A string
			B string
			C []string
		}
		storedValue := example{
			A: "AAAAA",
			B: "BBBBBB",
			C: []string{
				"11111",
				"22222",
				"33333",
			},
		}
		cache, err := NewCache(c.dir, c.file, c.expiredTime)
		if err != nil {
			t.Errorf("cannot create cache instance. error %+v", err)
		}

		err = cache.Store(&storedValue)
		if err != nil {
			t.Errorf("cannot store data into cache. error %+v", err)
		}

		if age, err := cache.Age(); err != nil && 0 <= age {
			t.Errorf("unexpected cache expired or error %+v", err)
		}

		if cache.Expired(c.expiredTime) && !cache.NotExpired(c.expiredTime) {
			t.Errorf("unexpected cache expired")
		}
	}
}
