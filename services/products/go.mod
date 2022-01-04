module github.com/lucasmls/ecommerce/services/products

go 1.17

replace github.com/lucasmls/ecommerce/shared => ../../shared

require (
	github.com/golang/mock v1.6.0
	github.com/lucasmls/ecommerce/shared v0.0.0-20211019010026-2ed6e2591d9f
	github.com/streadway/amqp v1.0.0
	go.opentelemetry.io/otel v1.2.0
	go.opentelemetry.io/otel/exporters/jaeger v1.2.0
	go.opentelemetry.io/otel/sdk v1.2.0
	go.opentelemetry.io/otel/trace v1.2.0
	go.uber.org/zap v1.19.1
	google.golang.org/grpc v1.42.0
	google.golang.org/protobuf v1.27.1
)

require (
	cloud.google.com/go v0.93.3 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.27.0 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210503060351-7fd8e65b6420 // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	golang.org/x/sys v0.0.0-20210823070655-63515b42dcdf // indirect
	golang.org/x/text v0.3.6 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20210828152312-66f60bf46e71 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
