package app

import (
	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.uber.org/zap"
)

type application struct {
	logger             *zap.Logger
	productsRepository domain.ProductsRepository
}

// NewApplication creates a new Application instance
func NewApplication(
	logger *zap.Logger,
	productsRepository domain.ProductsRepository,
) (application, error) {
	return application{
		logger:             logger,
		productsRepository: productsRepository,
	}, nil
}

// MustNewApplication creates a new Application instance
// It panics if any error is found
func MustNewApplication(
	logger *zap.Logger,
	productsRepository domain.ProductsRepository,
) application {
	app, err := NewApplication(logger, productsRepository)
	if err != nil {
		panic(err)
	}

	return app
}
