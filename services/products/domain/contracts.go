package domain

import "context"

// Application defines boundary interfaces of the application
// It should be called by the ports
type Application interface {
	// ListProducts fetches the list of Products
	ListProducts(context.Context, ListProductsFilter) ([]Product, error)

	// RegisterProduct registers a new Product
	RegisterProduct(context.Context, Product) (Product, error)

	// UpdateProduct updates a Product
	UpdateProduct(context.Context, Product) (Product, error)

	// DeleteProduct deletes a Product
	DeleteProduct(context.Context, int) error
}

// CLI defines the boundary interfaces of the application
// It should be called by the ports
type CLI interface {
	// RegisterProduct registers a new Product
	RegisterProduct(context.Context, Product) (Product, error)

	// DeleteProduct deletes a Product
	DeleteProduct(context.Context, int) error
}

type ProductsRepository interface {
	// Create creates a new Product in a data storage.
	Create(context.Context, Product) (Product, error)

	// Update updates a Product in a data storage.
	Update(context.Context, Product) (Product, error)

	// Delete deletes a Product from a data storage.
	Delete(context.Context, int) error

	// List all Products from a data storage.
	List(context.Context, ListProductsFilter) ([]Product, error)
}

// ListProductsFilter represents a filter passed to List
type ListProductsFilter struct {
	IDs []int
}
