package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/lucasmls/ecommerce/services/products/mocks"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type UpdateProductSuite struct {
	suite.Suite

	productsRepo *mocks.ProductsRepository
	app          domain.Application
}

func (s *UpdateProductSuite) SetupSuite() {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")
	s.productsRepo = &mocks.ProductsRepository{}

	s.app = MustNewApplication(loggerM, tracerM, s.productsRepo)
}

func (s *UpdateProductSuite) Test_UpdateProduct() {
	s.Run("Should return not found error in case the specified product isn't stored", func() {
		ctx := context.Background()
		product := domain.Product{
			ID:          1,
			Name:        "Macbook Air M1",
			Description: "Fast!",
			Price:       6800,
		}

		s.productsRepo.On("Update",
			mock.AnythingOfType("*context.valueCtx"),
			product,
		).Return(domain.Product{}, domain.ErrProductNotFound)

		_, err := s.app.UpdateProduct(ctx, product)

		s.Equal(domain.ErrProductNotFound, err)
	})

	s.Run("Should return the updated product", func() {
		ctx := context.Background()
		product := domain.Product{
			ID:          2,
			Name:        "Macbook Air M1 - Updated",
			Description: "Fast!!!",
			Price:       6800,
		}

		s.productsRepo.On("Update",
			mock.AnythingOfType("*context.valueCtx"),
			product,
		).Return(product, nil)

		got, err := s.app.UpdateProduct(ctx, product)

		s.NoError(err)
		s.Equal(product, got)
	})
}

func TestUpdateProductSuite(t *testing.T) {
	suite.Run(t, new(UpdateProductSuite))
}
