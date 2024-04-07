package main

import (
	"os"
	"path/filepath"
	"encoding/json"

	"golang.org/x/oauth2"
	log "github.com/sirupsen/logrus"
)

// Get the default cache directory.
func default_cache_dir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.WithError(err).Warn("failed to identify a home directory")
		return ""
	}

	return filepath.Join(home, ".local", "nspotify")
}

// Try to cache an access token.
func WriteCache(dir string, tok *oauth2.Token) error {
	err := os.Mkdir(dir, 0700)
	if err != nil && !os.IsExist(err) {
		log.WithError(err).Warnf("failed to make cache directory: %s", dir)
		return err
	}

	data, err := json.Marshal(tok)
	if err != nil {
		log.WithError(err).Warn("failed to marshall token")
		return err
	}

	err = os.WriteFile(filepath.Join(dir, "token.json"), data, 0600)
	if err != nil {
		log.WithError(err).Warn("failed to write cache file")
		return err
	}

	log.Debug("succeeded in caching token")

	return nil
}

// Try to read a cache file for an access token.
func ReadCache(dir string) (*oauth2.Token, error) {
	tok := &oauth2.Token{}

	full_path := filepath.Join(dir, "token.json")

	data, err := os.ReadFile(full_path)
	if err != nil {
		log.WithError(err).Warnf("failed to read cache file: %s", full_path)
		return nil, err
	}

	log.Debugf("found cache file: %s", full_path)

	err = json.Unmarshal(data, tok)
	if err != nil {
		log.WithError(err).Warn("failed to unmarshall token")
		return nil, err
	}

	log.Debug("succeeded in reading cached token")

	return tok, nil
}

