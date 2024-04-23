package proxy

import (
	"errors"
	"fmt"
	"injector/internal/drivers/logging"
	"io"
	"log/slog"
	"net"
	"time"
)

const (
	minPort = 1025
	maxPort = 65535
)

// InjectorProxyServer satisfies the Server interfaces providing the mechanism to fetch and modify incoming HTTP/1 and
// HTTP/2 requests before forwarding them to the receiving service.
type InjectorProxyServer struct {
	serverOptions *serverOptions

	tcpAddr  *net.TCPAddr
	listener net.Listener

	httpHandler Handler
	grpcHandler Handler
}

func NewInjectorProxyServer(port uint16, options ...Option) (*InjectorProxyServer, error) {
	if port < minPort || port > maxPort {
		return nil, fmt.Errorf("port %d is out of range", port)
	}

	opt := newDefaultOptions()

	for _, option := range options {
		if err := option(opt); err != nil {
			return nil, err
		}
	}

	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	return &InjectorProxyServer{
		serverOptions: opt,
		tcpAddr:       tcpAddr,
		httpHandler:   NewHTTPHandler(),
		grpcHandler:   NewGRPCHandler(),
	}, nil
}

// Serve initializes the TCP listener and ensures that new connections are served.
func (s *InjectorProxyServer) Serve() error {
	var err error

	s.listener, err = net.ListenTCP("tcp", s.tcpAddr)
	if err != nil {
		return err
	}

	go s.listen()

	return nil
}

func (s *InjectorProxyServer) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) || errors.Is(err, io.EOF) {
				return
			}

			slog.Error("failed to accept connection", logging.WithError(err))
			continue
		}

		if err = conn.SetReadDeadline(time.Now().Add(s.serverOptions.readDeadline)); err != nil {
			slog.Error("failed to set read deadline", logging.WithError(err))
			return
		}

		go s.handleConnection(conn)
	}
}

func (s *InjectorProxyServer) handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Error("failed to close connection", logging.WithError(err))
		}
	}()

	// read payload with limit
	buf := make([]byte, s.serverOptions.payloadSize)
	n, err := conn.Read(buf)
	if err != nil {
		slog.Error("failed to read connection payload", logging.WithError(err))
		return
	}

	handlerOptions, err := NewHandlerOptionsFromPayload(buf[:n])
	if err != nil {
		slog.Error("failed to create handler options from payload", logging.WithError(err))
	}

	if handlerOptions.isGRPC() {
		err = s.grpcHandler.Handle(conn, handlerOptions)
	} else {
		err = s.httpHandler.Handle(conn, handlerOptions)
	}

	if err != nil {
		slog.Error("failed to handle http request")
	}
}

// Stop closes down the TCP listener.
func (s *InjectorProxyServer) Stop() {
	if s.listener == nil {
		return
	}

	if err := s.listener.Close(); err != nil {
		slog.Error("failed to stop that server", logging.WithError(err))
	}
}
