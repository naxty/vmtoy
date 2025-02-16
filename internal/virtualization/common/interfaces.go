package common

type VMConfig interface {
	BuildCommand() []string
	GetName() string
	GetState() string
}

type VM interface {
	Start() error
	Stop() error
	Install() error
	Status() string
}

type Installer interface {
	Install() error
	Close() error
}
