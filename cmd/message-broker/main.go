package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/alexZaicev/message-broker/internal/adapters/broker"
	"github.com/alexZaicev/message-broker/internal/adapters/config"
	"github.com/alexZaicev/message-broker/internal/domain/entities"
	"github.com/alexZaicev/message-broker/internal/drivers/codes"
	"github.com/alexZaicev/message-broker/internal/drivers/logging"
)

func run() int {
	slog.Info("preparing message broker")

	mbCtx := entities.NewMessageBrokerContext()
	if err := config.LoadConfiguration(mbCtx); err != nil {
		slog.Error("message broker failed to configure", logging.WithError(err))
		return codes.Failure
	}

	messageBroker, err := broker.New(mbCtx)
	if err != nil {
		slog.Error("failed to create message broker", logging.WithError(err))
		return codes.Failure
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err = messageBroker.Start(); err != nil {
		slog.Error("failed to start message broker", logging.WithError(err))
		return codes.Failure
	}

	slog.Info("message broker ready to serve")

	<-ctx.Done()
	stop()

	slog.Info("stopping message broker")
	messageBroker.Stop()
	slog.Info("message broker stopped successfully")

	return codes.Success
}

func main() {
	logging.SetupConsoleLogger()
	os.Exit(run())
}
