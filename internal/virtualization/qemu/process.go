package qemu

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
)

type ProcessState struct {
	Running bool
	Error   error
}

type QemuProcess struct {
	cmd       *exec.Cmd
	stdin     io.WriteCloser
	stdout    io.ReadCloser
	done      chan struct{}
	state     *ProcessState
	stateLock sync.RWMutex
}

func NewQemuProcess(args []string) (*QemuProcess, error) {
	cmd := exec.Command("qemu-system-x86_64", args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		stdin.Close()
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	cmd.Stderr = os.Stderr

	return &QemuProcess{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		done:   make(chan struct{}),
		state:  &ProcessState{Running: false},
	}, nil
}

func (p *QemuProcess) Start() error {
	if err := p.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	p.stateLock.Lock()
	p.state.Running = true
	p.stateLock.Unlock()

	// Start output reader
	go p.handleOutput()

	// Start process monitor
	go p.monitor()
	return nil
}

func (p *QemuProcess) handleOutput() {
	reader := bufio.NewReader(p.stdout)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				p.setError(fmt.Errorf("output reading error: %w", err))
			}
			return
		}
		fmt.Print(line)
	}
}

func (p *QemuProcess) monitor() {
	err := p.cmd.Wait()
	p.stateLock.Lock()
	p.state.Running = false
	if err != nil {
		p.state.Error = err
	}
	p.stateLock.Unlock()
	close(p.done)
}

func (p *QemuProcess) Stop() error {
	if !p.IsRunning() {
		return nil
	}

	// Try graceful shutdown first
	if err := p.cmd.Process.Signal(syscall.SIGTERM); err != nil {
		// Force kill if graceful shutdown fails
		if err := p.cmd.Process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	}

	// Wait for process to finish
	<-p.done
	return p.state.Error
}

func (p *QemuProcess) IsRunning() bool {
	p.stateLock.RLock()
	defer p.stateLock.RUnlock()
	return p.state.Running
}

func (p *QemuProcess) setError(err error) {
	p.stateLock.Lock()
	p.state.Error = err
	p.stateLock.Unlock()
}

func (p *QemuProcess) GetError() error {
	p.stateLock.RLock()
	defer p.stateLock.RUnlock()
	return p.state.Error
}

func (p *QemuProcess) Send(message string) error {
	if !p.IsRunning() {
		return fmt.Errorf("process is not running")
	}

	// Ensure message ends with newline
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}

	_, err := p.stdin.Write([]byte(message))
	if err != nil {
		p.setError(fmt.Errorf("failed to write to process: %w", err))
		return err
	}

	return nil
}
