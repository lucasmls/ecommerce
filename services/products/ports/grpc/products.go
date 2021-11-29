package grpc_port

import (
	"context"
	"errors"

	"github.com/lucasmls/ecommerce/services/products/domain"
	pb "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// ProductsResolverInput ...
type ProductsResolverInput struct {
	Logger *zap.Logger
	App    domain.Application
	Tracer trace.Tracer
}

// ProductsResolver ...
type ProductsResolver struct {
	in ProductsResolverInput

	pb.UnimplementedProductsServiceServer
}

func NewProductsResolver(in ProductsResolverInput) (*ProductsResolver, error) {
	if in.Logger == nil {
		return nil, errors.New("missing required dependency: Logger")
	}

	return &ProductsResolver{
		in: in,
	}, nil
}

// NewProductsResolver creates a new ProductsResolver instance
// It panics if any error is found
func MustNewProductsResolver(in ProductsResolverInput) *ProductsResolver {
	app, err := NewProductsResolver(in)
	if err != nil {
		panic(err)
	}

	return app
}

func (r *ProductsResolver) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	ctx, span := r.in.Tracer.Start(ctx, "resolver.List")
	defer span.End()

	r.in.Logger.Info("received a request to list Products", zap.Strings("ids", req.Ids))

	products, err := r.in.App.ListProducts(ctx, domain.ListProductsFilter{
		IDs: req.Ids,
	})
	if err != nil {
		return nil, err
	}

	response := pb.ListResponse{
		Data: []*pb.Product{},
	}

	for _, product := range products {
		response.Data = append(response.Data, &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		})
	}

	return &response, nil
}

func (r *ProductsResolver) Register(ctx context.Context, req *pb.Product) (*pb.RegisterResponse, error) {
	ctx, span := r.in.Tracer.Start(ctx, "resolver.Register")
	defer span.End()

	r.in.Logger.Info("registering a new product", zap.Any("product", req))

	product, err := r.in.App.RegisterProduct(ctx, domain.Product{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return nil, err
	}

	response := &pb.RegisterResponse{
		Data: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}

	return response, nil
}

func (r *ProductsResolver) Update(ctx context.Context, req *pb.Product) (*pb.UpdateResponse, error) {
	ctx, span := r.in.Tracer.Start(ctx, "resolver.Update")
	defer span.End()

	r.in.Logger.Info("updating a product", zap.Any("product", req))

	product, err := r.in.App.UpdateProduct(ctx, domain.Product{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
	})
	if err != nil {
		return nil, err
	}

	response := &pb.UpdateResponse{
		Data: &pb.Product{
			Id:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
		},
	}

	return response, nil
}

func (r *ProductsResolver) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	ctx, span := r.in.Tracer.Start(ctx, "resolver.Delete")
	defer span.End()

	r.in.Logger.Info("deleting a product", zap.String("id", req.Id))

	err := r.in.App.DeleteProduct(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	response := &pb.DeleteResponse{
		Data: "Product deleted successfully",
	}

	return response, nil
}
