package broker

import (
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	"github.com/alexZaicev/message-broker/internal/domain/entities"
	"github.com/alexZaicev/message-broker/internal/drivers/logging"
	mbV1alpha1 "github.com/alexZaicev/message-broker/protobuf/go/messagebroker/v1alpha1"
)

func (b *MessageBroker) registerConsumer(conn net.Conn, errCount *atomic.Int32, req *mbV1alpha1.Request) {
	var subscribe mbV1alpha1.SubscribeRequest
	if err := proto.Unmarshal(req.Body, &subscribe); err != nil {
		slog.Error("failed to unmarshal subscribe request", logging.WithError(err))
		sendErrorResponse(conn, errCount, "Failed to read subscription request message")
		return
	}

	// check if queue exists
	queue, ok := b.mbCtx.FindQueueByName(subscribe.GetQueue())
	if !ok {
		sendErrorResponse(conn, errCount, fmt.Sprintf("Queue %s does not exist", subscribe.Queue))
		return
	}

	channel, ok := b.channels.Get(queue.ID)
	if !ok {
		// create new consumer channel
		b.channels.Add(NewChannel(conn, queue))
		sendAckResponse(conn)
		return
	}

	channel.RegisterConsumer(conn)
	sendAckResponse(conn)
}

func (b *MessageBroker) publishMessage(conn net.Conn, errCount *atomic.Int32, req *mbV1alpha1.Request) {
	var envelope mbV1alpha1.Envelope
	if err := proto.Unmarshal(req.Body, &envelope); err != nil {
		slog.Error("failed to unmarshal message envelope", logging.WithError(err))
		sendErrorResponse(conn, errCount, "Failed to read message envelope")
		return
	}

	exchange, ok := b.mbCtx.FindExchangeByName(envelope.GetExchange())
	if !ok {
		sendNackResponse(conn, MetadataKeyError, fmt.Sprintf("Exchange %s does not exist", envelope.GetExchange()))
		return
	}

	switch exchange.Kind {
	case entities.ExchangeKindQueue:
		b.publishMessageOverQueue(conn, exchange, &envelope)
	case entities.ExchangeKindTopic:
		b.publishMessageOverTopic(conn, exchange, &envelope)
	default:
		slog.Error("unknown exchange type", slog.Int("type", int(exchange.Kind)))
		sendErrorResponse(conn, errCount, "Failed to route message to exchange")
	}
}

func (b *MessageBroker) publishMessageOverQueue(
	conn net.Conn,
	exchange *entities.Exchange,
	envelope *mbV1alpha1.Envelope,
) {
	channel, ok := b.channels.Get(exchange.ID)
	if !ok {
		// TODO: maybe message should go into buffer?
		sendNackResponse(conn, MetadataKeyError, "Channel does not exist")
		return
	}

	sendAckResponse(conn)
	channel.Send(envelope)
}

func (b *MessageBroker) publishMessageOverTopic(
	conn net.Conn,
	exchange *entities.Exchange,
	envelope *mbV1alpha1.Envelope,
) {
	queues, err := b.mbCtx.FindQueuesForTopic(exchange.Name)
	if err != nil {
		sendNackResponse(conn, MetadataKeyError, "Topic does not exist")
	}

	sendAckResponse(conn)

	var wg sync.WaitGroup

	for _, queue := range queues {
		wg.Add(1)

		go func(queue *entities.Exchange) {
			defer wg.Done()
			channel, ok := b.channels.Get(queue.ID)
			if !ok {
				// TODO: maybe message should go into buffer?
				slog.Error("failed to find channel for queue", slog.String(logging.Queue, queue.Name))
				return
			}

			channel.Send(envelope)
		}(queue)
	}

	wg.Wait()
}
