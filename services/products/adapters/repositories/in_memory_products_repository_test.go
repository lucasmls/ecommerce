package repositories

import (
	"context"
	"math/rand"
	"testing"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/stretchr/testify/suite"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type NewInMemoryProductsRepositorySuite struct {
	suite.Suite

	loggerM *zap.Logger
	tracerM trace.Tracer
}

func (s *NewInMemoryProductsRepositorySuite) SetupSuite() {
	s.loggerM = zap.NewNop()
	s.tracerM = trace.NewNoopTracerProvider().Tracer("")
}

func (s *NewInMemoryProductsRepositorySuite) Test_NewInMemoryProductsRepository() {
	s.Run("Should not be able to instantiate ProductsRepository with no storage capacity", func() {
		_, err := NewInMemoryProductsRepository(s.loggerM, s.tracerM, 0)

		s.Error(err)
		s.Equal(err, ErrInvalidStorageSize)
	})

	s.Run("Should construct ProductsRepository with correct storage size", func() {
		expectedResult := InMemoryProductsRepository{
			Logger:      s.loggerM,
			Tracer:      s.tracerM,
			StorageSize: 10,
			storage:     map[int]domain.Product{},
		}

		got, err := NewInMemoryProductsRepository(s.loggerM, s.tracerM, 10)

		s.NoError(err)
		s.Equal(expectedResult, got)
	})
}

func (s *NewInMemoryProductsRepositorySuite) Test_MustNewInMemoryProductsRepository() {
	s.Run("Should not be able to instantiate ProductsRepository with no storage capacity", func() {
		defer func() {
			err := recover()
			s.NotNil(err)
			s.Equal(ErrInvalidStorageSize, err)
		}()

		MustNewInMemoryProductsRepository(s.loggerM, s.tracerM, 0)
	})

	s.Run("Should construct ProductsRepository with correct storage size", func() {
		expectedResult := InMemoryProductsRepository{
			Logger:      s.loggerM,
			Tracer:      s.tracerM,
			StorageSize: 10,
			storage:     map[int]domain.Product{},
		}

		got := MustNewInMemoryProductsRepository(s.loggerM, s.tracerM, 10)

		s.Equal(expectedResult, got)
	})
}

type CreateSuite struct {
	suite.Suite

	loggerM      *zap.Logger
	tracerM      trace.Tracer
	productsRepo domain.ProductsRepository
}

func (s *CreateSuite) SetupSuite() {
	s.loggerM = zap.NewNop()
	s.tracerM = trace.NewNoopTracerProvider().Tracer("")
	s.productsRepo = MustNewInMemoryProductsRepository(s.loggerM, s.tracerM, 2)
}

func (s *CreateSuite) Test_Create() {
	s.Run("Should store the given Product", func() {
		product := domain.Product{
			ID:          1,
			Name:        "Macbook Air M1",
			Description: "Fast!",
			Price:       7000,
		}

		createInput := product
		expectedResult := product

		ctx := context.Background()
		got, err := s.productsRepo.Create(ctx, createInput)

		s.NoError(err)
		s.Equal(expectedResult, got)
	})

	s.Run("Should store the given Product with a pseudo random id", func() {
		rand.Seed(1) // Used to have some control over pseudo-random ProductId generation
		expectedResult := domain.Product{
			ID:          81,
			Name:        "Macbook Air M1",
			Description: "Fast!",
			Price:       7000,
		}

		ctx := context.Background()
		createInput := domain.Product{
			Name:        "Macbook Air M1",
			Description: "Fast!",
			Price:       7000,
		}
		got, err := s.productsRepo.Create(ctx, createInput)

		s.NoError(err)
		s.Equal(expectedResult, got)
	})

	s.Run("Should not be able to store a Product when the storage capacity is full", func() {
		ctx := context.Background()
		createInput := domain.Product{
			Name:        "Macbook Air M1",
			Description: "Fast!",
			Price:       7000,
		}

		_, err := s.productsRepo.Create(ctx, createInput)
		s.Equal(ErrStorageLimitReached, err)
	})
}

type UpdateSuite struct {
	suite.Suite

	loggerM      *zap.Logger
	tracerM      trace.Tracer
	productsRepo domain.ProductsRepository
}

func (s *UpdateSuite) SetupSuite() {
	s.loggerM = zap.NewNop()
	s.tracerM = trace.NewNoopTracerProvider().Tracer("")
	s.productsRepo = MustNewInMemoryProductsRepository(s.loggerM, s.tracerM, 10)
}

func (s *UpdateSuite) SetupTest() {
	_, err := s.productsRepo.Create(context.Background(), domain.Product{
		ID:          1,
		Name:        "Iphone 12",
		Description: "Cool",
		Price:       4500,
	})

	s.NoError(err)
}

func (s *UpdateSuite) Test_Update() {
	s.Run("Should return not found error in case the specified product isn't stored", func() {
		ctx := context.Background()

		_, err := s.productsRepo.Update(ctx, domain.Product{ID: 2})
		s.Equal(domain.ErrProductNotFound, err)
	})

	s.Run("Should update the provided Product", func() {
		product := domain.Product{
			ID:          1,
			Name:        "Macbook Air M1",
			Description: "Fast!",
			Price:       7000,
		}

		expectedResult := product

		ctx := context.Background()
		updateInput := product

		got, err := s.productsRepo.Update(ctx, updateInput)

		s.NoError(err)
		s.Equal(expectedResult, got)
	})
}

type ListSuite struct {
	suite.Suite

	loggerM      *zap.Logger
	tracerM      trace.Tracer
	productsRepo domain.ProductsRepository
}

func (s *ListSuite) SetupSuite() {
	s.loggerM = zap.NewNop()
	s.tracerM = trace.NewNoopTracerProvider().Tracer("")
	s.productsRepo = MustNewInMemoryProductsRepository(s.loggerM, s.tracerM, 3)
}

func (s *ListSuite) SetupTest() {
	ctx := context.Background()

	products := []domain.Product{
		{1, "Iphone 13", "Cool", 4500},
		{2, "Macbook Pro M1 Max", "Fast!", 16500},
		{3, "Macbook Air M1", "Nice!", 6900},
	}

	for _, product := range products {
		_, err := s.productsRepo.Create(ctx, domain.Product{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
		s.NoError(err)
	}
}

func (s *ListSuite) Test_List() {
	s.Run("Should list all products that were stored", func() {
		expectedResult := []domain.Product{
			{1, "Iphone 13", "Cool", 4500},
			{2, "Macbook Pro M1 Max", "Fast!", 16500},
			{3, "Macbook Air M1", "Nice!", 6900},
		}

		ctx := context.Background()
		filterInput := domain.ListProductsFilter{IDs: []int{}}

		got, err := s.productsRepo.List(ctx, filterInput)

		s.NoError(err)
		s.ElementsMatch(expectedResult, got)
	})

	s.Run("Should list only the products that matches the provided filter", func() {
		expectedResult := []domain.Product{
			{1, "Iphone 13", "Cool", 4500},
		}

		ctx := context.Background()
		filterInput := domain.ListProductsFilter{IDs: []int{1}}

		got, err := s.productsRepo.List(ctx, filterInput)

		s.NoError(err)
		s.ElementsMatch(expectedResult, got)
	})
}

type DeleteSuite struct {
	suite.Suite

	loggerM      *zap.Logger
	tracerM      trace.Tracer
	productsRepo domain.ProductsRepository
}

func (s *DeleteSuite) SetupSuite() {
	s.loggerM = zap.NewNop()
	s.tracerM = trace.NewNoopTracerProvider().Tracer("")
	s.productsRepo = MustNewInMemoryProductsRepository(s.loggerM, s.tracerM, 10)
}

func (s *DeleteSuite) SetupTest() {
	_, err := s.productsRepo.Create(context.Background(), domain.Product{
		ID:          100,
		Name:        "Macbook Air M1",
		Description: "Fast!",
		Price:       7000,
	})

	s.NoError(err)
}

func (s *DeleteSuite) Test_Delete() {
	s.Run("Should return not found error in case the specified product isn't stored", func() {
		ctx := context.Background()

		err := s.productsRepo.Delete(ctx, 1)
		s.ErrorIs(err, domain.ErrProductNotFound)
	})

	s.Run("Should remove the specified Product", func() {
		ctx := context.Background()

		err := s.productsRepo.Delete(ctx, 100)
		s.NoError(err)
	})
}

func TestInMemoryProductsRepositorySuites(t *testing.T) {
	suite.Run(t, new(NewInMemoryProductsRepositorySuite))
	suite.Run(t, new(CreateSuite))
	suite.Run(t, new(ListSuite))
	suite.Run(t, new(UpdateSuite))
	suite.Run(t, new(DeleteSuite))
}
