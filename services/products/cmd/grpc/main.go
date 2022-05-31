package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/lucasmls/ecommerce/services/products/adapters/repositories"
	"github.com/lucasmls/ecommerce/services/products/app"
	resolvers "github.com/lucasmls/ecommerce/services/products/ports/grpc"
	protog "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"github.com/lucasmls/ecommerce/shared/env"
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

type ApplicationConfig struct {
	ServiceName              string `mapstructure:"SERVICE_NAME"`
	JaegerEndpoint           string `mapstructure:"JAEGER_ENDPOINT"`
	GrpcServerPort           int    `mapstructure:"GRPC_SERVER_PORT"`
	MetricsPort              int    `mapstructure:"METRICS_PORT"`
	PostgresConnectionString string `mapstructure:"PG_CONNECTION_STRING"`
}

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	config, err := env.LoadConfig[ApplicationConfig]()
	if err != nil {
		logger.Fatal("failed to load application config", zap.Error(err))
	}

	jaegerExporter, err := otelJaegerExporter.New(
		otelJaegerExporter.WithCollectorEndpoint(
			otelJaegerExporter.WithEndpoint(config.JaegerEndpoint),
		),
	)
	if err != nil {
		logger.Fatal("failed to instantiate Jaeger exporter", zap.Error(err))
	}

	tracingProvider := otelTraceSdk.NewTracerProvider(
		otelTraceSdk.WithBatcher(jaegerExporter),
		otelTraceSdk.WithSampler(otelTraceSdk.AlwaysSample()),
		otelTraceSdk.WithResource(otelSdkResource.NewWithAttributes(
			otelSemconv.SchemaURL,
			otelSemconv.ServiceNameKey.String(config.ServiceName),
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

	tracer := otel.Tracer(config.ServiceName)

	// productsInMemoryRepository := repositories.MustNewInMemoryProductsRepository(logger, tracer, 10)
	pgProductsRepository := repositories.MustNewPgProductsRepository(config.PostgresConnectionString)

	application := app.MustNewApplication(logger, tracer, pgProductsRepository)
	productsResolver := resolvers.MustNewProductsResolver(logger, tracer, application)

	server := grpc.MustNewServer(grpc.ServerInput{
		Port:   config.GrpcServerPort,
		Logger: logger,
		Registrator: func(server gGRPC.ServiceRegistrar) {
			protog.RegisterProductsServiceServer(server, productsResolver)
		},
	})

	go func() {
		port := fmt.Sprintf(":%d", config.MetricsPort)

		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(port, nil)
	}()

	if err := server.Run(ctx); err != nil {
		logger.Fatal("failed to run gRPC server", zap.Error(err))
	}
}
