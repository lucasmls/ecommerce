package domain

import "context"

type ProductsRepository interface {
	// Create creates a new Product in a data storage.
	Create(context.Context, Product) (Product, error)

	// Update updates a Product in a data storage.
	Update(context.Context, Product) (Product, error)

	// Delete deletes a Product from a data storage.
	Delete(context.Context, string) error

	// List list all Products from a data storage.
	List(context.Context, ListProductsFilter) ([]Product, error)
}

// ListProductsFilter represents a filter passed to List
type ListProductsFilter struct {
	IDs []string
}
