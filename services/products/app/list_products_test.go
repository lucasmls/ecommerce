package app

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/lucasmls/ecommerce/services/products/mocks"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func TestListProducts(t *testing.T) {
	tt := []struct {
		name                string
		expectedResult      error
		productsRepositoryM func(context.Context, *gomock.Controller) domain.ProductsRepository
	}{
		{
			name:           "Fails when repository.List returns a error",
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
			logger, _ := zap.NewDevelopment()

			a := application{
				in: ApplicationInput{
					Logger:             logger,
					ProductsRepository: tc.productsRepositoryM(ctx, ctrl),
					Tracer:             trace.NewNoopTracerProvider().Tracer(""),
				},
			}

			_, err := a.ListProducts(ctx, domain.ListProductsFilter{})
			if err.Error() != tc.expectedResult.Error() {
				t.Errorf("application.ListProducts() = %v, want %v", err, tc.expectedResult)
			}
		})
	}
}
