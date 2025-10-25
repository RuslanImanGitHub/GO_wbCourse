package main

import (
	sarKaf "L0/backend/kafka"
	fullS "L0/backend/models/compositeView/service"
	deliveryRepo "L0/backend/models/deliveryAddress/postgresql"
	deliveryService "L0/backend/models/deliveryAddress/service"
	ioRepo "L0/backend/models/itemDetails/postgresql"
	ioService "L0/backend/models/itemDetails/service"
	"L0/backend/models/order"
	orderRepo "L0/backend/models/order/postgresql"
	orderService "L0/backend/models/order/service"
	paymentRepo "L0/backend/models/payment/postgresql"
	paymentService "L0/backend/models/payment/service"
	"L0/backend/pkg/dbClient/postgresql"
	"L0/backend/router"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/golang-lru/v2/expirable"
)

func main() {

	ctx := context.Background()

	//DB
	pool, err := postgresql.NewClient(ctx)
	if err != nil {
		log.Fatal("Couldn't establish DB Connection")
	}

	orderRepo := orderRepo.NewRepository(pool)
	orderS := orderService.NewService(orderRepo)

	deliveryRepo := deliveryRepo.NewRepository(pool)
	deliveryS := deliveryService.NewService(deliveryRepo)

	paymentRepo := paymentRepo.NewRepository(pool)
	paymentS := paymentService.NewService(paymentRepo)

	ioRepo := ioRepo.NewRepository(pool)
	ioS := ioService.NewService(ioRepo)

	fullS := fullS.NewService(pool, *deliveryS, *orderS, *paymentS, *ioS)
	log.Println("Connected to DB")

	//Cache
	cache := expirable.NewLRU[string, order.Order](1000, nil, 0)
	cacheOrd, err := fullS.FindN(ctx, pool, 100)
	if err != nil {
		log.Println("Failed to fill cache")

	}

	for _, ord := range cacheOrd {
		cache.Add(ord.Order_uid, ord)
	}

	//kafka controller
	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	consumer, err := sarKaf.NewConsumer(brokers)
	if err != nil {
		log.Fatal("Failed to create KafkaConsumer", err)
	}
	defer consumer.Close()

	producer, err := sarKaf.NewProducer(brokers)
	if err != nil {
		log.Fatal("Failed to create KafkaProducer", err)
	}
	defer producer.Close()

	kafkaCont := &sarKaf.SarKafCont{
		Producer: producer,
		Consumer: consumer,
	}

	//server
	router := router.StartRouter(ctx, fullS, cache, kafkaCont)
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}

	fmt.Println("App started successfully!")
}
