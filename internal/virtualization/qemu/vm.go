package qemu

import (
	"fmt"
	"strings"
)

type PortForward struct {
	HostPort  int
	GuestPort int
	Protocol  string // tcp or udp
}

type QemuVMConfig struct {
	Name         string
	Memory       int
	CPUs         int
	DiskImage    string
	ISO          string
	Accelerator  string
	NetworkType  string
	GraphicMode  bool
	PortForwards []PortForward
	Drive        string
	State        string
	Distro       string
}

func NewQemuVMConfig(name string) *QemuVMConfig {
	return &QemuVMConfig{
		Name:        name,
		Memory:      1024,
		CPUs:        1,
		Accelerator: "tcg",
		NetworkType: "user",
		State:       "new",
	}
}

func (vm *QemuVMConfig) buildPortForwards() string {
	if len(vm.PortForwards) == 0 {
		return ""
	}

	var forwards []string
	for _, pf := range vm.PortForwards {
		forwards = append(forwards, fmt.Sprintf(",hostfwd=%s::%d-:%d",
			pf.Protocol, pf.HostPort, pf.GuestPort))
	}
	return strings.Join(forwards, "")
}

func (c *QemuVMConfig) BuildCommand() []string {
	args := []string{
		"-name", c.Name,
		"-m", fmt.Sprintf("%d", c.Memory),
		"-smp", fmt.Sprintf("%d", c.CPUs),
		"-accel", c.Accelerator,
		//"-serial", "telnet:localhost:4321,server,nowait",
		"-monitor", "tcp:localhost:4320,server,nowait",
		"-serial", "tcp:localhost:4321,server,nowait",
	}

	if c.DiskImage != "" {
		args = append(args, "-hda", c.DiskImage)
	}

	if c.ISO != "" {
		args = append(args, "-cdrom", c.ISO, "-boot", "d")
	}

	if len(c.PortForwards) > 0 {
		netdev := fmt.Sprintf("user,id=net0%s", c.buildPortForwards())
		args = append(args,
			"-netdev", netdev,
			"-device", "virtio-net-pci,netdev=net0")
	}

	if !c.GraphicMode {
		args = append(args, "-nographic")
	}

	return args
}

func (c *QemuVMConfig) GetName() string  { return c.Name }
func (c *QemuVMConfig) GetState() string { return c.State }
