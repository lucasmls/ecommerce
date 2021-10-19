package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/lucasmls/ecommerce/services/bff/ports/graphql/generated"
	"github.com/lucasmls/ecommerce/services/bff/ports/graphql/model"
)

func (r *mutationResolver) RegisterProduct(ctx context.Context, input model.RegisterProductInput) (*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, input model.UpdateProductInput) (*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) RemoveProduct(ctx context.Context, input model.RemoveProductInput) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
