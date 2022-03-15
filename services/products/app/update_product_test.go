package app

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/lucasmls/ecommerce/services/products/mocks"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func Test_UpdateProduct(t *testing.T) {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name                string
			logger              *zap.Logger
			tracer              trace.Tracer
			product             domain.Product
			expectedResult      error
			productsRepositoryM func(ctx context.Context, ctrl *gomock.Controller, product domain.Product) domain.ProductsRepository
		}{
			{
				name:           "Should return not found error in case the specified product isn't stored",
				expectedResult: domain.ErrProductNotFound,
				logger:         loggerM,
				tracer:         tracerM,
				product: domain.Product{
					ID:          123,
					Name:        "Macbook Air M1",
					Description: "Cool!",
					Price:       6800,
				},
				productsRepositoryM: func(ctx context.Context, ctrl *gomock.Controller, product domain.Product) domain.ProductsRepository {
					productsRepoM := mocks.NewMockProductsRepository(ctrl)

					productsRepoM.EXPECT().
						Update(gomock.Any(), product).
						Return(domain.Product{}, domain.ErrProductNotFound)

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

				_, err := a.UpdateProduct(ctx, tc.product)
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
			productsRepositoryM func(ctx context.Context, ctrl *gomock.Controller, product domain.Product) domain.ProductsRepository
		}{
			{
				name:   "Should return the updated product",
				logger: loggerM,
				tracer: tracerM,
				product: domain.Product{
					ID:          123,
					Name:        "Macbook Air M1",
					Description: "Cool!",
					Price:       6800,
				},
				expectedResult: domain.Product{
					ID:          123,
					Name:        "Macbook Air M1",
					Description: "Cool!",
					Price:       6800,
				},
				productsRepositoryM: func(ctx context.Context, ctrl *gomock.Controller, product domain.Product) domain.ProductsRepository {
					productsRepoM := mocks.NewMockProductsRepository(ctrl)

					productsRepoM.EXPECT().
						Update(gomock.Any(), product).
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

				got, err := a.UpdateProduct(ctx, tc.product)
				assert.NoError(t, err)
				assert.Equal(t, got, tc.expectedResult)
			})
		}
	})
}
