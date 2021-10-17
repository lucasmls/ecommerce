package app

import (
	"context"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.uber.org/zap"
)

func (a application) RegisterProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	a.logger.Info("registering a new product", zap.Any("product", product))

	registeredProduct, err := a.productsRepository.Create(ctx, product)
	if err != nil {
		return domain.Product{}, err
	}

	return registeredProduct, nil
}
