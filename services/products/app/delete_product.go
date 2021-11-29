package app

import (
	"context"

	"go.uber.org/zap"
)

func (a application) DeleteProduct(ctx context.Context, id string) error {

	a.in.Logger.Info("deleting a product", zap.Any("id", id))

	err := a.in.ProductsRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
