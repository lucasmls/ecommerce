package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	graph "github.com/lucasmls/ecommerce/services/bff/ports/graphql"
	"github.com/lucasmls/ecommerce/services/bff/ports/graphql/generated"
	productsPb "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"github.com/lucasmls/ecommerce/shared/grpc"
	"go.uber.org/zap"
)

const defaultPort = "8080"

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	productsServiceGRPCClient, err := grpc.NewClient(grpc.ClientInput{
		Address: "localhost:8081",
		Logger:  logger,
	})
	if err != nil {
		logger.Error("failed to build products service gRPC client", zap.Error(err))
		return
	}

	productsServiceConn, err := productsServiceGRPCClient.Connect(ctx)
	if err != nil {
		logger.Error("failed to connect into products service gRPC server.", zap.Error(err))
		return
	}

	productsService := productsPb.NewProductsServiceClient(productsServiceConn)

	graphQlResolver := &graph.Resolver{
		Logger:          logger,
		ProductsService: productsService,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphQlResolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
