package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/naxty/vmtoy/internal/config"
)

const vmPrefix = "vm-"

type VM struct {
	Name string
	Path string
}

type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) List() {
	// list all VMs
	entries, err := os.ReadDir(config.VMTOY_WORKDIR)
	if err != nil {
		fmt.Errorf("error reading directory: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), vmPrefix) {
			fmt.Println(entry.Name())
		}
	}
}

func (m *Manager) Exists(vmName string) bool {
	// Ensure the vmName starts with the proper prefix.
	if !strings.HasPrefix(vmName, vmPrefix) {
		vmName = vmPrefix + vmName
	}

	vmPath := filepath.Join(config.VMTOY_WORKDIR, vmName)
	info, err := os.Stat(vmPath)
	if err != nil {
		return false
	}
	return info.IsDir()
}
