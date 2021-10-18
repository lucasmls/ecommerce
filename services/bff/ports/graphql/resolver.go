package graph

import "github.com/lucasmls/ecommerce/services/bff/ports/graphql/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	data []*model.Product
}
