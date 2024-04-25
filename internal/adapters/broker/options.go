package broker

import "fmt"

const (
	defaultBrokerPort    = 6800
	defaultBufferSize    = 10 * 1024 * 1024
	defaultMaxErrorCount = 5
)

type brokerOptions struct {
	port          uint16
	bufferSize    uint64
	maxErrorCount int32
}

func defaultBrokerOptions() *brokerOptions {
	return &brokerOptions{
		port:          defaultBrokerPort,
		bufferSize:    defaultBufferSize,
		maxErrorCount: defaultMaxErrorCount,
	}
}

type Option func(*brokerOptions) error

func WithPort(port uint16) Option {
	return func(options *brokerOptions) error {
		if port < 1024 && port > 65535 {
			return fmt.Errorf("invalid port number: %d", port)
		}

		options.port = port

		return nil
	}
}

func WithBufferSize(bufferSize uint64) Option {
	return func(options *brokerOptions) error {
		options.bufferSize = bufferSize
		return nil
	}
}
