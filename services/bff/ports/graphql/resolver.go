package graph

import (
	productsPb "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	Logger          *zap.Logger
	ProductsService productsPb.ProductsServiceClient
}
