package main

import (
	"context"
	"net/http"

	"github.com/lucasmls/ecommerce/services/products/adapters/repositories"
	"github.com/lucasmls/ecommerce/services/products/app"
	resolvers "github.com/lucasmls/ecommerce/services/products/ports/grpc"
	protog "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"github.com/lucasmls/ecommerce/shared/grpc"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otel "go.opentelemetry.io/otel"
	otelJaegerExporter "go.opentelemetry.io/otel/exporters/jaeger"
	otelPropagation "go.opentelemetry.io/otel/propagation"
	otelSdkResource "go.opentelemetry.io/otel/sdk/resource"
	otelTraceSdk "go.opentelemetry.io/otel/sdk/trace"
	otelSemconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	gGRPC "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	jaegerExporter, err := otelJaegerExporter.New(
		otelJaegerExporter.WithCollectorEndpoint(
			otelJaegerExporter.WithEndpoint("http://localhost:14268/api/traces"),
		),
	)
	if err != nil {
		logger.Error("failed to instantiate Jaeger exporter", zap.Error(err))
		return
	}

	tracingProvider := otelTraceSdk.NewTracerProvider(
		otelTraceSdk.WithBatcher(jaegerExporter),
		otelTraceSdk.WithSampler(otelTraceSdk.AlwaysSample()),
		otelTraceSdk.WithResource(otelSdkResource.NewWithAttributes(
			otelSemconv.SchemaURL,
			otelSemconv.ServiceNameKey.String("products"),
		)),
	)

	defer func() {
		_ = tracingProvider.Shutdown(ctx)
	}()

	otel.SetTracerProvider(tracingProvider)
	otel.SetTextMapPropagator(otelPropagation.NewCompositeTextMapPropagator(
		otelPropagation.TraceContext{},
		otelPropagation.Baggage{},
	))

	tracer := otel.Tracer("products")

	productsInMemoryRepository := repositories.MustNewInMemoryProductsRepository(logger, tracer, 10)

	// pgProductsRepository, err := repositories.NewPgProductsRepository(
	// 	"dbname=products_service user=postgres password=postgres host=localhost port=5432 sslmode=disable",
	// )
	// if err != nil {
	// 	logger.Error("failed to instantiate PgProductsRepository", zap.Error(err))
	// 	return
	// }

	application := app.MustNewApplication(logger, tracer, productsInMemoryRepository)
	productsResolver := resolvers.MustNewProductsResolver(logger, tracer, application)

	server := grpc.MustNewServer(grpc.ServerInput{
		Port:   8081,
		Logger: logger,
		Registrator: func(server gGRPC.ServiceRegistrar) {
			protog.RegisterProductsServiceServer(server, productsResolver)
		},
	})

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":2112", nil)
	}()

	if err := server.Run(ctx); err != nil {
		logger.Error("failed to run gRPC server", zap.Error(err))
		return
	}
}
