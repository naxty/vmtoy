package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	workspace *Workspace
	settings  map[string]interface{}
}

func NewConfig() (*Config, error) {
	ws, err := NewWorkspace()
	if err != nil {
		return nil, err
	}

	if err := ws.EnsureExists(); err != nil {
		return nil, err
	}

	cfg := &Config{
		workspace: ws,
		settings:  make(map[string]interface{}),
	}

	// Load existing config if present
	data, err := os.ReadFile(ws.ConfigPath)
	if err == nil {
		json.Unmarshal(data, &cfg.settings)
	}

	return cfg, nil
}

func (c *Config) GetWorkspace() *Workspace {
	return c.workspace
}

func (c *Config) Save() error {
	data, err := json.MarshalIndent(c.settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(c.workspace.ConfigPath, data, 0644)
}
