package repositories

import (
	"context"
	"math/rand"
	"testing"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func TestNewInMemoryProductsRepository(t *testing.T) {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name           string
			logger         *zap.Logger
			tracer         trace.Tracer
			storageSize    int
			expectedResult error
		}{
			{
				name:           "Should not be able to instantiate ProductsRepository with no storage capacity",
				logger:         loggerM,
				tracer:         tracerM,
				storageSize:    0,
				expectedResult: ErrInvalidStorageSize,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				_, err := NewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)
				assert.Equal(t, err, tc.expectedResult)
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name           string
			logger         *zap.Logger
			tracer         trace.Tracer
			storageSize    int
			expectedResult InMemoryProductsRepository
		}{
			{
				name:        "Should construct ProductsRepository with correct storage size",
				logger:      loggerM,
				tracer:      tracerM,
				storageSize: 10,
				expectedResult: InMemoryProductsRepository{
					Logger:      loggerM,
					Tracer:      tracerM,
					StorageSize: 10,
					storage:     map[int]domain.Product{},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				got, err := NewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)

				assert.NoError(t, err)
				assert.Equal(t, got, tc.expectedResult)
			})
		}
	})
}

func TestMustNewInMemoryProductsRepository(t *testing.T) {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name           string
			logger         *zap.Logger
			tracer         trace.Tracer
			storageSize    int
			expectedResult error
		}{
			{
				name:           "Should not be able to instantiate ProductsRepository with no storage capacity",
				logger:         loggerM,
				tracer:         tracerM,
				storageSize:    0,
				expectedResult: ErrInvalidStorageSize,
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				defer func() {
					err := recover()
					assert.Equal(t, err, tc.expectedResult)
				}()

				MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name           string
			logger         *zap.Logger
			tracer         trace.Tracer
			storageSize    int
			expectedResult InMemoryProductsRepository
		}{
			{
				name:        "Should construct ProductsRepository with correct storage size",
				logger:      loggerM,
				tracer:      tracerM,
				storageSize: 10,
				expectedResult: InMemoryProductsRepository{
					Logger:      loggerM,
					Tracer:      tracerM,
					StorageSize: 10,
					storage:     map[int]domain.Product{},
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				got := MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)
				assert.Equal(t, got, tc.expectedResult)
			})
		}
	})
}

func TestInMemoryProductsRepository_Create(t *testing.T) {
	rand.Seed(1) // Used to have some control over pseudo-random Product Id generation

	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name           string
			storageSize    int
			logger         *zap.Logger
			tracer         trace.Tracer
			ctx            context.Context
			product        domain.Product
			expectedResult error
			beforeEach     func(ctx context.Context, productsRepo domain.ProductsRepository)
		}{
			{
				name:        "Should not store a Product when the storage capacity is full",
				logger:      loggerM,
				tracer:      tracerM,
				storageSize: 1,
				ctx:         context.Background(),
				product: domain.Product{
					ID:          1,
					Name:        "Macbook Air M1",
					Description: "Fast!",
					Price:       7000,
				},
				expectedResult: ErrStorageLimitReached,
				beforeEach: func(ctx context.Context, productsRepo domain.ProductsRepository) {
					productsRepo.Create(ctx, domain.Product{
						ID:          1,
						Name:        "Iphone 12",
						Description: "Cool",
						Price:       4500,
					})
				},
			},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				productsRepo := MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)

				if tc.beforeEach != nil {
					tc.beforeEach(tc.ctx, productsRepo)
				}

				_, err := productsRepo.Create(tc.ctx, tc.product)
				assert.EqualError(t, err, tc.expectedResult.Error())
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name           string
			storageSize    int
			logger         *zap.Logger
			tracer         trace.Tracer
			ctx            context.Context
			product        domain.Product
			expectedResult domain.Product
		}{
			{
				name:        "Should store the given Product",
				logger:      loggerM,
				tracer:      tracerM,
				storageSize: 2,
				ctx:         context.Background(),
				product: domain.Product{
					ID:          1,
					Name:        "Macbook Air M1",
					Description: "Fast!",
					Price:       7000,
				},
				expectedResult: domain.Product{
					ID:          1,
					Name:        "Macbook Air M1",
					Description: "Fast!",
					Price:       7000,
				},
			},
			{
				name:        "Should store the given Product with a pseudo random id",
				logger:      loggerM,
				tracer:      tracerM,
				storageSize: 2,
				ctx:         context.Background(),
				product: domain.Product{
					Name:        "Macbook Air M1",
					Description: "Fast!",
					Price:       7000,
				},
				expectedResult: domain.Product{
					ID:          81,
					Name:        "Macbook Air M1",
					Description: "Fast!",
					Price:       7000,
				},
			},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				productsRepo := MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)

				got, err := productsRepo.Create(tc.ctx, tc.product)
				assert.NoError(t, err)
				assert.Equal(t, got, tc.expectedResult)
			})
		}
	})
}

func TestInMemoryProductsRepository_Update(t *testing.T) {
	loggerM, _ := zap.NewDevelopment()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tests := []struct {
			name           string
			storageSize    int
			logger         *zap.Logger
			tracer         trace.Tracer
			ctx            context.Context
			product        domain.Product
			expectedResult error
		}{
			{
				name:        "Should return not found error in case the specified product isn't stored",
				storageSize: 10,
				logger:      loggerM,
				tracer:      tracerM,
				ctx:         context.Background(),
				product: domain.Product{
					ID:          1,
					Name:        "Macbook Air M1",
					Description: "Fast!",
					Price:       7000,
				},
				expectedResult: domain.ErrProductNotFound,
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				productsRepo := MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)

				_, err := productsRepo.Update(tc.ctx, tc.product)
				assert.EqualError(t, err, tc.expectedResult.Error())
			})
		}
	})

	t.Run("Successful test", func(t *testing.T) {
		tests := []struct {
			name           string
			storageSize    int
			logger         *zap.Logger
			tracer         trace.Tracer
			ctx            context.Context
			product        domain.Product
			expectedResult domain.Product
			beforeEach     func(ctx context.Context, productsRepo domain.ProductsRepository)
		}{
			{
				name:        "Should update the provided Product",
				storageSize: 10,
				logger:      loggerM,
				tracer:      tracerM,
				ctx:         context.Background(),
				product: domain.Product{
					ID:          1,
					Name:        "Macbook Air M1",
					Description: "Fast!",
					Price:       7000,
				},
				expectedResult: domain.Product{
					ID:          1,
					Name:        "Macbook Air M1",
					Description: "Fast!",
					Price:       7000,
				},
				beforeEach: func(ctx context.Context, productsRepo domain.ProductsRepository) {
					productsRepo.Create(ctx, domain.Product{
						ID:          1,
						Name:        "Iphone 12",
						Description: "Cool",
						Price:       4500,
					})
				},
			},
		}
		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				productsRepo := MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)

				if tc.beforeEach != nil {
					tc.beforeEach(tc.ctx, productsRepo)
				}

				got, err := productsRepo.Update(tc.ctx, tc.product)
				assert.NoError(t, err)
				assert.Equal(t, got, tc.expectedResult)
			})
		}
	})
}

func TestInMemoryProductsRepository_Delete(t *testing.T) {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name           string
			logger         *zap.Logger
			tracer         trace.Tracer
			storageSize    int
			ctx            context.Context
			productId      int
			expectedResult error
		}{
			{
				name:           "Should return not found error in case the specified product isn't stored",
				logger:         loggerM,
				tracer:         tracerM,
				storageSize:    1,
				ctx:            context.Background(),
				productId:      1,
				expectedResult: domain.ErrProductNotFound,
			},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				productsRepo := MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)

				err := productsRepo.Delete(tc.ctx, tc.productId)
				assert.EqualError(t, err, tc.expectedResult.Error())
			})
		}
	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name        string
			logger      *zap.Logger
			tracer      trace.Tracer
			storageSize int
			ctx         context.Context
			productId   int
			beforeEach  func(ctx context.Context, productsRepo domain.ProductsRepository, productId int)
		}{
			{
				name:        "Should remove the specified Product",
				logger:      loggerM,
				tracer:      tracerM,
				storageSize: 2,
				ctx:         context.Background(),
				productId:   1,
				beforeEach: func(ctx context.Context, productsRepo domain.ProductsRepository, productId int) {
					productsRepo.Create(ctx, domain.Product{
						ID:          productId,
						Name:        "Macbook Air M1",
						Description: "Fast!",
						Price:       7000,
					})
				},
			},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				productsRepo := MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)

				if tc.beforeEach != nil {
					tc.beforeEach(tc.ctx, productsRepo, tc.productId)
				}

				err := productsRepo.Delete(tc.ctx, tc.productId)
				assert.NoError(t, err)
			})
		}
	})
}

func TestInMemoryProductsRepository_List(t *testing.T) {
	loggerM, _ := zap.NewDevelopment()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	seedProducts := func(ctx context.Context, productsRepo domain.ProductsRepository) {
		productsRepo.Create(ctx, domain.Product{
			ID:          1,
			Name:        "Iphone 13",
			Description: "Cool",
			Price:       4500,
		})
		productsRepo.Create(ctx, domain.Product{
			ID:          2,
			Name:        "Macbook Pro M1 Max",
			Description: "Fast!",
			Price:       16500,
		})
		productsRepo.Create(ctx, domain.Product{
			ID:          3,
			Name:        "Macbook Air M1",
			Description: "Nice!",
			Price:       6900,
		})
	}

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name           string
			storageSize    int
			logger         *zap.Logger
			tracer         trace.Tracer
			ctx            context.Context
			filter         domain.ListProductsFilter
			expectedResult []domain.Product
			beforeEach     func(ctx context.Context, productsRepo domain.ProductsRepository)
		}{
			{
				name:        "Should list all products that were stored",
				logger:      loggerM,
				tracer:      tracerM,
				storageSize: 10,
				ctx:         context.Background(),
				filter: domain.ListProductsFilter{
					IDs: []int{},
				},
				expectedResult: []domain.Product{
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
				},
				beforeEach: seedProducts,
			},
			{
				name:        "Should list only the products that match the provided filter",
				logger:      loggerM,
				tracer:      tracerM,
				storageSize: 10,
				ctx:         context.Background(),
				filter: domain.ListProductsFilter{
					IDs: []int{
						1,
					},
				},
				expectedResult: []domain.Product{
					{
						ID:          1,
						Name:        "Iphone 13",
						Description: "Cool",
						Price:       4500,
					},
				},
				beforeEach: seedProducts,
			},
		}
		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				productsRepo := MustNewInMemoryProductsRepository(tc.logger, tc.tracer, tc.storageSize)

				if tc.beforeEach != nil {
					tc.beforeEach(tc.ctx, productsRepo)
				}

				got, err := productsRepo.List(tc.ctx, tc.filter)

				assert.NoError(t, err)
				assert.ElementsMatch(t, got, tc.expectedResult)
			})
		}
	})
}
