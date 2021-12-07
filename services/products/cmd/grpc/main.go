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

	"go.opentelemetry.io/otel"
	jaegerExporter "go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	tracingSdkResource "go.opentelemetry.io/otel/sdk/resource"
	tracingSdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

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
		tracingSdk.WithSampler(tracingSdk.AlwaysSample()),
		tracingSdk.WithResource(tracingSdkResource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("products"),
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

	tracer := otel.Tracer("products")

	productsInMemoryRepository := repositories.MustNewInMemoryProductsRepository(repositories.ProductsRepositoryInput{
		Logger: logger,
		Tracer: tracer,
		Size:   10,
	})

	application := app.MustNewApplication(app.ApplicationInput{
		Logger:             logger,
		Tracer:             tracer,
		ProductsRepository: productsInMemoryRepository,
	})

	productsResolver := resolvers.MustNewProductsResolver(resolvers.ProductsResolverInput{
		Logger: logger,
		Tracer: tracer,
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
