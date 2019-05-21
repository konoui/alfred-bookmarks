package cache

import (
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

// Cache implements a simple store/load API, saving data to specified directory.
type Cache struct {
	Dir         string
	File        string
	ExpiredTime time.Duration
}

// NewCache creates a new cache Instance
func NewCache(dir, file string, expiredTime time.Duration) (*Cache, error) {
	if !pathExists(dir) {
		return &Cache{}, fmt.Errorf("%s directory does not exist", dir)
	}
	return &Cache{
		Dir:         dir,
		File:        file,
		ExpiredTime: expiredTime,
	}, nil
}

// Store save data into cache
func (c Cache) Store(v interface{}) error {
	p := c.path()
	f, err := os.OpenFile(p, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = gob.NewEncoder(f).Encode(v); err != nil {
		log.Printf("cannot save data into cache (%s). error %+v", p, err)
		return err
	}
	log.Printf("saving data into cache (%s) is success", p)
	return nil
}

// Load read data saved cache into v
func (c Cache) Load(v interface{}) error {
	p := c.path()
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = gob.NewDecoder(f).Decode(v); err != nil {
		log.Printf("cannot read data from cache (%s). error %+v", p, err)
		return err
	}
	log.Printf("reading data from cache %s is success", p)
	return nil
}

// Clear remove cache file if exists
func (c Cache) Clear() error {
	p := c.path()
	if pathExists(p) {
		return os.Remove(p)
	}
	return nil
}

// NotExpired return true if cache is no expired
func (c Cache) NotExpired(maxAge time.Duration) bool {
	return !c.Expired(maxAge)
}

// Expired return true if cache is expired
func (c Cache) Expired(maxAge time.Duration) bool {
	age, err := c.Age()
	if err != nil {
		return true
	}
	return age > maxAge
}

// Age return the time since the data is cached at
func (c Cache) Age() (time.Duration, error) {
	p := c.path()
	fi, err := os.Stat(p)
	if err != nil {
		return time.Duration(0), err
	}
	return time.Since(fi.ModTime()), nil
}

// Exists return true if the cache file exists
func (c Cache) Exists() bool {
	return pathExists(c.path())
}

// path return the path of cache file
func (c Cache) path() string {
	return filepath.Join(c.Dir, c.File)
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
