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

func TestListProducts(t *testing.T) {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

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
		{
			ID:          3,
			Name:        "Macbook Air M1",
			Description: "Nice!",
			Price:       6900,
		},
	}

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name                string
			logger              *zap.Logger
			tracer              trace.Tracer
			expectedResult      error
			productsRepositoryM func(context.Context, *gomock.Controller) domain.ProductsRepository
		}{
			{
				name:           "Should fail when repository.List returns a error",
				logger:         loggerM,
				tracer:         tracerM,
				expectedResult: errors.New("failed to list products from the datastore"),
				productsRepositoryM: func(c context.Context, g *gomock.Controller) domain.ProductsRepository {
					productsRepoM := mocks.NewMockProductsRepository(g)

					productsRepoM.EXPECT().
						List(gomock.Any(), domain.ListProductsFilter{}).
						Return([]domain.Product{}, errors.New("failed to list products from the datastore"))

					return productsRepoM
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				ctrl := gomock.NewController(t)

				a := MustNewApplication(tc.logger, tc.tracer, tc.productsRepositoryM(ctx, ctrl))

				_, err := a.ListProducts(ctx, domain.ListProductsFilter{})
				assert.EqualError(t, err, tc.expectedResult.Error())
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name                string
			logger              *zap.Logger
			tracer              trace.Tracer
			expectedResult      []domain.Product
			productsRepositoryM func(context.Context, *gomock.Controller) domain.ProductsRepository
		}{
			{
				name:           "Should return the Products returned by the repository layer",
				logger:         loggerM,
				tracer:         tracerM,
				expectedResult: products,
				productsRepositoryM: func(c context.Context, g *gomock.Controller) domain.ProductsRepository {
					productsRepoM := mocks.NewMockProductsRepository(g)

					productsRepoM.EXPECT().
						List(gomock.Any(), domain.ListProductsFilter{}).
						Return(products, nil)

					return productsRepoM
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				ctrl := gomock.NewController(t)

				a := MustNewApplication(tc.logger, tc.tracer, tc.productsRepositoryM(ctx, ctrl))

				got, err := a.ListProducts(ctx, domain.ListProductsFilter{})
				assert.NoError(t, err)
				assert.Equal(t, got, tc.expectedResult)
			})
		}
	})
}
