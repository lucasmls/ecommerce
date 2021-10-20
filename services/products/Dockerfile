# --- Builder
FROM golang:1.17.2 as builder
ENV CGO_ENABLED=0
WORKDIR /app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download
COPY . .

RUN go build -o /grpc_server -mod=readonly ./cmd/grpc/main.go

# ---
FROM alpine as grpc_server

RUN apk add --no-cache ca-certificates
COPY --from=builder /grpc_server /grpc_server
ENTRYPOINT [ "/grpc_server" ]
