package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Site      string `json:"site"`
	Session   string `json:"session"`
	CSRFToken string `json:"csrfToken"`
}

type fileConfig struct {
	Site      string `json:"site"`
	Session   string `json:"session"`
	CSRFToken string `json:"csrfToken"`
}

func Load(workspace string) (Config, error) {
	cfg := Config{
		Site: "https://leetcode.com",
	}

	fileCfg, err := loadFromFile(filepath.Join(workspace, ".leetcode.json"))
	if err != nil {
		return Config{}, err
	}
	if fileCfg != nil {
		if fileCfg.Site != "" {
			cfg.Site = fileCfg.Site
		}
		cfg.Session = fileCfg.Session
		cfg.CSRFToken = fileCfg.CSRFToken
	}

	if v := os.Getenv("LC_SITE"); v != "" {
		cfg.Site = v
	}
	if v := os.Getenv("LEETCODE_SESSION"); v != "" {
		cfg.Session = v
	}
	if v := os.Getenv("LEETCODE_CSRF_TOKEN"); v != "" {
		cfg.CSRFToken = v
	}

	cfg.Site = strings.TrimSpace(strings.TrimRight(cfg.Site, "/"))
	if cfg.Site == "" {
		return Config{}, errors.New("site is empty")
	}

	return cfg, nil
}

func (c Config) ValidateForSubmit() error {
	if c.Session == "" || c.CSRFToken == "" {
		return errors.New("missing auth, set LEETCODE_SESSION and LEETCODE_CSRF_TOKEN (or .leetcode.json)")
	}
	return nil
}

func loadFromFile(path string) (*fileConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}

	var c fileConfig
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}
