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

func Test_DeleteProduct(t *testing.T) {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name                string
			logger              *zap.Logger
			tracer              trace.Tracer
			productId           int
			expectedResult      error
			productsRepositoryM func(ctx context.Context, ctrl *gomock.Controller, productId int) domain.ProductsRepository
		}{
			{
				name:           "Should return not found error in case the specified product isn't stored",
				expectedResult: domain.ErrProductNotFound,
				logger:         loggerM,
				tracer:         tracerM,
				productId:      123,
				productsRepositoryM: func(ctx context.Context, ctrl *gomock.Controller, productId int) domain.ProductsRepository {
					productsRepoM := mocks.NewMockProductsRepository(ctrl)

					productsRepoM.EXPECT().
						Delete(gomock.Any(), productId).
						Return(domain.ErrProductNotFound)

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
					tc.productsRepositoryM(ctx, ctrl, tc.productId),
				)

				err := a.DeleteProduct(ctx, tc.productId)
				assert.EqualError(t, err, tc.expectedResult.Error())
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name                string
			logger              *zap.Logger
			tracer              trace.Tracer
			productId           int
			expectedResult      error
			productsRepositoryM func(ctx context.Context, ctrl *gomock.Controller, productId int) domain.ProductsRepository
		}{
			{
				name:           "Should delete the specified product",
				logger:         loggerM,
				tracer:         tracerM,
				expectedResult: nil,
				productId:      123,
				productsRepositoryM: func(ctx context.Context, ctrl *gomock.Controller, productId int) domain.ProductsRepository {
					productsRepoM := mocks.NewMockProductsRepository(ctrl)

					productsRepoM.EXPECT().
						Delete(gomock.Any(), productId).
						Return(nil)

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
					tc.productsRepositoryM(ctx, ctrl, tc.productId),
				)

				err := a.DeleteProduct(ctx, tc.productId)
				assert.NoError(t, err, tc.expectedResult)
			})
		}
	})
}
