package app

import (
	"context"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.uber.org/zap"
)

func (a application) ListProducts(ctx context.Context, filter domain.ListProductsFilter) ([]domain.Product, error) {
	a.logger.Info("listing products", zap.Any("filter", filter))

	products, err := a.productsRepository.List(ctx, filter)
	if err != nil {
		return nil, err
	}

	return products, nil
}
