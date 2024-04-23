package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	testservicev1 "injector/protobuf/go/testservice/v1"
	"time"
)

const (
	defaultTimeout = time.Second * 5
)

func newTestServiceClient(ctx context.Context) (testservicev1.TestServiceClient, error) {
	conn, err := grpc.DialContext(ctx, "localhost:8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return testservicev1.NewTestServiceClient(conn), nil
}
