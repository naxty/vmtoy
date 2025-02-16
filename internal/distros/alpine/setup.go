package alpine

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/naxty/vmtoy/internal/distros/common"
)

// AlpineSetup handles the automated installation of Alpine Linux
type AlpineSetup struct {
	conn    common.Connection
	reader  *bufio.Reader
	timeout time.Duration
}

// InstallStep represents a single step in the installation process
type InstallStep struct {
	Command string
	Trigger string
	Timeout time.Duration
}

// NewAlpineSetup creates a new Alpine Linux setup handler
func NewAlpineSetup(host string, port int) (*AlpineSetup, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to VM: %w", err)
	}

	return &AlpineSetup{
		conn:    conn,
		reader:  bufio.NewReader(conn),
		timeout: 30 * time.Second,
	}, nil
}

func (a *AlpineSetup) Close() error {
	return a.conn.Close()
}

func (a *AlpineSetup) send(command string) error {
	fmt.Printf("\n[SENDING] %q\n", command)
	_, err := a.conn.Write([]byte(command + "\n"))
	return err
}

func (a *AlpineSetup) waitForString(trigger string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var buffer strings.Builder

	for time.Now().Before(deadline) {
		b := make([]byte, 1)
		_, err := a.conn.Read(b)
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("read error: %w", err)
			}
			return nil
		}

		buffer.Write(b)
		fmt.Print(string(b))

		if strings.Contains(buffer.String(), trigger) {
			return nil
		}
	}
	return fmt.Errorf("timeout waiting for: %s", trigger)
}

func (a *AlpineSetup) sendAndWait(command, trigger string, timeout time.Duration) error {
	if err := a.send(command); err != nil {
		return fmt.Errorf("send failed: %w", err)
	}
	if trigger == "" {
		time.Sleep(1 * time.Second)
		return nil
	}
	return a.waitForString(trigger, timeout)
}

// defaultInstallSteps returns the default Alpine Linux installation steps
func defaultInstallSteps() []InstallStep {
	return []InstallStep{
		{"", "login:", 10 * time.Second},     // Wait for initial login prompt
		{"root", "Welcome", 3 * time.Second}, // Login
		//{"", "Password:", 10 * time.Second},           // Wait for password prompt
		{"echo 'new-hostname' > /etc/hostname", "localhost", 1 * time.Second},
		{"setup-interfaces", "Available interfaces are:", 5 * time.Second},
		{"eth0", "Ip address for eth0?", 5 * time.Second},
		{"dhcp", "manual network configuration?", 5 * time.Second},
		{"n", "localhost:~#", 5 * time.Second},
		// Start networking
		{"rc-service networking --quiet start", "localhost:~#", 10 * time.Second},

		// Setup timezone
		{"setup-timezone -z UTC", "localhost:~#", 5 * time.Second},
		// Setup repositories

		//{"setup-apkrepos -1", "localhost:~#", 30 * time.Second},
		{"echo 'http://dl-cdn.alpinelinux.org/alpine/v3.21/main' >> /etc/apk/repositories", "localhost:~#", 5 * time.Second},
		{"echo '#http://dl-cdn.alpinelinux.org/alpine/v3.21/community' >> /etc/apk/repositories", "localhost:~#", 5 * time.Second},
		{"cat /etc/apk/repositories", "main", 5 * time.Second},
		// Setup SSH
		//{"setup-sshd", "Which ssh server?", 3 * time.Second},
		//{"openssh", "Allow root ssh login?", 3 * time.Second},
		//{"yes", "Enter ssh key or URL for root", 3 * time.Second},
		//{"", "ok ]", 10 * time.Second},
		{"apk add openssh", "localhost", 10 * time.Second},
		{"rc-update add sshd", "localhost", 10 * time.Second},
		//{"rc-service sshd --quiet start", "localhost", 10 * time.Second},

		//{"", "", 1 * time.Second},
		// Setup disk
		{"setup-disk -m sys /dev/sda", "WARNING: Erase the above disk", 5 * time.Second},
		{"y", "Installation is complete", 30 * time.Second},
	}
}

// Install performs the Alpine Linux installation
func (a *AlpineSetup) Install() error {
	steps := defaultInstallSteps()
	fmt.Println("\n=== Starting Alpine Linux Installation ===")

	for i, step := range steps {
		fmt.Printf("\n[Step %d/%d] Executing: %s (waiting %d for: %s)\n",
			i+1, len(steps), step.Command, step.Timeout, step.Trigger)

		if err := a.sendAndWait(step.Command, step.Trigger, step.Timeout); err != nil {
			return fmt.Errorf("step %d failed: %w", i+1, err)
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}
