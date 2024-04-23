package proxy

import (
	"bytes"
	"injector/internal/drivers/logging"
	"io"
	"log/slog"
	"net"
	"net/http"
)

type HTTPHandler struct {
}

func NewHTTPHandler() *HTTPHandler {
	return &HTTPHandler{}
}

func (h *HTTPHandler) Handle(conn net.Conn, options *HandlerOptions) error {
	logger := slog.With(
		slog.String(logging.Protocol, "http"),
		slog.String(logging.Path, options.Request.URL.RequestURI()),
		slog.String(logging.Method, options.Request.Method),
		slog.String(logging.Source, options.getSourceAddr()),
		slog.String(logging.Destination, options.getDestinationAddr()),
	)

	logger.Info("preparing to handle http request")

	resp := http.Response{
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		ProtoMinor: 0,
		Body:       io.NopCloser(bytes.NewBufferString("Hello World")),
	}
	resp.StatusCode = http.StatusOK
	return resp.Write(conn)
}
