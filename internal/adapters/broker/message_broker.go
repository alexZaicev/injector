package broker

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	"github.com/alexZaicev/message-broker/internal/domain/entities"
	mbErrors "github.com/alexZaicev/message-broker/internal/drivers/errors"
	"github.com/alexZaicev/message-broker/internal/drivers/logging"
	mbV1alpha1 "github.com/alexZaicev/message-broker/protobuf/go/messagebroker/v1alpha1"
)

type MessageBroker struct {
	brokerOptions *brokerOptions

	listener net.Listener
	mbCtx    *entities.MessageBrokerContext
	channels *ChannelMap
}

func New(mbCtx *entities.MessageBrokerContext, options ...Option) (*MessageBroker, error) {
	opt := defaultBrokerOptions()
	for _, option := range options {
		if err := option(opt); err != nil {
			return nil, err
		}
	}

	return &MessageBroker{
		mbCtx:         mbCtx,
		brokerOptions: opt,
		channels:      newChannelMap(),
	}, nil
}

func (b *MessageBroker) Start() error {
	tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", b.brokerOptions.port))
	if err != nil {
		return err
	}

	listener, err := net.ListenTCP("tcp", tcpAddr)
	if err != nil {
		return err
	}

	b.listener = listener

	go b.accept()

	return nil
}

func (b *MessageBroker) accept() {
	for {
		conn, err := b.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}

			slog.Debug("failed to accept connection", logging.WithError(err))
			continue
		}

		// accept incoming connection
		go b.process(conn)
	}
}

func (b *MessageBroker) process(conn net.Conn) {
	defer func() {
		sendDisconnectResponse(conn)
		if err := conn.Close(); err != nil {
			slog.Error("failed to close connection", logging.WithError(err))
		}
	}()

	errCount := &atomic.Int32{}

	for {
		// drop connection if error count exceed maximum
		if errCount.Load() >= b.brokerOptions.maxErrorCount {
			slog.Debug("connection dropped due to maximum error count")
			return
		}

		buf := make([]byte, b.brokerOptions.bufferSize)
		n, err := conn.Read(buf)
		if err != nil {
			if mbErrors.IsClosedConn(err) || mbErrors.IsEOF(err) {
				return
			}

			slog.Error("failed to read received message request", logging.WithError(err))
			sendErrorResponse(conn, errCount, "Failed to read request message")
			continue
		}

		var request mbV1alpha1.Request
		if err = proto.Unmarshal(buf[:n], &request); err != nil {
			slog.Error("failed to unmarshal received message request", logging.WithError(err))
			sendErrorResponse(conn, errCount, "Failed to read request message")
			continue
		}

		switch request.GetType() {
		case mbV1alpha1.EnumRequestType_ENUM_REQUEST_TYPE_SUBSCRIBE:
			b.registerConsumer(conn, errCount, &request)
		case mbV1alpha1.EnumRequestType_ENUM_REQUEST_TYPE_MESSAGE:
			b.publishMessage(conn, errCount, &request)
		default:
			slog.Error("unknown request type", slog.String("type", string(request.GetType())))
			sendErrorResponse(conn, errCount, "Failed to read request message")
		}
	}
}

func (b *MessageBroker) Stop() {
	for _, queueID := range b.channels.Keys() {
		ch, ok := b.channels.Get(queueID)
		if ok {
			ch.Close()
		}

		b.channels.Remove(queueID)
	}

	if err := b.listener.Close(); err != nil {
		slog.Error("failed to close listener", logging.WithError(err))
	}
}
