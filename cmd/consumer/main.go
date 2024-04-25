package main

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"google.golang.org/protobuf/proto"

	"github.com/alexZaicev/message-broker/internal/drivers/codes"
	"github.com/alexZaicev/message-broker/internal/drivers/logging"
	mbV1alpha1 "github.com/alexZaicev/message-broker/protobuf/go/messagebroker/v1alpha1"
)

func subscribe(conn net.Conn, queue string) error {
	request, err := getSubscriptionRequest(queue)
	if err != nil {
		slog.Error("failed to get subscription request", logging.WithError(err))
		return err
	}

	if _, writeErr := conn.Write(request); writeErr != nil {
		slog.Error("failed to write subscription request", logging.WithError(writeErr))
		return writeErr
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	var response mbV1alpha1.Response

	if err = proto.Unmarshal(buf[:n], &response); err != nil {
		return err
	}

	if response.GetStatus() != mbV1alpha1.EnumStatus_ENUM_STATUS_ACK {
		return fmt.Errorf("failed to subscribe, status not ACK")
	}

	return nil
}

func getSubscriptionRequest(queue string) ([]byte, error) {
	subscription, err := proto.Marshal(&mbV1alpha1.SubscribeRequest{
		Queue: queue,
	})
	if err != nil {
		return nil, err
	}

	request, err := proto.Marshal(&mbV1alpha1.Request{
		Type:     mbV1alpha1.EnumRequestType_ENUM_REQUEST_TYPE_SUBSCRIBE,
		Body:     subscription,
		Metadata: nil,
	})
	if err != nil {
		return nil, err
	}

	return request, nil
}

func run() int {
	tcpAddr, err := net.ResolveTCPAddr("tcp", ":6800")
	if err != nil {
		slog.Error("failed to resolve tcp address", logging.WithError(err))
		return codes.Failure
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		slog.Error("failed to dial tcp", logging.WithError(err))
		return codes.Failure
	}
	defer conn.Close()

	if err = subscribe(conn, "echoQueue"); err != nil {
		slog.Error("failed to subscribe to queue", logging.WithError(err))
		return codes.Failure
	}

	for {
		buf := make([]byte, 1024)

		n, readErr := conn.Read(buf)
		if readErr != nil {
			return codes.Failure
		}

		var msg mbV1alpha1.Envelope

		if readErr = proto.Unmarshal(buf[:n], &msg); readErr != nil {
			return codes.Failure
		}

		slog.Info(string(msg.GetPayload()))
	}
}

func main() {
	logging.SetupConsoleLogger()
	os.Exit(run())
}
