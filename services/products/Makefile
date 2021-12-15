
include .env
export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

deps:
	@ echo
	@ echo "Downloading dependencies..."
	@ echo
	@ go get -v ./...

update-deps:
	@ echo
	@ echo "Updating dependencies..."
	@ echo
	@ go get -u ./...

grpc-server:
	@ echo
	@ echo "Starting products gRPC server..."
	@ echo
	@ go run ./cmd/grpc/main.go

test:
	@ echo
	@ echo "Starting running tests..."
	@ echo
	@ go test -cover ./...

gen-grpc-proto:
	@ echo "Generating ./ports/grpc/*.proto into ./ports/grpc/gen/*.pb.go ..."
	@ protoc \
		--go_out=. \
    --go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
    ports/grpc/proto/*.proto

gen-rmq-proto:
	@ echo "Generating ./ports/rmq/*.proto into ./ports/rmq/gen/*.pb.go ..."
	@ protoc \
		--go_out=. \
    --go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
    ports/rmq/proto/*.proto

gen-proto: gen-grpc-proto gen-rmq-proto

mock:
	@ echo "Starting building mocks..."
	@ echo
	@ mkdir -p mocks
	@ rm mocks/*.go || true && \
		mockgen -source=./domain/contracts.go -destination=./mocks/domain.go -package=mocks

%:
	@:
