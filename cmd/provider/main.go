package main

import (
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/alexZaicev/message-broker/internal/drivers/codes"
	"github.com/alexZaicev/message-broker/internal/drivers/logging"
	mbV1alpha1 "github.com/alexZaicev/message-broker/protobuf/go/messagebroker/v1alpha1"
)

func getMessageBytes(payload []byte) ([]byte, error) {
	envelope, marshalErr := proto.Marshal(&mbV1alpha1.Envelope{
		Exchange: "echoQueue",
		Payload:  payload,
	})
	if marshalErr != nil {
		return nil, marshalErr
	}

	request, marshalErr := proto.Marshal(&mbV1alpha1.Request{
		Type:     mbV1alpha1.EnumRequestType_ENUM_REQUEST_TYPE_MESSAGE,
		Body:     envelope,
		Metadata: nil,
	})
	if marshalErr != nil {
		return nil, marshalErr
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

	// send message every 3 seconds
	for {
		msg, marshalErr := getMessageBytes([]byte("hello world"))
		if marshalErr != nil {
			slog.Error("failed to encode message", logging.WithError(marshalErr))
			return codes.Failure
		}

		if _, writeErr := conn.Write(msg); writeErr != nil {
			slog.Error("failed to write message", logging.WithError(writeErr))
			return codes.Failure
		}

		buf := make([]byte, 1024)

		n, readErr := conn.Read(buf)
		if readErr != nil {
			slog.Error("failed to read response", logging.WithError(readErr))
			return codes.Failure
		}

		var response mbV1alpha1.Response

		if unmarshalErr := proto.Unmarshal(buf[:n], &response); unmarshalErr != nil {
			slog.Info("received response", logging.WithError(unmarshalErr))
			return codes.Failure
		}

		slog.Info("response received", slog.String("status", string(response.GetStatus())))

		time.Sleep(3 * time.Second)
	}
}

func main() {
	logging.SetupConsoleLogger()
	os.Exit(run())
}
