package grpc_port

import (
	"context"
	"errors"

	pb "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
	"go.uber.org/zap"
)

// ProductsResolverInput ...
type ProductsResolverInput struct {
	Logger *zap.Logger
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

func (r *ProductsResolver) List(ctx context.Context, req *pb.ListRequest) (*pb.ListResponse, error) {
	r.in.Logger.Info("received a request to list Products", zap.Strings("ids", req.Ids))
	r.in.Logger.Error("failed to list products", zap.Error(errors.New("testing zap.Error")))
	r.in.Logger.Debug("debug log!", zap.Int("int_value", 10), zap.Float32("float_value", 12))
	r.in.Logger.Warn("warn log!", zap.Bool("bool_value", true))

	response := pb.ListResponse{
		Data: []*pb.Product{
			{Id: "1", Name: "Macbook Air M1", Description: "Fast!", Price: 6900},
			{Id: "2", Name: "Macbook Pro M1", Description: "Super fast!", Price: 9000},
		},
	}
	return &response, nil
}
