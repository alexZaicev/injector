package proxy

import (
	"fmt"
	"time"
)

const (
	defaultReadDeadlineInSeconds = 60
	defaultPayloadSize           = 10 * 1024 * 1024
)

type serverOptions struct {
	readDeadline time.Duration
	payloadSize  int64
}

type Option func(*serverOptions) error

func newDefaultOptions() *serverOptions {
	return &serverOptions{
		readDeadline: defaultReadDeadlineInSeconds * time.Second,
		payloadSize:  defaultPayloadSize,
	}
}

func WithReadDeadline(deadlineInSeconds uint64) Option {
	return func(o *serverOptions) error {
		if deadlineInSeconds == 0 {
			return fmt.Errorf("dealine must not be zero")
		}

		o.readDeadline = time.Duration(deadlineInSeconds) * time.Second
		return nil
	}
}

func WithPayloadSize(size int64) Option {
	return func(options *serverOptions) error {
		if size == 0 {
			return fmt.Errorf("payload size must not be zero")
		}

		options.payloadSize = size
		return nil
	}
}
