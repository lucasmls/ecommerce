package main

import (
	"log"

	protoMessages "github.com/lucasmls/ecommerce/services/products/ports/rmq/proto"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

const (
	productsExchange        string = "products"
	registerProductsRouting string = "register"
	amqpConnectionString    string = "amqp:guest:guest@localhost:5672/"
)

func main() {
	ampqConnection, err := amqp.Dial(amqpConnectionString)
	if err != nil {
		log.Fatal(err)
	}

	defer ampqConnection.Close()

	amqpChannel, err := ampqConnection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	err = amqpChannel.ExchangeDeclare(
		productsExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}

	newProductMessage := &protoMessages.Product{
		Id:          1,
		Name:        "Macbook Air M1",
		Description: "Fast!",
		Price:       6800,
	}

	newProductPb, err := proto.Marshal(newProductMessage)
	if err != nil {
		log.Fatal(err)
	}

	err = amqpChannel.Publish(
		productsExchange,
		registerProductsRouting,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         newProductPb,
		},
	)
	if err != nil {
		log.Fatal(err)
	}
}
