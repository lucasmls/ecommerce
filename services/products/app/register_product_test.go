package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/suite"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/lucasmls/ecommerce/services/products/mocks"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type RegisterProductSuite struct {
	suite.Suite

	productsRepo *mocks.ProductsRepository
	app          domain.Application
}

func (s *RegisterProductSuite) SetupSuite() {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")
	s.productsRepo = &mocks.ProductsRepository{}

	s.app = NewApplication(loggerM, tracerM, s.productsRepo)
}

func (s *RegisterProductSuite) Test_RegisterProduct() {
	s.Run("Should fail when repository.Create returns any error", func() {
		ctx := context.Background()
		product := domain.Product{
			ID:          1,
			Name:        "Macbook Air M1",
			Description: "Fast!",
			Price:       6800,
		}

		s.productsRepo.On("Create",
			mock.AnythingOfType("*context.valueCtx"),
			product,
		).Return(domain.Product{}, errors.New("failed to register Product"))

		_, err := s.app.RegisterProduct(ctx, product)

		s.Equal(errors.New("failed to register Product"), err)
	})

	s.Run("Should register a new Product", func() {
		ctx := context.Background()
		product := domain.Product{
			ID:          2,
			Name:        "Macbook Air M1",
			Description: "Fast!",
			Price:       6800,
		}

		s.productsRepo.On("Create",
			mock.AnythingOfType("*context.valueCtx"),
			product,
		).Return(product, nil)

		got, err := s.app.RegisterProduct(ctx, product)

		s.NoError(err)
		s.Equal(product, got)
	})
}

func TestRegisterProductSuite(t *testing.T) {
	suite.Run(t, new(RegisterProductSuite))
}
