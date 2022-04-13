package grpc_port

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/lucasmls/ecommerce/shared/grpc"
	"github.com/stretchr/testify/assert"
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

func Test_Delete(t *testing.T) {
	loggerM := zap.NewNop()
	tracerM := trace.NewNoopTracerProvider().Tracer("")

	t.Run("Failure tests", func(t *testing.T) {
		tt := []struct {
			name           string
			productId      int
			expectedResult error
			appM           func(ctx context.Context, ctrl *gomock.Controller, productId int) *mocks.MockApplication
		}{
			{
				name:           "Should return not found error in case the specified product isn't stored",
				productId:      1,
				expectedResult: status.Error(codes.NotFound, domain.ErrProductNotFound.Error()),
				appM: func(ctx context.Context, ctrl *gomock.Controller, productId int) *mocks.MockApplication {
					appM := mocks.NewMockApplication(ctrl)

					// We have a problem here that if we provide a productId that does not match
					// with the expected, the test simply times out.
					// This happens becase some way our unit test does not get to know that
					// id didn't match...
					appM.EXPECT().
						DeleteProduct(gomock.Any(), productId).
						Return(domain.ErrProductNotFound)

					return appM
				},
			},
			{
				name:           "Should return a generic error in case we receive a error that we're not aware of",
				productId:      1,
				expectedResult: status.Error(codes.Internal, "Internal server error"),
				appM: func(ctx context.Context, ctrl *gomock.Controller, productId int) *mocks.MockApplication {
					appM := mocks.NewMockApplication(ctrl)

					appM.EXPECT().
						DeleteProduct(gomock.Any(), productId).
						Return(errors.New("mock error"))

					return appM
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				ctrl := gomock.NewController(t)

				grpcClient, grpcConnection := makeGrpcClientAndServer(
					ctx,
					loggerM,
					tracerM,
					tc.appM(ctx, ctrl, tc.productId),
				)

				defer grpcConnection.Close()

				_, err := grpcClient.Delete(ctx, &protog.DeleteRequest{
					Id: int32(tc.productId),
				})

				assert.Equal(t, tc.expectedResult, err)
			})
		}

	})

	t.Run("Successful tests", func(t *testing.T) {
		tt := []struct {
			name           string
			productId      int
			expectedResult *protog.DeleteResponse
			appM           func(ctx context.Context, ctrl *gomock.Controller, productId int) *mocks.MockApplication
		}{
			{
				name:      "Should successfully delete the specified product",
				productId: 1,
				expectedResult: &protog.DeleteResponse{
					Data: "Product deleted successfully",
				},
				appM: func(ctx context.Context, ctrl *gomock.Controller, productId int) *mocks.MockApplication {
					appM := mocks.NewMockApplication(ctrl)

					appM.EXPECT().
						DeleteProduct(gomock.Any(), productId).
						Return(nil)

					return appM
				},
			},
		}

		for _, tc := range tt {
			t.Run(tc.name, func(t *testing.T) {
				ctx := context.Background()
				ctrl := gomock.NewController(t)

				grpcClient, grpcConnection := makeGrpcClientAndServer(
					ctx,
					loggerM,
					tracerM,
					tc.appM(ctx, ctrl, tc.productId),
				)

				defer grpcConnection.Close()

				response, err := grpcClient.Delete(ctx, &protog.DeleteRequest{
					Id: int32(tc.productId),
				})

				assert.NoError(t, err)
				assert.True(
					t,
					proto.Equal(tc.expectedResult, response),
				)
			})
		}
	})
}

func makeGrpcClientAndServer(
	ctx context.Context,
	loggerM *zap.Logger,
	tracerM trace.Tracer,
	appM *mocks.MockApplication,
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
	appM *mocks.MockApplication,
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
