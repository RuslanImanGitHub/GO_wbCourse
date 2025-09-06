package main

import (
	sarKaf "L0/internal"
	fullS "L0/models/compositeView/service"
	deliveryRepo "L0/models/deliveryAddress/postgresql"
	deliveryService "L0/models/deliveryAddress/service"
	ioRepo "L0/models/itemDetails/postgresql"
	ioService "L0/models/itemDetails/service"
	"L0/models/order"
	orderRepo "L0/models/order/postgresql"
	orderService "L0/models/order/service"
	paymentRepo "L0/models/payment/postgresql"
	paymentService "L0/models/payment/service"
	"L0/pkg/dbClient/postgresql"
	"L0/router"
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/spf13/viper"
)

func main() {
	//Config
	if err := initConfig(); err != nil {
		log.Fatalf("Error initializing config: %s", err.Error())
	}

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
	cache := expirable.NewLRU[string, order.Order](viper.GetInt("cache.objsInCache"), nil, 0)
	cacheOrd, err := fullS.FindN(ctx, pool, viper.GetInt("cache.cacheBuffer"))
	if err != nil {
		log.Println("Failed to fill cache")

	}

	for _, ord := range cacheOrd {
		cache.Add(ord.Order_uid, ord)
	}

	//kafka controller
	brokers := []string{fmt.Sprintf("%s:%s", viper.GetString("kafka.host"), viper.GetString("kafka.kafkaPort"))}
	consumer, err := sarKaf.NewConsumer(brokers)
	if err != nil {
		log.Fatal("Failed to create KafkaConsumer")
	}
	defer consumer.Close()

	producer, err := sarKaf.NewProducer(brokers)
	if err != nil {
		log.Fatal("Failed to create KafkaProducer")
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

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
