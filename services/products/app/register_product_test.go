package app

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/lucasmls/ecommerce/services/products/mocks"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Test_RegisterProduct(t *testing.T) {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name                string
			logger              *zap.Logger
			tracer              trace.Tracer
			product             domain.Product
			expectedResult      error
			productsRepositoryM func(context.Context, *gomock.Controller, domain.Product) domain.ProductsRepository
		}{
			{
				name:   "Should fail when repository.Create returns a error",
				logger: loggerM,
				tracer: tracerM,
				product: domain.Product{
					ID:          1,
					Name:        "Iphone 13",
					Description: "Cool",
					Price:       4500,
				},
				expectedResult: errors.New("Failed to register Product"),
				productsRepositoryM: func(c context.Context, g *gomock.Controller, product domain.Product) domain.ProductsRepository {
					productsRepoM := mocks.NewMockProductsRepository(g)

					productsRepoM.EXPECT().
						Create(gomock.Any(), product).
						Return(domain.Product{}, errors.New("Failed to register Product"))

					return productsRepoM
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				ctrl := gomock.NewController(t)

				a := MustNewApplication(
					tc.logger,
					tc.tracer,
					tc.productsRepositoryM(ctx, ctrl, tc.product),
				)

				_, err := a.RegisterProduct(ctx, tc.product)
				assert.EqualError(t, err, tc.expectedResult.Error())
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name                string
			logger              *zap.Logger
			tracer              trace.Tracer
			product             domain.Product
			expectedResult      domain.Product
			productsRepositoryM func(context.Context, *gomock.Controller, domain.Product) domain.ProductsRepository
		}{
			{
				name:   "Should register a new Product",
				logger: loggerM,
				tracer: tracerM,
				product: domain.Product{
					ID:          1,
					Name:        "Iphone 13",
					Description: "Cool",
					Price:       4500,
				},
				expectedResult: domain.Product{
					ID:          1,
					Name:        "Iphone 13",
					Description: "Cool",
					Price:       4500,
				},
				productsRepositoryM: func(c context.Context, g *gomock.Controller, product domain.Product) domain.ProductsRepository {
					productsRepoM := mocks.NewMockProductsRepository(g)

					productsRepoM.EXPECT().
						Create(gomock.Any(), product).
						Return(product, nil)

					return productsRepoM
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				ctrl := gomock.NewController(t)

				a := MustNewApplication(
					tc.logger,
					tc.tracer,
					tc.productsRepositoryM(ctx, ctrl, tc.product),
				)

				got, err := a.RegisterProduct(ctx, tc.product)
				assert.NoError(t, err)
				assert.Equal(t, got, tc.expectedResult)
			})
		}
	})
}
