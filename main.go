package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"time"
)

const (
	Iterations = 2
	TopicName = "track-events"
	Brokers = "127.0.0.1:9092"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("file .env not found")
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": Brokers})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	topic := TopicName
	for i := 1; i<=Iterations;i++ {
		value, _ := json.Marshal(map[string]interface{}{
			"userId": fmt.Sprintf("%x", md5.Sum([]byte(string(i)))),
			"anonymousId": nil,
			"integrations": map[string]interface{}{
				"all": true,
				"mixpanel": true,
				"firebase": true,
				"googleAnalytics": true,
				"kinesisFirehose": true,
			},
			"traits": map[string]interface{}{
				"email": "teste@picpay.comn",
				"first_name": "firstname",
				"last_name": "lastname",
				"username": "username",
				"created_at": time.Now().Format(time.RFC3339),
			},
			"context": map[string]interface{}{
				"ip": "127.0.0.1",
			},
			"timestamp": time.Now().Format(time.RFC3339),
		})

		_ = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value,
		}, nil)
	}

	p.Flush(1)
}

