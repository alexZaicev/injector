package main

import (
	"context"
	"fmt"
	"github.com/hashicorp/cli"
	"injector/internal/drivers/codes"
	testservicev1 "injector/protobuf/go/testservice/v1"
)

type EchoCommand struct {
	ui cli.Ui
}

func (c *EchoCommand) Help() string {
	return ""
}

func (c *EchoCommand) Run(args []string) int {
	if len(args) != 1 {
		c.ui.Error("expected exactly one argument representing echo message")
		return codes.Failure
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	client, err := newTestServiceClient(ctx)
	if err != nil {
		c.ui.Error(fmt.Sprintf("failed to dial test service: %s", err))
		return codes.Failure
	}

	resp, err := client.Echo(ctx, &testservicev1.EchoRequest{
		Message: args[0],
	})
	if err != nil {
		c.ui.Error(fmt.Sprintf("failed to perform echo request: %s", err))
		return codes.Failure
	}

	c.ui.Info(fmt.Sprintf("received message: %s", resp.Message))

	return codes.Success
}

func (c *EchoCommand) Synopsis() string {
	return ""
}
