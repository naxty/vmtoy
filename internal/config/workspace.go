package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// Directory structure
	WorkspaceDirName = ".vmtoy"
	ImagesDirName    = "images"
	VMsDirName       = "vms"
	CacheDirName     = "cache"
	ConfigFileName   = "config.json"
)

type Workspace struct {
	RootDir    string
	ImagesDir  string
	VMsDir     string
	CacheDir   string
	ConfigPath string
}

func NewWorkspace() (*Workspace, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	ws := &Workspace{
		RootDir:    filepath.Join(homeDir, WorkspaceDirName),
		ImagesDir:  filepath.Join(homeDir, WorkspaceDirName, ImagesDirName),
		VMsDir:     filepath.Join(homeDir, WorkspaceDirName, VMsDirName),
		CacheDir:   filepath.Join(homeDir, WorkspaceDirName, CacheDirName),
		ConfigPath: filepath.Join(homeDir, WorkspaceDirName, ConfigFileName),
	}

	return ws, nil
}

func (w *Workspace) Initialize() error {
	dirs := []string{
		w.RootDir,
		w.ImagesDir,
		w.VMsDir,
		w.CacheDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create default config file if it doesn't exist
	if _, err := os.Stat(w.ConfigPath); os.IsNotExist(err) {
		if err := os.WriteFile(w.ConfigPath, []byte("{}"), 0644); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
	}

	return nil
}

func (w *Workspace) EnsureExists() error {
	if _, err := os.Stat(w.RootDir); os.IsNotExist(err) {
		return w.Initialize()
	}
	return nil
}

func (w *Workspace) CleanWorkspace() error {
	return os.RemoveAll(w.RootDir)
}
