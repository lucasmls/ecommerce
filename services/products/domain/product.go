package domain

import "errors"

// Product represents a product in the system.
type Product struct {
	ID          int
	Name        string
	Description string
	Price       int
}

var (
	ErrProductNotFound = errors.New("product-not-found")
)
