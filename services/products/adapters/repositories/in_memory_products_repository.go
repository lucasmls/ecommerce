package repositories

import (
	"context"
	"errors"
	"math/rand"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var (
	ErrInvalidStorageSize  = errors.New("invalid-storage-size")
	ErrStorageLimitReached = errors.New("storage-limit-reached")
	ErrProductNotFound     = errors.New("product-not-found")
)

type InMemoryProductsRepository struct {
	Logger      *zap.Logger
	Tracer      trace.Tracer
	StorageSize int

	storage map[int]domain.Product
}

// NewInMemoryProductsRepository creates a new InMemoryProductsRepository.
func NewInMemoryProductsRepository(
	logger *zap.Logger,
	tracer trace.Tracer,
	storageSize int,
) (domain.ProductsRepository, error) {
	if storageSize == 0 {
		return InMemoryProductsRepository{}, ErrInvalidStorageSize
	}

	return InMemoryProductsRepository{
		Logger:      logger,
		Tracer:      tracer,
		StorageSize: storageSize,
		storage:     make(map[int]domain.Product, storageSize),
	}, nil
}

// MustNewInMemoryProductsRepository creates a new InMemoryProductsRepository.
// It panics if any error is found.
func MustNewInMemoryProductsRepository(
	logger *zap.Logger,
	tracer trace.Tracer,
	storageSize int,
) domain.ProductsRepository {
	repo, err := NewInMemoryProductsRepository(logger, tracer, storageSize)
	if err != nil {
		panic(err)
	}

	return repo
}

// Create creates a new Product in-memory.
func (r InMemoryProductsRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	_, span := r.Tracer.Start(ctx, "repository.Create")
	defer span.End()

	if len(r.storage) == r.StorageSize {
		return domain.Product{}, ErrStorageLimitReached
	}

	if product.ID == 0 {
		product.ID = rand.Intn(1000)
	}

	r.storage[product.ID] = product
	return product, nil
}

// Update updates a product in-memory.
func (r InMemoryProductsRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	_, span := r.Tracer.Start(ctx, "repository.Update")
	defer span.End()

	if _, ok := r.storage[product.ID]; !ok {
		return domain.Product{}, ErrProductNotFound
	}

	r.storage[product.ID] = product
	return product, nil
}

// Delete deletes a Product from memory.
func (r InMemoryProductsRepository) Delete(ctx context.Context, id int) error {
	_, span := r.Tracer.Start(ctx, "repository.Delete")
	defer span.End()

	_, found := r.storage[id]
	if !found {
		return ErrProductNotFound
	}

	delete(r.storage, id)

	return nil
}

// List all products from memory.
func (r InMemoryProductsRepository) List(ctx context.Context, filter domain.ListProductsFilter) ([]domain.Product, error) {
	_, span := r.Tracer.Start(ctx, "repository.List")
	defer span.End()

	filterIndex := map[int]bool{}
	for _, id := range filter.IDs {
		filterIndex[id] = true
	}

	var result []domain.Product
	for _, product := range r.storage {
		if len(filter.IDs) == 0 {
			result = append(result, product)
			continue
		}

		if _, ok := filterIndex[product.ID]; ok {
			result = append(result, product)
		}
	}

	return result, nil
}
