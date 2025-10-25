package main

import (
	sarKaf "L0/backend/kafka"
	fullS "L0/backend/models/compositeView/service"
	"L0/backend/models/compositeView/storage"
	"L0/backend/models/deliveryAddress"
	deliveryRepo "L0/backend/models/deliveryAddress/postgresql"
	deliveryService "L0/backend/models/deliveryAddress/service"
	itemdetails "L0/backend/models/itemDetails"
	ioRepo "L0/backend/models/itemDetails/postgresql"
	ioService "L0/backend/models/itemDetails/service"
	"L0/backend/models/order"
	orderRepo "L0/backend/models/order/postgresql"
	orderService "L0/backend/models/order/service"
	"L0/backend/models/payment"
	paymentRepo "L0/backend/models/payment/postgresql"
	paymentService "L0/backend/models/payment/service"
	"L0/backend/pkg/dbClient/postgresql"
	"L0/frontend/validation"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/jackc/pgx/v4/pgxpool"
)

var (
	orderFull  storage.Repository
	ctx        context.Context
	globalPool *pgxpool.Pool
)

func main() {
	log.Println("KafkaService started")

	//DB
	ctx = context.TODO()

	pool, err := postgresql.NewClient(ctx)
	if err != nil {
		log.Fatal("Couldn't establish DB Connection")
	}
	globalPool = pool

	orderRepo := orderRepo.NewRepository(pool)
	orderS := orderService.NewService(orderRepo)

	deliveryRepo := deliveryRepo.NewRepository(pool)
	deliveryS := deliveryService.NewService(deliveryRepo)

	paymentRepo := paymentRepo.NewRepository(pool)
	paymentS := paymentService.NewService(paymentRepo)

	ioRepo := ioRepo.NewRepository(pool)
	ioS := ioService.NewService(ioRepo)

	orderFull = fullS.NewService(pool, *deliveryS, *orderS, *paymentS, *ioS)

	//Kafka
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	brokers := []string{os.Getenv("KAFKA_BROKERS")}
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

	doneChan := make(chan struct{})

	go sarKaf.ReadMsgLoop(kafkaCont.Producer, kafkaCont.Consumer, doneChan, getOrderById, "getOrderById", "resGetOrderById")
	go sarKaf.ReadMsgLoop(kafkaCont.Producer, kafkaCont.Consumer, doneChan, postFakeOrder, "postOrder", "resPostOrder")
	go func(chan os.Signal, chan struct{}) {
		<-sigChan
		log.Println("Service stopped")
		doneChan <- struct{}{}
	}(sigChan, doneChan)

	<-doneChan
}

func postFakeOrder(id string) (interface{}, error) {

	source := rand.NewSource(time.Now().UnixNano())
	rnd := rand.New(source)

	var (
		order    order.Order
		delivery deliveryAddress.DeliveryAddress
		payment  payment.Payment
		items    []itemdetails.ItemDetails
	)

	err := gofakeit.Struct(&order)
	if err != nil {
		log.Printf("Ошибка генерации: %v\n", err)
		return -1, err
	}

	err = gofakeit.Struct(&delivery)
	if err != nil {
		log.Printf("Ошибка генерации: %v\n", err)
		return -1, err
	}

	err = gofakeit.Struct(&payment)
	if err != nil {
		log.Printf("Ошибка генерации: %v\n", err)
		return -1, err
	}

	itemsNumber := rnd.Intn(10) + 1
	for i := 0; i < itemsNumber; i++ {

		var item itemdetails.ItemDetails
		err = gofakeit.Struct(&item)
		if err != nil {
			log.Printf("Ошибка генерации: %v\n", err)
			return -1, err
		}
		item.Track_number = order.Track_number

		if err = validation.Val.Struct(&item); err != nil {
			log.Println("ItemOrder failed validation check: ", err)
			return -1, err
		}
		items = append(items, item)
	}

	order.Delivery = delivery
	order.Delivery_id = delivery.Delivery_id
	payment.Transaction = order.Order_uid
	order.Payment = payment
	order.Items = items

	if err = validation.Val.Struct(&order); err != nil {
		log.Printf("Order with uid (%s) failed validation check: %s", order.Order_uid, err)
		return -1, err
	}

	err = orderFull.Create(ctx, globalPool, &order, &order.Delivery, &order.Payment, order.Items)
	if err != nil {
		log.Println(err)
		return -1, err
	}

	return fmt.Sprint("Generated an order - ", order.Order_uid), nil
}

func getOrderById(id string) (interface{}, error) {
	id = id[1 : len(id)-1]
	result, err := orderFull.FindOne(ctx, globalPool, id)
	if err != nil {
		log.Println(err)
		return -1, err
	}
	return result, nil
}
