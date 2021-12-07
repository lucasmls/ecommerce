package graph

import (
	productsPb "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Logger          *zap.Logger
	Tracer          trace.Tracer
	ProductsService productsPb.ProductsServiceClient
}
