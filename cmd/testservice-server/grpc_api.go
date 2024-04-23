package main

import (
	"context"
	"google.golang.org/grpc"
	testservicev1 "injector/protobuf/go/testservice/v1"
)

type TestServiceAPI struct {
	testservicev1.UnimplementedTestServiceServer
}

func NewTestServiceAPI() *TestServiceAPI {
	return &TestServiceAPI{}
}

func (api *TestServiceAPI) RegisterService(server *grpc.Server) {
	testservicev1.RegisterTestServiceServer(server, api)
}

func (api *TestServiceAPI) Echo(_ context.Context, req *testservicev1.EchoRequest) (*testservicev1.EchoResponse, error) {
	return &testservicev1.EchoResponse{
		Message: req.Message,
	}, nil
}
