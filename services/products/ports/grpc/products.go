package grpc_port

import (
	"context"
	"errors"

	"github.com/lucasmls/ecommerce/services/products/domain"
	pb "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InternalServerError = status.Error(codes.Internal, "Internal server error")
)

var (
	ErrMissingLogger = errors.New("missing required dependency: Logger")
	ErrMisingTracer  = errors.New("missing required dependency: Tracer")
)

// ProductsResolver ...
type ProductsResolver struct {
	Logger *zap.Logger
	Tracer trace.Tracer
	App    domain.Application

	pb.UnimplementedProductsServiceServer
}

func NewProductsResolver(
	logger *zap.Logger,
	tracer trace.Tracer,
	app domain.Application,
) (*ProductsResolver, error) {
	if logger == nil {
		return nil, ErrMissingLogger
	}

	if tracer == nil {
		return nil, ErrMisingTracer
	}

	return &ProductsResolver{
		Logger: logger,
		Tracer: tracer,
		App:    app,
	}, nil
}

// MustNewProductsResolver creates a new ProductsResolver instance.
//It panics if any error is found
func MustNewProductsResolver(
	logger *zap.Logger,
	tracer trace.Tracer,
	app domain.Application,
) *ProductsResolver {
	productsResolver, err := NewProductsResolver(logger, tracer, app)
	if err != nil {
		panic(err)
	}

	return productsResolver
}

func (r *ProductsResolver) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	ctx, span := r.Tracer.Start(ctx, "resolver.List")
	defer span.End()

	var filter domain.ListProductsFilter
	for _, id := range req.Ids {
		filter.IDs = append(filter.IDs, int(id))
	}

	r.Logger.Info("received a request to list Products", zap.Any("filter", filter))

	products, err := r.App.ListProducts(ctx, filter)
	if err != nil {
		r.Logger.Sugar().Errorw(
			"failed to list products",
			zap.Error(err),
			zap.Any("filter", filter),
		)

		return nil, InternalServerError
	}

	response := pb.ListResponse{
		Data: []*pb.Product{},
	}

	for _, product := range products {
		response.Data = append(response.Data, &pb.Product{
			Id:          int32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       int32(product.Price),
		})
	}

	return &response, nil
}

func (r *ProductsResolver) Register(ctx context.Context, req *pb.Product) (*pb.RegisterResponse, error) {
	ctx, span := r.Tracer.Start(ctx, "resolver.Register")
	defer span.End()

	r.Logger.Info("registering a new product", zap.Any("req", req))

	product, err := r.App.RegisterProduct(ctx, domain.Product{
		ID:          int(req.Id),
		Name:        req.Name,
		Description: req.Description,
		Price:       int(req.Price),
	})
	if err != nil {
		r.Logger.Sugar().Errorw(
			"failed to register the provided product",
			zap.Error(err),
			zap.Any("product", product),
		)

		return nil, InternalServerError
	}

	response := &pb.RegisterResponse{
		Data: &pb.Product{
			Id:          int32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       int32(product.Price),
		},
	}

	return response, nil
}

func (r *ProductsResolver) Update(ctx context.Context, req *pb.Product) (*pb.UpdateResponse, error) {
	ctx, span := r.Tracer.Start(ctx, "resolver.Update")
	defer span.End()

	r.Logger.Info("updating a product", zap.Any("req", req))

	product, err := r.App.UpdateProduct(ctx, domain.Product{
		ID:          int(req.Id),
		Name:        req.Name,
		Description: req.Description,
		Price:       int(req.Price),
	})
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			r.Logger.Debug(
				"the provided product to be updated was not found",
				zap.Error(err),
				zap.Any("product", product),
			)

			return nil, status.Error(codes.NotFound, err.Error())
		}

		r.Logger.Sugar().Errorw(
			"failed to update the provided product",
			zap.Error(err),
			zap.Any("product", product),
		)

		return nil, InternalServerError
	}

	response := &pb.UpdateResponse{
		Data: &pb.Product{
			Id:          int32(product.ID),
			Name:        product.Name,
			Description: product.Description,
			Price:       int32(product.Price),
		},
	}

	return response, nil
}

func (r *ProductsResolver) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	ctx, span := r.Tracer.Start(ctx, "resolver.Delete")
	defer span.End()

	productId := int(req.Id)

	r.Logger.Info("deleting a product", zap.Int("id", productId))

	err := r.App.DeleteProduct(ctx, productId)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

		r.Logger.Sugar().Errorw(
			"failed to delete the product",
			zap.Error(err),
			zap.Int("productId", productId),
		)

		return nil, InternalServerError
	}

	response := &pb.DeleteResponse{
		Data: "Product deleted successfully",
	}

	return response, nil
}
