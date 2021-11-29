package main

import (
	"context"

	"github.com/lucasmls/ecommerce/shared/grpc"
	"go.uber.org/zap"

	"github.com/lucasmls/ecommerce/services/products/adapters/repositories"
	"github.com/lucasmls/ecommerce/services/products/app"
	resolvers "github.com/lucasmls/ecommerce/services/products/ports/grpc"
	pb "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	gGRPC "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	productsInMemoryRepository := repositories.MustNewInMemoryProductsRepository(10)
	application := app.MustNewApplication(logger, productsInMemoryRepository)

	productsResolver := resolvers.MustNewProductsResolver(resolvers.ProductsResolverInput{
		Logger: logger,
		App:    application,
	})

	server, err := grpc.NewServer(grpc.ServerInput{
		Port:   8081,
		Logger: logger,
		Registrator: func(server gGRPC.ServiceRegistrar) {
			pb.RegisterProductsServiceServer(server, productsResolver)
		},
	})
	if err != nil {
		logger.Error("failed to instantiate gRPC server", zap.Error(err))
		return
	}

	if err := server.Run(ctx); err != nil {
		logger.Error("failed to run gRPC server", zap.Error(err))
		return
	}
}
