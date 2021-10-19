package grpc

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	gGRPC "google.golang.org/grpc"
)

// ClientInput is the input (aka dependencies) needed to create a gRPC client
type ClientInput struct {
	Address string
	Logger  *zap.Logger
}

// Client is the gRPC client itself
type Client struct {
	in ClientInput
}

// NewClient is the gRPC client constructor
func NewClient(in ClientInput) (*Client, error) {
	return &Client{in: in}, nil
}

// Connect into a gRPC server
func (c Client) Connect(ctx context.Context) (*grpc.ClientConn, error) {
	c.in.Logger.Info("connecting to gRPc server", zap.String("address", c.in.Address))

	conn, err := gGRPC.Dial(c.in.Address, grpc.WithInsecure())
	if err != nil {
		c.in.Logger.Error("failed to connect into gRPC server", zap.Error(err))
		return nil, err
	}

	return conn, nil
}
