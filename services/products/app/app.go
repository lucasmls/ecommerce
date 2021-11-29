package app

import (
	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ApplicationInput holds all the dependencies needed to
// instantiate the Application
type ApplicationInput struct {
	Logger *zap.Logger

	ProductsRepository domain.ProductsRepository
}

type application struct {
	in ApplicationInput
}

// NewApplication creates a new Application instance
func NewApplication(in ApplicationInput) (application, error) {
	return application{
		in: in,
	}, nil
}

// MustNewApplication creates a new Application instance
// It panics if any error is found
func MustNewApplication(in ApplicationInput) application {
	app, err := NewApplication(in)
	if err != nil {
		panic(err)
	}

	return app
}
