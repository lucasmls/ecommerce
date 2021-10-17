package app

import (
	"context"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.uber.org/zap"
)

func (a application) UpdateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	a.logger.Info("updating a product", zap.Any("product", product))

	updatedProduct, err := a.productsRepository.Update(ctx, product)
	if err != nil {
		return domain.Product{}, err
	}

	return updatedProduct, nil
}
