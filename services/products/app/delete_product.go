package app

import (
	"context"

	"go.uber.org/zap"
)

func (a application) DeleteProduct(ctx context.Context, id int) error {
	ctx, span := a.Tracer.Start(ctx, "app.DeleteProduct")
	defer span.End()

	a.Logger.Info("deleting a product", zap.Any("id", id))

	err := a.ProductsRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
