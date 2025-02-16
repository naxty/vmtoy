package common

type Connection interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Close() error
}

type Installer interface {
	Install() error
	Close() error
}
