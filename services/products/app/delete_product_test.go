package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/lucasmls/ecommerce/services/products/mocks"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type DeleteProductSuite struct {
	suite.Suite

	productsRepo *mocks.ProductsRepository
	app          domain.Application
}

func (s *DeleteProductSuite) SetupSuite() {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")
	s.productsRepo = &mocks.ProductsRepository{}

	s.app = NewApplication(loggerM, tracerM, s.productsRepo)
}

func (s *DeleteProductSuite) Test_DeleteProduct() {
	s.Run("Should return not found error in case the specified product isn't stored", func() {
		ctx := context.Background()
		productId := 1

		s.productsRepo.
			On("Delete",
				mock.AnythingOfType("*context.valueCtx"),
				productId,
			).
			Return(domain.ErrProductNotFound)

		err := s.app.DeleteProduct(ctx, productId)

		s.Equal(domain.ErrProductNotFound, err)
	})

	s.Run("Should delete the specified product", func() {
		ctx := context.Background()
		productId := 2

		s.productsRepo.
			On("Delete",
				mock.AnythingOfType("*context.valueCtx"),
				productId,
			).
			Return(nil)

		err := s.app.DeleteProduct(ctx, productId)

		s.NoError(err)
	})
}

func TestDeleteProductSuite(t *testing.T) {
	suite.Run(t, new(DeleteProductSuite))
}
