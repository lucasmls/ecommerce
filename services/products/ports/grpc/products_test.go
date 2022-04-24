package grpc_port

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/lucasmls/ecommerce/shared/grpc"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	gGRPC "google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"github.com/lucasmls/ecommerce/services/products/domain"
	"github.com/lucasmls/ecommerce/services/products/mocks"
	protog "github.com/lucasmls/ecommerce/services/products/ports/grpc/proto"
)

type ProductsResolverSuite struct {
	suite.Suite

	app    *mocks.Application
	logger *zap.Logger
	tracer trace.Tracer

	grpcClient     protog.ProductsServiceClient
	grpcConnection *gGRPC.ClientConn
}

func (s *ProductsResolverSuite) SetupSuite() {
	s.logger = zap.NewNop()
	s.tracer = trace.NewNoopTracerProvider().Tracer("")
	s.app = &mocks.Application{}

	grpcClient, grpcConnection := makeGrpcClientAndServer(
		context.Background(),
		s.logger,
		s.tracer,
		s.app,
	)

	s.grpcClient = grpcClient
	s.grpcConnection = grpcConnection
}

func (s *ProductsResolverSuite) TearDownSuite() {
	err := s.grpcConnection.Close()
	s.NoError(err)
}

func (s *ProductsResolverSuite) Test_NewProductsResolver() {
	s.Run("Should fail to instantiate the ProductsResolver in case a Logger isn't provided", func() {
		_, err := NewProductsResolver(nil, nil, nil)

		s.Equal(ErrMissingLogger, err)
	})

	s.Run("Should fail to instantiate the ProductsResolver in case a Tracer isn't provided", func() {
		_, err := NewProductsResolver(s.logger, nil, nil)

		s.Equal(ErrMisingTracer, err)
	})
	s.Run("Should successfully instantiate the ProductsResolver", func() {
		_, err := NewProductsResolver(s.logger, s.tracer, nil)

		s.NoError(err)
	})
}

func (s *ProductsResolverSuite) Test_MustNewProductsResolver() {
	s.Run("Should fail to instantiate the ProductsResolver in case a Logger isn't provided", func() {
		defer func() {
			err := recover()
			s.NotNil(err)
			s.Equal(ErrMissingLogger, err)
		}()

		MustNewProductsResolver(nil, nil, nil)
	})

	s.Run("Should fail to instantiate the ProductsResolver in case a Tracer isn't provided", func() {
		defer func() {
			err := recover()
			s.NotNil(err)
			s.Equal(ErrMisingTracer, err)
		}()

		MustNewProductsResolver(s.logger, nil, nil)
	})
	s.Run("Should successfully instantiate the ProductsResolver", func() {
		defer func() {
			err := recover()
			s.Nil(err)
		}()

		MustNewProductsResolver(s.logger, s.tracer, nil)
	})
}

func (s *ProductsResolverSuite) Test_Delete() {
	s.Run("Should return not found error in case the specified product isn't stored", func() {
		ctx := context.Background()
		productId := 1

		expectedResult := status.Error(codes.NotFound, domain.ErrProductNotFound.Error())

		s.app.
			On("DeleteProduct",
				mock.AnythingOfType("*context.valueCtx"),
				productId,
			).
			Return(domain.ErrProductNotFound)

		_, err := s.grpcClient.Delete(ctx, &protog.DeleteRequest{
			Id: int32(productId),
		})

		s.Equal(expectedResult, err)
	})

	s.Run("Should return a generic error in case we receive a error that we're not aware of", func() {
		ctx := context.Background()
		productId := 2

		expectedResult := status.Error(codes.Internal, "Internal server error")

		s.app.
			On("DeleteProduct",
				mock.AnythingOfType("*context.valueCtx"),
				productId,
			).
			Return(errors.New("mock error"))

		_, err := s.grpcClient.Delete(ctx, &protog.DeleteRequest{
			Id: int32(productId),
		})

		s.Equal(expectedResult, err)
	})

	s.Run("Should successfully delete the specified product", func() {
		ctx := context.Background()
		productId := 3

		expectedResult := &protog.DeleteResponse{
			Data: "Product deleted successfully",
		}

		s.app.
			On("DeleteProduct",
				mock.AnythingOfType("*context.valueCtx"),
				productId,
			).
			Return(nil)

		got, err := s.grpcClient.Delete(ctx, &protog.DeleteRequest{
			Id: int32(productId),
		})

		s.NoError(err)
		s.True(proto.Equal(expectedResult, got))
	})
}

func (s *ProductsResolverSuite) Test_List() {
	s.Run("Should return a generic error in case we receive a error that we're not aware of", func() {
		ctx := context.Background()
		filter := domain.ListProductsFilter{}

		expectedResult := status.Error(codes.Internal, "Internal server error")

		s.app.
			On("ListProducts",
				mock.AnythingOfType("*context.valueCtx"),
				filter,
			).
			Return(nil, errors.New("mock error"))

		_, err := s.grpcClient.List(ctx, &protog.ListRequest{
			Ids: []int32{},
		})

		s.Equal(expectedResult, err)
	})

	s.Run("Should successfully list the Products", func() {
		ctx := context.Background()
		filter := domain.ListProductsFilter{
			IDs: []int{1, 2},
		}

		products := []domain.Product{
			{
				ID:          1,
				Name:        "Iphone 13",
				Description: "Cool",
				Price:       4500,
			},
			{
				ID:          2,
				Name:        "Macbook Pro M1 Max",
				Description: "Fast!",
				Price:       16500,
			},
		}

		expectedResult := protog.ListResponse{
			Data: []*protog.Product{},
		}

		for _, product := range products {
			expectedResult.Data = append(expectedResult.Data, &protog.Product{
				Id:          int32(product.ID),
				Name:        product.Name,
				Description: product.Description,
				Price:       int32(product.Price),
			})
		}

		s.app.
			On("ListProducts",
				mock.AnythingOfType("*context.valueCtx"),
				filter,
			).
			Return(products, nil)

		got, err := s.grpcClient.List(ctx, &protog.ListRequest{
			Ids: []int32{1, 2},
		})

		s.NoError(err)
		for i := range got.Data {
			got := got.Data[i]
			expected := expectedResult.Data[i]

			s.True(proto.Equal(expected, got))
		}
	})
}

func (s *ProductsResolverSuite) Test_Register() {
	s.Run("Should return a generic error in case we receive a error that we're not aware of", func() {
		ctx := context.Background()
		req := &protog.Product{
			Id:          1,
			Name:        "Macbook Air M1",
			Description: "Fast",
			Price:       6800,
		}

		product := domain.Product{
			ID:          int(req.Id),
			Name:        req.Name,
			Description: req.Description,
			Price:       int(req.Price),
		}

		expectedResult := status.Error(codes.Internal, "Internal server error")

		s.app.
			On("RegisterProduct",
				mock.AnythingOfType("*context.valueCtx"),
				product,
			).
			Return(domain.Product{}, errors.New("mock error"))

		_, err := s.grpcClient.Register(ctx, req)

		s.Equal(expectedResult, err)
	})

	s.Run("Should successfully register the provided Product", func() {
		ctx := context.Background()
		req := &protog.Product{
			Id:          2,
			Name:        "Macbook Air M1",
			Description: "Fast",
			Price:       6800,
		}

		expectedResult := &protog.RegisterResponse{
			Data: &protog.Product{
				Id:          int32(req.Id),
				Name:        req.Name,
				Description: req.Description,
				Price:       int32(req.Price),
			},
		}

		product := domain.Product{
			ID:          int(req.Id),
			Name:        req.Name,
			Description: req.Description,
			Price:       int(req.Price),
		}

		s.app.
			On("RegisterProduct",
				mock.AnythingOfType("*context.valueCtx"),
				product,
			).
			Return(product, nil)

		got, err := s.grpcClient.Register(ctx, req)

		s.NoError(err)
		s.True(proto.Equal(expectedResult, got))
	})
}

func (s *ProductsResolverSuite) Test_Update() {
	s.Run("Should return not found in case the provided Product isn't stored", func() {
		ctx := context.Background()
		req := &protog.Product{
			Id:          1,
			Name:        "Macbook Air M1",
			Description: "Fast",
			Price:       6800,
		}

		product := domain.Product{
			ID:          int(req.Id),
			Name:        req.Name,
			Description: req.Description,
			Price:       int(req.Price),
		}

		expectedResult := status.Error(codes.NotFound, domain.ErrProductNotFound.Error())

		s.app.
			On("UpdateProduct",
				mock.AnythingOfType("*context.valueCtx"),
				product,
			).
			Return(domain.Product{}, domain.ErrProductNotFound)

		_, err := s.grpcClient.Update(ctx, req)

		s.Equal(expectedResult, err)
	})

	s.Run("Should return a generic error in case we receive a error that we're not aware of", func() {
		ctx := context.Background()
		req := &protog.Product{
			Id:          2,
			Name:        "Macbook Air M1",
			Description: "Fast",
			Price:       6800,
		}

		product := domain.Product{
			ID:          int(req.Id),
			Name:        req.Name,
			Description: req.Description,
			Price:       int(req.Price),
		}

		expectedResult := status.Error(codes.Internal, "Internal server error")

		s.app.
			On("UpdateProduct",
				mock.AnythingOfType("*context.valueCtx"),
				product,
			).
			Return(domain.Product{}, errors.New("mock error"))

		_, err := s.grpcClient.Update(ctx, req)

		s.Equal(expectedResult, err)
	})

	s.Run("Should successfully update the provided Product", func() {
		ctx := context.Background()
		req := &protog.Product{
			Id:          3,
			Name:        "Macbook Air M1",
			Description: "Fast",
			Price:       6800,
		}

		expectedResult := &protog.UpdateResponse{
			Data: &protog.Product{
				Id:          int32(req.Id),
				Name:        req.Name,
				Description: req.Description,
				Price:       int32(req.Price),
			},
		}

		product := domain.Product{
			ID:          int(req.Id),
			Name:        req.Name,
			Description: req.Description,
			Price:       int(req.Price),
		}

		s.app.
			On("UpdateProduct",
				mock.AnythingOfType("*context.valueCtx"),
				product,
			).
			Return(product, nil)

		got, err := s.grpcClient.Update(ctx, req)

		s.NoError(err)
		s.True(proto.Equal(expectedResult, got))
	})
}

func TestProductsResolverSuite(t *testing.T) {
	suite.Run(t, new(ProductsResolverSuite))
}

func makeGrpcClientAndServer(
	ctx context.Context,
	loggerM *zap.Logger,
	tracerM trace.Tracer,
	appM domain.Application,
) (protog.ProductsServiceClient, *gGRPC.ClientConn) {
	listener := makeGrpcServer(ctx, loggerM, tracerM, appM)
	grpcClient, grpcConnection := makeGrpcClient(ctx, loggerM, listener)

	return grpcClient, grpcConnection
}

func makeGrpcClient(
	ctx context.Context,
	loggerM *zap.Logger,
	listener *bufconn.Listener,
) (protog.ProductsServiceClient, *gGRPC.ClientConn) {
	productsServiceGRPCClient := grpc.MustNewClient(grpc.ClientInput{
		Logger: loggerM,
		AdditionalDialOptions: []gGRPC.DialOption{
			gGRPC.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
				return listener.Dial()
			}),
		},
	})

	productsServiceConn := productsServiceGRPCClient.MustConnect(ctx)

	grpcClient := protog.NewProductsServiceClient(productsServiceConn)
	return grpcClient, productsServiceConn
}

func makeGrpcServer(
	ctx context.Context,
	loggerM *zap.Logger,
	tracerM trace.Tracer,
	appM domain.Application,
) *bufconn.Listener {
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.MustNewServer(grpc.ServerInput{
		Port:   1234,
		Logger: loggerM,
		Registrator: func(server gGRPC.ServiceRegistrar) {
			protog.RegisterProductsServiceServer(server, &ProductsResolver{
				Logger: loggerM,
				Tracer: tracerM,
				App:    appM,
			})
		},
		Listener: listener,
	})

	go func() {
		if err := grpcServer.Run(ctx); err != nil {
			panic(err)
		}
	}()

	return listener
}
