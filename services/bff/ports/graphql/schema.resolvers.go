package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/lucasmls/ecommerce/services/bff/ports/graphql/generated"
	"github.com/lucasmls/ecommerce/services/bff/ports/graphql/model"
)

func (m *mutationResolver) RegisterProduct(ctx context.Context, input model.RegisterProductInput) (*model.Product, error) {
	p := &model.Product{
		ID:          fmt.Sprint(rand.Int()),
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}

	m.data = append(m.data, p)

	return p, nil
}

func (q *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	return q.data, nil
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
