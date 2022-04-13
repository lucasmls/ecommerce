package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcPrometheusInterceptors "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	gGRPC "google.golang.org/grpc"
)

// ServiceInput is the input (aka dependencies) to create a GRPC server.
type ServerInput struct {
	Port        int
	Registrator func(server gGRPC.ServiceRegistrar)
	Logger      *zap.Logger

	Listener net.Listener
}

// Server is the GRPC server itself.
type Server struct {
	in ServerInput

	server *gGRPC.Server
}

// NewServer is the GRPC Server constructor.
func NewServer(in ServerInput) (*Server, error) {
	if in.Port <= 80 {
		return nil, errors.New("the GRPC server port should be greater than or equal to 80")
	}

	if in.Registrator == nil {
		return nil, errors.New("missing required dependency: Registrator")
	}

	grpcPrometheusInterceptors.EnableHandlingTimeHistogram()
	grpcServer := gGRPC.NewServer(
		gGRPC.StreamInterceptor(
			grpcMiddleware.ChainStreamServer(
				otelgrpc.StreamServerInterceptor(),
				grpcPrometheusInterceptors.StreamServerInterceptor,
			),
		),
		gGRPC.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				otelgrpc.UnaryServerInterceptor(),
				grpcPrometheusInterceptors.UnaryServerInterceptor,
			),
		),
	)

	return &Server{
		in:     in,
		server: grpcServer,
	}, nil
}

// MustNewServer is the GRPC Server constructor.
// It panics if any error is found.
func MustNewServer(in ServerInput) *Server {
	server, err := NewServer(in)
	if err != nil {
		panic(err)
	}

	return server
}

func (s Server) Run(ctx context.Context) error {
	s.in.Registrator(s.server)

	address := fmt.Sprintf(":%d", s.in.Port)

	listener := s.in.Listener
	if listener == nil {
		l, err := net.Listen("tcp", address)
		if err != nil {
			s.in.Logger.Error("failed to start to listen TCP address", zap.Error(err))
			return err
		}

		listener = l
	}

	s.in.Logger.Info("gRPC server started:", zap.Int("port", s.in.Port))

	if err := s.server.Serve(listener); err != nil {
		s.in.Logger.Error("failed to serve gRPC server", zap.Error(err))
		return err
	}

	return nil
}
