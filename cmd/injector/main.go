package main

import (
	"context"
	"injector/internal/adapters/proxy"
	"injector/internal/drivers/cli"
	"injector/internal/drivers/logging"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func initLogger() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})))
}

func run() int {
	server, err := proxy.NewInjectorProxyServer(8080)
	if err != nil {
		slog.Error("failed to create proxy server", logging.WithError(err))
		return cli.Failure
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	if err = server.Serve(); err != nil {
		slog.Error("failed to start proxy server", logging.WithError(err))
		return cli.Failure
	}

	slog.Info("proxy server started")

	<-ctx.Done()
	stop()

	slog.Info("stopping proxy server")
	server.Stop()
	slog.Info("proxy server stopped")

	return cli.Success
}

func main() {
	initLogger()
	os.Exit(run())
}
