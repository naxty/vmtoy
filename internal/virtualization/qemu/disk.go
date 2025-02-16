package qemu

import "os/exec"

// CreateDisk creates a new QCOW2 disk image
func CreateDisk(path string, size string) error {
	args := []string{"create", "-f", "qcow2", path, size}
	cmd := exec.Command("qemu-img", args...)
	return cmd.Run()
}
