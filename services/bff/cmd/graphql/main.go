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

	"go.opentelemetry.io/otel"
	jaegerExporter "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	tracingSdkResource "go.opentelemetry.io/otel/sdk/resource"
	tracingSdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
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

	jaegerExporter, err := jaegerExporter.New(
		jaegerExporter.WithCollectorEndpoint(
			jaegerExporter.WithEndpoint("http://localhost:14268/api/traces"),
		),
	)
	if err != nil {
		logger.Error("failed to instantiate Jaeger exporter", zap.Error(err))
		return
	}

	tracingProvider := tracingSdk.NewTracerProvider(
		tracingSdk.WithBatcher(jaegerExporter),
		tracingSdk.WithResource(tracingSdkResource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("bff"),
		)),
	)

	defer func() {
		_ = tracingProvider.Shutdown(ctx)
	}()

	otel.SetTracerProvider(tracingProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer := otel.Tracer("bff")

	productsServiceGRPCClient := grpc.MustNewClient(grpc.ClientInput{
		Address: "localhost:8081",
		Logger:  logger,
	})

	productsServiceConn := productsServiceGRPCClient.MustConnect(ctx)

	productsService := productsPb.NewProductsServiceClient(productsServiceConn)

	graphQlResolver := &graph.Resolver{
		Logger:          logger,
		Tracer:          tracer,
		ProductsService: productsService,
	}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: graphQlResolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
