package app

import (
	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// application holds all dependencies the Application needs to work
type application struct {
	Logger *zap.Logger
	Tracer trace.Tracer

	ProductsRepository domain.ProductsRepository
}

// NewApplication creates a new Application instance
func NewApplication(
	logger *zap.Logger,
	tracer trace.Tracer,
	productsRepository domain.ProductsRepository,
) application {
	return application{
		Logger:             logger,
		Tracer:             tracer,
		ProductsRepository: productsRepository,
	}
}

// MustNewApplication creates a new Application instance
// It panics if any error is found
func MustNewApplication(
	logger *zap.Logger,
	tracer trace.Tracer,
	productsRepository domain.ProductsRepository,
) application {
	app := NewApplication(logger, tracer, productsRepository)
	return app
}
