package qemu

import (
	"fmt"

	"github.com/naxty/vmtoy/internal/distros/alpine"
	distros "github.com/naxty/vmtoy/internal/distros/common"
)

type QemuWrapper struct {
	config  *QemuVMConfig
	process *QemuProcess
}

func New(config *QemuVMConfig) *QemuWrapper {
	return &QemuWrapper{
		config: config,
	}
}

func (q *QemuWrapper) Start() error {
	if q.process != nil {
		return fmt.Errorf("VM is already running")
	}

	process, err := NewQemuProcess(q.config.BuildCommand())
	if err != nil {
		return fmt.Errorf("failed to create QEMU process: %w", err)
	}

	if err := process.Start(); err != nil {
		return fmt.Errorf("failed to start QEMU: %w", err)
	}

	q.process = process
	q.config.State = "running"
	return nil
}

func (q *QemuWrapper) Stop() error {
	if q.process == nil {
		return nil
	}

	if err := q.process.Stop(); err != nil {
		return fmt.Errorf("failed to stop VM: %w", err)
	}

	q.process = nil
	q.config.State = "stopped"
	return nil
}

func (q *QemuWrapper) Install() error {
	// Create installer based on distribution
	var err error
	var installer distros.Installer
	switch q.config.Distro {
	case "alpine":
		installer, err = alpine.NewAlpineSetup("localhost", 4321)
	default:
		return fmt.Errorf("unsupported distribution: %s", q.config.Distro)
	}

	if err != nil {
		return fmt.Errorf("failed to create installer: %w", err)
	}
	defer installer.Close()

	if err := installer.Install(); err != nil {
		return fmt.Errorf("installation failed: %w", err)
	}

	q.config.State = "installed"
	return nil
}

func (q *QemuWrapper) Status() string {
	return q.config.GetState()
}
