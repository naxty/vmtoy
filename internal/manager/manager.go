package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/naxty/vmtoy/internal/config"
)

const vmPrefix = "vm-"

type VMMetadata struct {
	Name      string    `json:"name"`
	Created   time.Time `json:"created"`
	LastUsed  time.Time `json:"last_used"`
	ImagePath string    `json:"image_path"`
	Status    string    `json:"status"`
}

type VM struct {
	Name     string
	Path     string
	Metadata VMMetadata
}

type Manager struct {
	config    *config.Config
	workspace *config.Workspace
}

func NewManager() (*Manager, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to create manager: %w", err)
	}

	return &Manager{
		config:    cfg,
		workspace: cfg.GetWorkspace(),
	}, nil
}

func (m *Manager) Create(name string, imagePath string) error {
	if !strings.HasPrefix(name, vmPrefix) {
		name = vmPrefix + name
	}

	vmPath := filepath.Join(m.workspace.VMsDir, name)
	if _, err := os.Stat(vmPath); !os.IsNotExist(err) {
		return fmt.Errorf("VM %s already exists", name)
	}

	// Create VM directory
	if err := os.MkdirAll(vmPath, 0755); err != nil {
		return fmt.Errorf("failed to create VM directory: %w", err)
	}

	// Create metadata
	metadata := VMMetadata{
		Name:      name,
		Created:   time.Now(),
		LastUsed:  time.Now(),
		ImagePath: imagePath,
		Status:    "created",
	}

	// Save metadata
	return m.saveVMMetadata(vmPath, metadata)
}

func (m *Manager) List() ([]VM, error) {
	entries, err := os.ReadDir(m.workspace.VMsDir)
	if err != nil {
		return nil, fmt.Errorf("error reading VM directory: %w", err)
	}

	vms := []VM{}
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), vmPrefix) {
			vm, err := m.loadVM(entry.Name())
			if err != nil {
				continue
			}
			vms = append(vms, vm)
		}
	}
	return vms, nil
}

func (m *Manager) Delete(name string) error {
	if !strings.HasPrefix(name, vmPrefix) {
		name = vmPrefix + name
	}

	vmPath := filepath.Join(m.workspace.VMsDir, name)
	if _, err := os.Stat(vmPath); os.IsNotExist(err) {
		return fmt.Errorf("VM %s does not exist", name)
	}

	return os.RemoveAll(vmPath)
}

func (m *Manager) Exists(name string) bool {
	if !strings.HasPrefix(name, vmPrefix) {
		name = vmPrefix + name
	}

	vmPath := filepath.Join(m.workspace.VMsDir, name)
	_, err := os.Stat(vmPath)
	return err == nil
}

func (m *Manager) loadVM(name string) (VM, error) {
	vmPath := filepath.Join(m.workspace.VMsDir, name)
	metadata, err := m.loadVMMetadata(vmPath)
	if err != nil {
		return VM{}, err
	}

	return VM{
		Name:     name,
		Path:     vmPath,
		Metadata: metadata,
	}, nil
}

func (m *Manager) saveVMMetadata(vmPath string, metadata VMMetadata) error {
	data, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	metadataPath := filepath.Join(vmPath, "metadata.json")
	return os.WriteFile(metadataPath, data, 0644)
}

func (m *Manager) loadVMMetadata(vmPath string) (VMMetadata, error) {
	var metadata VMMetadata
	data, err := os.ReadFile(filepath.Join(vmPath, "metadata.json"))
	if err != nil {
		return metadata, fmt.Errorf("failed to read metadata: %w", err)
	}

	if err := json.Unmarshal(data, &metadata); err != nil {
		return metadata, fmt.Errorf("failed to unmarshal metadata: %w", err)
	}

	return metadata, nil
}

func (m *Manager) LoadVM(name string) (*VM, error) {
	if !strings.HasPrefix(name, vmPrefix) {
		name = vmPrefix + name
	}

	vm, err := m.loadVM(name)
	if err != nil {
		return nil, fmt.Errorf("failed to load VM %s: %w", name, err)
	}

	return &vm, nil
}

func (m *Manager) GetWorkspace() *config.Workspace {
	return m.workspace
}
