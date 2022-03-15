package main

import (
	"context"
	"fmt"
	"log"

	"github.com/lucasmls/ecommerce/services/products/adapters/repositories"
	"github.com/lucasmls/ecommerce/services/products/app"
	rmqPort "github.com/lucasmls/ecommerce/services/products/ports/rmq"
	protoMessages "github.com/lucasmls/ecommerce/services/products/ports/rmq/proto"
	"github.com/streadway/amqp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	jaegerExporter "go.opentelemetry.io/otel/exporters/jaeger"
	tracingSdkResource "go.opentelemetry.io/otel/sdk/resource"
	tracingSdk "go.opentelemetry.io/otel/sdk/trace"
)

const (
	productsExchange        string = "products"
	registerProductsRouting string = "register"
	amqpConnectionString    string = "amqp:guest:guest@localhost:5672/"
	jaegerEndpoint          string = "http://localhost:14268/api/traces"
)

func main() {
	ctx := context.Background()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	jaegerExporter, err := jaegerExporter.New(
		jaegerExporter.WithCollectorEndpoint(
			jaegerExporter.WithEndpoint(jaegerEndpoint),
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

	productsInMemoryRepository := repositories.MustNewInMemoryProductsRepository(logger, tracer, 10)

	application := app.MustNewApplication(app.ApplicationInput{
		Logger:             logger,
		Tracer:             tracer,
		ProductsRepository: productsInMemoryRepository,
	})

	rmqProductsConsumer := rmqPort.MustNewProductsConsumer(rmqPort.ProductsConsumerInput{
		Logger: logger,
		Tracer: tracer,
		App:    application,
	})

	ampqConnection, err := amqp.Dial(amqpConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer ampqConnection.Close()

	amqpChannel, err := ampqConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = amqpChannel.ExchangeDeclare(
		productsExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	productsQueue, err := amqpChannel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	err = amqpChannel.QueueBind(
		productsQueue.Name,
		registerProductsRouting,
		productsExchange,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	messages, err := amqpChannel.Consume(
		productsQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	for message := range messages {
		p := protoMessages.Product{}

		err := proto.Unmarshal(message.Body, &p)
		if err != nil {
			fmt.Println(err)
			continue
		}

		rmqProductsConsumer.Register(ctx, &p)
	}
}
