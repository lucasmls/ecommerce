package app

import (
	"context"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.uber.org/zap"
)

func (a application) UpdateProduct(ctx context.Context, product domain.Product) (domain.Product, error) {
	ctx, span := a.Tracer.Start(ctx, "app.UpdateProduct")
	defer span.End()

	a.Logger.Info("updating a product", zap.Any("product", product))

	updatedProduct, err := a.ProductsRepository.Update(ctx, product)
	if err != nil {
		return domain.Product{}, err
	}

	return updatedProduct, nil
}
