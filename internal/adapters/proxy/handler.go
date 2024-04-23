package proxy

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"strings"
)

const (
	XRealIP       = "X-Real-IP"
	XForwardedFor = "X-Forwarded-For"
)

type HandlerOptions struct {
	Payload io.Reader
	Request *http.Request
}

func NewHandlerOptionsFromPayload(payload []byte) (*HandlerOptions, error) {
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(payload)))
	if err != nil {
		return nil, err
	}

	return &HandlerOptions{
		Payload: bytes.NewReader(payload),
		Request: req,
	}, nil
}

func (o HandlerOptions) isGRPC() bool {
	return o.Request.ProtoMajor >= 2
}

func (o HandlerOptions) getSourceAddr() string {
	val := o.Request.Header.Get(XRealIP)
	if val != "" {
		return val
	}

	val = o.Request.Header.Get(XForwardedFor)
	if val != "" {
		// X-Forwarded-For header value is a collection IPs where the request travelled
		// The left most IP address, would be where the request originated from
		ips := strings.Split(val, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	return o.Request.RemoteAddr
}

func (o HandlerOptions) getDestinationAddr() string {
	return o.Request.Host
}

type Handler interface {
	Handle(net.Conn, *HandlerOptions) error
}
