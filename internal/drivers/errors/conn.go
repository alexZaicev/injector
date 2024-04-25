package errors

import (
	"errors"
	"io"
	"net"
)

func IsClosedConn(err error) bool {
	return errors.Is(err, net.ErrClosed)
}

func IsEOF(err error) bool {
	return errors.Is(err, io.EOF)
}
