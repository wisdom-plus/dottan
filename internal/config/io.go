package config

import (
	"errors"
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

var ErrAlreayExists = errors.New("config file already exists")

func ConfigPath(appName string) (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil (
		return "", err
	)
	return filepath.Join(dir, appName, "config.toml"), nil
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Save(path string, cfg *Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	b, err := toml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o600)
}
