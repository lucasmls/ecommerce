package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/lucasmls/ecommerce/services/bff/ports/graphql/generated"
	"github.com/lucasmls/ecommerce/services/bff/ports/graphql/model"
	grpc_protobuf "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"go.uber.org/zap"
)

func (m *mutationResolver) RegisterProduct(ctx context.Context, input model.RegisterProductInput) (*model.Product, error) {
	m.Logger.Info("registering a new product", zap.Any("input", input))

	registeredProduct, err := m.ProductsService.Register(ctx, &grpc_protobuf.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       float32(input.Price),
	})
	if err != nil {
		return nil, err
	}

	response := &model.Product{
		ID:          registeredProduct.Data.Id,
		Name:        registeredProduct.Data.Name,
		Description: registeredProduct.Data.Description,
		Price:       float64(registeredProduct.Data.Price),
	}

	return response, nil
}

func (m *mutationResolver) UpdateProduct(ctx context.Context, input model.UpdateProductInput) (*model.Product, error) {
	m.Logger.Info("updating a product", zap.Any("input", input))

	updatedProduct, err := m.ProductsService.Update(ctx, &grpc_protobuf.Product{
		Id:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Price:       float32(input.Price),
	})
	if err != nil {
		return nil, err
	}

	response := &model.Product{
		ID:          updatedProduct.Data.Id,
		Name:        updatedProduct.Data.Name,
		Description: updatedProduct.Data.Description,
		Price:       float64(updatedProduct.Data.Price),
	}

	return response, nil
}

func (q *queryResolver) Products(ctx context.Context) ([]*model.Product, error) {
	q.Logger.Info("querying products")

	products, err := q.ProductsService.List(ctx, &grpc_protobuf.ListRequest{})
	if err != nil {
		return nil, err
	}

	response := []*model.Product{}
	for _, product := range products.Data {
		response = append(response, &model.Product{
			ID:          product.Id,
			Name:        product.Name,
			Description: product.Description,
			Price:       float64(product.Price),
		})
	}

	return response, nil
}

func (m *mutationResolver) RemoveProduct(ctx context.Context, input model.RemoveProductInput) (string, error) {
	m.Logger.Info("removing a product", zap.String("id", input.ID))

	deleteReponse, err := m.ProductsService.Delete(ctx, &grpc_protobuf.DeleteRequest{Id: input.ID})
	if err != nil {
		return "", err
	}

	return deleteReponse.Data, nil
}

// Mutation returns generated1.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated1.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
