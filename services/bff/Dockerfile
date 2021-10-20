# --- Builder
FROM golang:1.17.2 as builder
ENV CGO_ENABLED=0
WORKDIR /app

COPY ./go.mod ./go.mod
COPY ./go.sum ./go.sum
RUN go mod download
COPY . .

RUN go build -o /graphql_server -mod=readonly ./cmd/graphql/main.go

# ---
FROM alpine as graphql_server

RUN apk add --no-cache ca-certificates
COPY --from=builder /graphql_server /graphql_server
ENTRYPOINT [ "/graphql_server" ]
