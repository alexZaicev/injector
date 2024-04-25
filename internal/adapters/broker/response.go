package broker

import (
	"log/slog"
	"net"
	"sync/atomic"

	"google.golang.org/protobuf/proto"

	mbErrors "github.com/alexZaicev/message-broker/internal/drivers/errors"
	"github.com/alexZaicev/message-broker/internal/drivers/logging"
	mbV1alpha1 "github.com/alexZaicev/message-broker/protobuf/go/messagebroker/v1alpha1"
)

type Metadata map[string]string

func newResponse(status mbV1alpha1.EnumStatus, metadata Metadata) ([]byte, error) {
	resp, err := proto.Marshal(&mbV1alpha1.Response{
		Status:   status,
		Metadata: metadata,
	})
	if err != nil {
		slog.Error("failed to marshal response", logging.WithError(err))
		return nil, err
	}

	return resp, nil
}

func sendAckResponse(conn net.Conn, metadataPairs ...string) {
	sendResponse(conn, mbV1alpha1.EnumStatus_ENUM_STATUS_ACK, metadataPairs...)
}

func sendNackResponse(conn net.Conn, metadataPairs ...string) {
	sendResponse(conn, mbV1alpha1.EnumStatus_ENUM_STATUS_NACK, metadataPairs...)
}

func sendDisconnectResponse(conn net.Conn) {
	sendResponse(conn, mbV1alpha1.EnumStatus_ENUM_STATUS_DISCONNECT)
}

func sendErrorResponse(conn net.Conn, errCount *atomic.Int32, message string) {
	defer func() {
		errCount.Add(1)
	}()

	sendResponse(conn, mbV1alpha1.EnumStatus_ENUM_STATUS_ERR, MetadataKeyError, message)
}

func sendResponse(conn net.Conn, status mbV1alpha1.EnumStatus, metadataPairs ...string) {
	if len(metadataPairs)%2 != 0 {
		panic("metadata length must be even")
	}

	metadata := Metadata{}

	for i := 0; i < len(metadataPairs); i += 2 {
		metadata[metadataPairs[i]] = metadataPairs[i+1]
	}

	resp, err := newResponse(status, metadata)
	if err != nil {
		return
	}

	if _, err = conn.Write(resp); err != nil {
		if !mbErrors.IsClosedConn(err) {
			slog.Error("failed to write response", logging.WithError(err))
		}
	}
}
