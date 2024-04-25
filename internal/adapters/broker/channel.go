package broker

import (
	"errors"
	"log/slog"
	"net"
	"sync"

	"google.golang.org/protobuf/proto"

	"github.com/alexZaicev/message-broker/internal/domain/entities"
	mbErrors "github.com/alexZaicev/message-broker/internal/drivers/errors"
	"github.com/alexZaicev/message-broker/internal/drivers/logging"
	mbV1alpha1 "github.com/alexZaicev/message-broker/protobuf/go/messagebroker/v1alpha1"
)

type Channel struct {
	queue *entities.Exchange

	mu        sync.Mutex
	consumers []net.Conn
}

func NewChannel(conn net.Conn, queue *entities.Exchange) *Channel {
	return &Channel{
		queue:     queue,
		consumers: []net.Conn{conn},
	}
}

func (ch *Channel) RegisterConsumer(conn net.Conn) {
	ch.consumers = append(ch.consumers, conn)
}

func (ch *Channel) Send(envelope *mbV1alpha1.Envelope) {
	logger := slog.With(slog.String(logging.Queue, ch.queue.Name))

	var (
		wg              sync.WaitGroup
		mu              sync.Mutex
		failedConsumers []int
	)

	for idx, consumer := range ch.consumers {
		wg.Add(1)

		go func(idx int, consumer net.Conn) {
			defer wg.Done()

			bytes, err := proto.Marshal(envelope)
			if err != nil {
				logger.Error("failed to marshal response", logging.WithError(err))
				return
			}
			if _, err = consumer.Write(bytes); err != nil {
				if errors.Is(err, net.ErrClosed) {
					logger.Debug("consumer disconnected from queue")
					mu.Lock()
					failedConsumers = append(failedConsumers, idx)
					mu.Unlock()

					return
				}

				logger.Error("failed to send message to consumer", logging.WithError(err))
				return
			}
		}(idx, consumer)
	}

	wg.Wait()

	// remove disconnected consumers
	ch.mu.Lock()
	for i := len(failedConsumers) - 1; i >= 0; i-- {
		f := failedConsumers[i]
		ch.consumers = append(ch.consumers[:f], ch.consumers[f+1:]...)
	}
	ch.mu.Unlock()
}

func (ch *Channel) Close() {
	logger := slog.With(slog.String(logging.Queue, ch.queue.Name))
	logger.Debug("closing channel")

	var wg sync.WaitGroup

	for _, consumer := range ch.consumers {
		if err := consumer.Close(); err != nil {
			if !mbErrors.IsClosedConn(err) {
				slog.Error("failed to disconnect consumer", logging.WithError(err))
			}
		}
	}

	wg.Wait()

	logger.Debug("channel closed successfully")
}
