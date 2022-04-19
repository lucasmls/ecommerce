package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/lucasmls/ecommerce/services/products/mocks"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type ListProductsSuite struct {
	suite.Suite

	productsRepo *mocks.ProductsRepository
	app          domain.Application
}

func (s *ListProductsSuite) SetupSuite() {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")
	s.productsRepo = &mocks.ProductsRepository{}

	s.app = NewApplication(loggerM, tracerM, s.productsRepo)
}

func (s *ListProductsSuite) Test_ListProducts() {
	s.Run("Should fail when repository.List returns any error", func() {
		ctx := context.Background()
		filter := domain.ListProductsFilter{}

		s.productsRepo.
			On("List",
				mock.AnythingOfType("*context.valueCtx"),
				filter,
			).
			Return([]domain.Product{}, errors.New("failed to list products from the datastore"))

		_, err := s.app.ListProducts(ctx, filter)

		s.Equal(errors.New("failed to list products from the datastore"), err)
	})

	s.Run("Should return the Products returned by the repository layer", func() {
		products := []domain.Product{
			{
				ID:          1,
				Name:        "Iphone 13",
				Description: "Cool",
				Price:       4500,
			},
			{
				ID:          2,
				Name:        "Macbook Pro M1 Max",
				Description: "Fast!",
				Price:       16500,
			},
		}

		ctx := context.Background()
		filter := domain.ListProductsFilter{IDs: []int{1, 2}}

		s.productsRepo.On("List",
			mock.AnythingOfType("*context.valueCtx"),
			filter,
		).
			Return(products, nil)

		got, err := s.app.ListProducts(ctx, filter)

		s.NoError(err)
		s.ElementsMatch(products, got)
	})
}

func TestListProductsSuite(t *testing.T) {
	suite.Run(t, new(ListProductsSuite))
}
