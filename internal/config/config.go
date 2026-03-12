package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Profile struct {
	URL   string `json:"url"`
	Token string `json:"token,omitempty"`
	Email string `json:"email,omitempty"`
}

type Config struct {
	Profiles      map[string]*Profile `json:"profiles"`
	ActiveProfile string              `json:"active_profile"`
}

func defaultConfig() *Config {
	return &Config{
		Profiles: map[string]*Profile{
			"default": {URL: "http://localhost:3000"},
		},
		ActiveProfile: "default",
	}
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "metabase", "config.json")
}

func Load() (*Config, error) {
	data, err := os.ReadFile(configPath())
	if os.IsNotExist(err) {
		return defaultConfig(), nil
	}
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if cfg.Profiles == nil {
		cfg.Profiles = map[string]*Profile{"default": {URL: "http://localhost:3000"}}
	}
	return &cfg, nil
}

func (c *Config) Save() error {
	path := configPath()
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func (c *Config) GetProfile(name string) *Profile {
	if p, ok := c.Profiles[name]; ok {
		return p
	}
	if p, ok := c.Profiles["default"]; ok {
		return p
	}
	p := &Profile{URL: "http://localhost:3000"}
	c.Profiles[name] = p
	return p
}
