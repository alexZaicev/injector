package main

import (
	"context"
	"google.golang.org/grpc"
	"injector/internal/drivers/codes"
	"injector/internal/drivers/logging"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func run() int {
	listener, err := net.Listen("tcp", ":8090")
	if err != nil {
		return codes.Failure
	}

	server := grpc.NewServer()

	api := NewTestServiceAPI()
	api.RegisterService(server)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err = server.Serve(listener); err != nil {
			slog.Error("failed to start grpc test server", logging.WithError(err))
		}
	}()

	slog.Info("grpc test server started")

	<-ctx.Done()
	stop()

	slog.Info("stopping grpc test server")
	server.Stop()
	slog.Info("grpc test server stopped")

	return codes.Success
}

func main() {
	logging.SetupJSONLogger()
	os.Exit(run())
}
