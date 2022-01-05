package repositories

import (
	"context"
	"errors"
	"github.com/lucasmls/ecommerce/services/products/domain"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"math/rand"
)

var (
	ErrInvalidStorageSize = errors.New("invalid-storage-size")
	ErrProductNotFound    = errors.New("product-not-found")
)

// ProductsRepositoryInput holds all the dependencies needed to
// instantiate ProductsRepository
type ProductsRepositoryInput struct {
	Logger *zap.Logger
	Tracer trace.Tracer

	Size int
}

type InMemoryProductsRepository struct {
	in      ProductsRepositoryInput
	storage map[int]domain.Product
}

// NewInMemoryProductsRepository creates a new InMemoryProductsRepository.
func NewInMemoryProductsRepository(in ProductsRepositoryInput) (InMemoryProductsRepository, error) {
	if in.Size == 0 {
		return InMemoryProductsRepository{}, ErrInvalidStorageSize
	}

	storage := make(map[int]domain.Product, in.Size)

	storage[10] = domain.Product{
		ID:          10,
		Name:        "Macbook Air M1",
		Description: "Fast!",
		Price:       6900,
	}

	return InMemoryProductsRepository{
		in:      in,
		storage: storage,
	}, nil
}

// MustNewInMemoryProductsRepository creates a new InMemoryProductsRepository.
// It panics if any error is found.
func MustNewInMemoryProductsRepository(in ProductsRepositoryInput) InMemoryProductsRepository {
	repo, err := NewInMemoryProductsRepository(in)
	if err != nil {
		panic(err)
	}

	return repo
}

// Create creates a new Product in-memory.
func (r InMemoryProductsRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	_, span := r.in.Tracer.Start(ctx, "repository.Create")
	defer span.End()

	id := rand.Intn(1000)
	product.ID = id

	r.storage[id] = product
	return product, nil
}

// Update updates a product in-memory.
func (r InMemoryProductsRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	_, span := r.in.Tracer.Start(ctx, "repository.Update")
	defer span.End()

	_, ok := r.storage[product.ID]
	if !ok {
		return domain.Product{}, ErrProductNotFound
	}

	r.storage[product.ID] = product
	return product, nil
}

// Delete deletes a Product from memory.
func (r InMemoryProductsRepository) Delete(ctx context.Context, id int) error {
	_, span := r.in.Tracer.Start(ctx, "repository.Delete")
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
	_, span := r.in.Tracer.Start(ctx, "repository.List")
	defer span.End()

	var result []domain.Product
	for _, v := range r.storage {
		result = append(result, v)
	}

	return result, nil
}
