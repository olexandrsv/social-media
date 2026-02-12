package connection

import "io"

type Connection interface {
	io.ReadWriteCloser
}

type MessageConnection interface {
	Send([]byte) error
	Receive() ([]byte, error)
	Close() error
}
