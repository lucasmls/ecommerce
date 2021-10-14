package main

import (
	"context"

	"github.com/lucasmls/ecommerce/shared/grpc"
	"go.uber.org/zap"

	resolvers "github.com/lucasmls/ecommerce/services/products/ports/grpc"
	pb "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	gGRPC "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	productsResolver, err := resolvers.NewProductsResolver(resolvers.ProductsResolverInput{
		Logger: logger,
	})
	if err != nil {
		logger.Error("failed to instantiate Users resolver.", zap.Error(err))
		return
	}

	server, err := grpc.NewServer(grpc.ServerInput{
		Port: 8080,
		Registrator: func(server gGRPC.ServiceRegistrar) {
			pb.RegisterProductsServiceServer(server, productsResolver)
		},
	})
	if err != nil {
		logger.Error("failed to instantiate gRPC server.", zap.Error(err))
		return
	}

	if err := server.Run(ctx); err != nil {
		logger.Error("failed to run gRPC server.", zap.Error(err))
		return
	}
}
