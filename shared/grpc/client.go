package grpc

import (
	"context"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	gGRPC "google.golang.org/grpc"
)

// ClientInput is the input (aka dependencies) needed to create a gRPC client
type ClientInput struct {
	Address string
	Logger  *zap.Logger

	AdditionalDialOptions []gGRPC.DialOption
}

// Client is the gRPC client itself
type Client struct {
	in ClientInput
}

// NewClient is the gRPC client constructor
func NewClient(in ClientInput) (*Client, error) {
	return &Client{in: in}, nil
}

// MustNewClient is the gRPC client constructor
// It panics if any error is found
func MustNewClient(in ClientInput) *Client {
	client, err := NewClient(in)
	if err != nil {
		panic(err)
	}

	return client
}

// Connect into a gRPC server
func (c Client) Connect(ctx context.Context) (*gGRPC.ClientConn, error) {
	c.in.Logger.Info("connecting to gRPC server", zap.String("address", c.in.Address))

	dialOptions := []gGRPC.DialOption{
		gGRPC.WithInsecure(),
		gGRPC.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
		gGRPC.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()),
	}

	dialOptions = append(dialOptions, c.in.AdditionalDialOptions...)

	conn, err := gGRPC.DialContext(
		ctx,
		c.in.Address,
		dialOptions...,
	)
	if err != nil {
		c.in.Logger.Error("failed to connect into gRPC server", zap.Error(err))
		return nil, err
	}

	return conn, nil
}

// MustConnect into a gRPC server
// It panics if any error is found
func (c Client) MustConnect(ctx context.Context) *gGRPC.ClientConn {
	connection, err := c.Connect(ctx)
	if err != nil {
		panic(err)
	}

	return connection
}
