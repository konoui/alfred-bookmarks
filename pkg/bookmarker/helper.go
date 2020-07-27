package bookmarker

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

func validateURL(s string) (u *url.URL, err error) {
	u, err = url.Parse(s)
	// Ignore invalid URLs
	if err != nil {
		return
	}
	if u.Host == "" {
		return u, errors.New("hostname is enpty")
	}
	return
}

// getLatestFile returns a path to latest files in dir
func getLatestFile(dir string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}
	latestIndex := 0
	for i, file := range files {
		if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
			continue
		}
		if time.Since(file.ModTime()) <= time.Since(files[latestIndex].ModTime()) {
			latestIndex = i
		}
	}

	return filepath.Join(dir, files[latestIndex].Name()), nil
}

// searchSuffixDir returns a directory name of suffix ignoring case-sensitive
func searchSuffixDir(dir, suffux string) (string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if name := file.Name(); file.IsDir() &&
			strings.HasSuffix(strings.ToLower(name), strings.ToLower(suffux)) {
			return name, nil
		}
	}

	return "", fmt.Errorf("not found a directory of suffix (%s) in %s directory", suffux, dir)
}
