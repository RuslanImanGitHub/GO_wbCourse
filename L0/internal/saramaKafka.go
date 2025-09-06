package sarKaf

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
)

type SarKafCont struct {
	Producer sarama.SyncProducer
	Consumer sarama.Consumer
}

func NewProducer(brokers []string) (sarama.SyncProducer, error) {
	cfg := sarama.NewConfig()

	cfg.Net.ReadTimeout = 10 * time.Second
	cfg.Net.WriteTimeout = 10 * time.Second

	cfg.Producer.Return.Successes = true
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Retry.Max = 5
	cfg.Producer.Return.Errors = true

	log.Println("Created new Producer")

	return sarama.NewSyncProducer(brokers, cfg)
}

func NewConsumer(broker []string) (sarama.Consumer, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Errors = true

	log.Println("Created new Consumer")

	return sarama.NewConsumer(broker, cfg)
}

func PushMsg(prod sarama.SyncProducer, msg []byte, topic string) error {
	kafMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}
	partiotion, offset, err := prod.SendMessage(kafMsg)
	if err != nil {
		log.Println(err)
		return err
	}

	log.Printf("Message sent: %s / Topic: %s, Partition: %d, Offset: %d\n", kafMsg.Value, topic, partiotion, offset)
	return nil
}

func EndlessServerLoop(prod sarama.SyncProducer, cons sarama.Consumer,
	finisherCahn chan struct{}, function func(string) (interface{}, error),
	reqTopic string, resTopic string) {

	partConsumer, err := cons.ConsumePartition(reqTopic, 0, sarama.OffsetNewest)
	if err != nil {
		log.Println(err)
		return
	}
	go func() {
		for {
			select {
			case err := <-partConsumer.Errors():
				cErr := err.Err
				if cErr != nil {
					log.Println(cErr)
					errBytes, err := json.Marshal(cErr)
					if err != nil {
						log.Println(err)
						continue
					}
					err = PushMsg(prod, errBytes, resTopic)
					if err != nil {
						log.Println(err)
						continue
					}
				}
			case msg := <-partConsumer.Messages():
				log.Printf("Message received: %s / Topic: %s, Partition: %d \n", string(msg.Value), string(msg.Topic), int(msg.Partition))
				result, err := function(string(msg.Value))
				if err != nil {
					log.Println(err)
					result = err.Error()
				}

				resultBytes, err := json.Marshal(result)
				if err != nil {
					log.Println(err)
					continue
				}

				err = PushMsg(prod, resultBytes, resTopic)
				if err != nil {
					log.Println(err)
					continue
				}
			}
		}
	}()

	<-finisherCahn
}

func ControllerRequest(prod sarama.SyncProducer, cons sarama.Consumer, value []byte, reqTopic string, resTopic string) string {
	err := PushMsg(prod, value, reqTopic)
	if err != nil {
		return err.Error()
	}

	partConsumer, err := cons.ConsumePartition(resTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return err.Error()
	}

	var response string
	doneChan := make(chan struct{})

	go func() {
		select {
		case err := <-partConsumer.Errors():
			response = err.Error()
			partConsumer.Close()
			doneChan <- struct{}{}

		case msg := <-partConsumer.Messages():
			log.Printf("Message received: %s / Topic: %s, Partition: %d \n", string(msg.Value), string(msg.Topic), int(msg.Partition))

			response = string(msg.Value)
			partConsumer.Close()
			doneChan <- struct{}{}
		}
	}()
	<-doneChan
	return response
}
