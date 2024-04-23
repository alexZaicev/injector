package proxy

import (
	"fmt"
	"golang.org/x/net/http2"
	"injector/internal/drivers/logging"
	"log/slog"
	"net"
)

type GRPCHandler struct {
}

func NewGRPCHandler() *GRPCHandler {
	return &GRPCHandler{}
}

func (h *GRPCHandler) Handle(conn net.Conn, options *HandlerOptions) error {
	logger := slog.With(
		slog.String(logging.Protocol, "http"),
		slog.String(logging.Path, options.Request.URL.RequestURI()),
		slog.String(logging.Method, options.Request.Method),
		slog.String(logging.Source, options.getSourceAddr()),
		slog.String(logging.Destination, options.getDestinationAddr()),
	)

	framer := http2.NewFramer(conn, conn)
	if err := framer.WriteSettingsAck(); err != nil {
		logger.Error("failed to send settings ack", logging.WithError(err))
		return fmt.Errorf("failed to send settings ack")
	}

	return nil
}
