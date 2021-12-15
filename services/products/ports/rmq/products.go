package rmq_port

import (
	"context"
	"errors"

	"github.com/lucasmls/ecommerce/services/products/domain"
	protoMessages "github.com/lucasmls/ecommerce/services/products/ports/rmq/proto"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ProductsConsumerInput ...
type ProductsConsumerInput struct {
	Logger *zap.Logger
	App    domain.Application
	Tracer trace.Tracer
}

// ProductsConsumer ...
type ProductsConsumer struct {
	in ProductsConsumerInput
}

// NewProductsConsumer creates a new ProductsConsumer instance
func NewProductsConsumer(in ProductsConsumerInput) (*ProductsConsumer, error) {
	if in.Logger == nil {
		return nil, errors.New("missing required dependency: Logger")
	}

	return &ProductsConsumer{
		in: in,
	}, nil
}

// NewProductsConsumer creates a new ProductsConsumer instance
// It panics if any error is found
func MustNewProductsConsumer(in ProductsConsumerInput) *ProductsConsumer {
	app, err := NewProductsConsumer(in)
	if err != nil {
		panic(err)
	}

	return app
}

func (r *ProductsConsumer) Register(ctx context.Context, req *protoMessages.Product) {
	ctx, span := r.in.Tracer.Start(ctx, "consumer.Register")
	defer span.End()

	r.in.Logger.Info("registering a new product", zap.Any("product", req))

	product, err := r.in.App.RegisterProduct(ctx, domain.Product{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		r.in.Logger.Error("failed to register a new product", zap.Error(err))
	}

	r.in.Logger.Info("product registered successfully", zap.Any("product", product))
}
