package repositories

import (
	"context"
	"errors"
	"math/rand"
	"strconv"

	"github.com/lucasmls/ecommerce/services/products/domain"
)

type InMemoryProductsRepository struct {
	storage map[string]domain.Product
}

// NewInMemoryProductsRepository creates a new InMemoryProductsRepository.
func NewInMemoryProductsRepository(size int) (InMemoryProductsRepository, error) {
	if size == 0 {
		return InMemoryProductsRepository{}, errors.New("invalid-storage-size")
	}

	storage := make(map[string]domain.Product, size)

	storage["10"] = domain.Product{
		ID:          "10",
		Name:        "Macbook Air M1",
		Description: "Fast!",
		Price:       6900,
	}

	return InMemoryProductsRepository{
		storage: storage,
	}, nil
}

// MustNewInMemoryProductsRepository creates a new InMemoryProductsRepository.
// It panics if any error is found.
func MustNewInMemoryProductsRepository(size int) InMemoryProductsRepository {
	repo, err := NewInMemoryProductsRepository(size)
	if err != nil {
		panic(err)
	}

	return repo
}

// Create creates a new Product in-memory.
func (r InMemoryProductsRepository) Create(ctx context.Context, product domain.Product) (domain.Product, error) {
	id := strconv.Itoa(rand.Intn(1000))
	product.ID = id

	r.storage[id] = product
	return product, nil
}

// Update updates a product in-memory.
func (r InMemoryProductsRepository) Update(ctx context.Context, product domain.Product) (domain.Product, error) {
	r.storage[product.ID] = product
	return product, nil
}

// Delete deletes a Product from memory.
func (r InMemoryProductsRepository) Delete(ctx context.Context, id string) error {
	_, found := r.storage[id]
	if !found {
		return errors.New("not-found")
	}

	delete(r.storage, id)

	return nil
}

// List list all products from memory.
func (r InMemoryProductsRepository) List(ctx context.Context, filter domain.ListProductsFilter) ([]domain.Product, error) {
	result := []domain.Product{}

	for _, v := range r.storage {
		result = append(result, v)
	}

	return result, nil
}
