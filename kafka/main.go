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
	Iterations = 20
	TopicName = "identify-events-approved"
	Brokers = "127.0.0.1:9092"
	//Brokers = "127.0.0.1:9092"
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
			"userId": i,
			"anonymousId": fmt.Sprintf("%x", md5.Sum([]byte(string(i*1000)))),
			"keyType": "CONSUMER",
			"integrations": map[string]interface{}{
				"all": true,
				"mixpanel": true,
			},
			"traits": map[string]interface{}{
				"appsflyer_id": fmt.Sprintf("%x", md5.Sum([]byte(string(i)))),
				"advertising_id": fmt.Sprintf("%x", md5.Sum([]byte(string(time.Now().Format(time.RFC3339))))),
				"os": "ios",
				"document": "333.111.222-32",
				"cnpj": "33.111.222/0001-32",
				"email": "adrianolaselva@gmail.com",
				"cvv": "335",
				"credit_card": "5232 9609 5260 4191",
				"card": "5232960952604191",
			},
			"context": map[string]interface{}{
				"ip": "127.0.0.1",
				"appsflyer_id": fmt.Sprintf("%x", md5.Sum([]byte(string(i)))),
				"advertising_id": fmt.Sprintf("%x", md5.Sum([]byte(string(time.Now().Format(time.RFC3339))))),
				"idfa": fmt.Sprintf("%x", md5.Sum([]byte(string(time.Now().Format(time.RFC3339))))),
			},
			"timestamp": time.Now().Format(time.RFC3339),
		})

		//value, _ := json.Marshal(map[string]interface{}{
		//	"userId": fmt.Sprintf("%x", md5.Sum([]byte(string(i)))),
		//	"anonymousId": nil,
		//	"integrations": map[string]interface{}{
		//		"all": true,
		//		"mixpanel": true,
		//		"firebase": true,
		//		"googleAnalytics": true,
		//		"kinesisFirehose": true,
		//	},
		//	"traits": map[string]interface{}{
		//		"email": "teste@picpay.comn",
		//		"first_name": "firstname",
		//		"last_name": "lastname",
		//		"username": "username",
		//		"created_at": time.Now().Format(time.RFC3339),
		//	},
		//	"context": map[string]interface{}{
		//		"ip": "127.0.0.1",
		//	},
		//	"timestamp": time.Now().Format(time.RFC3339),
		//})

		//value := []byte(PayloadAndroidJSON)

		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Key: []byte(fmt.Sprintf("%x", md5.Sum(value))),
			Value:          value,
		}, nil)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	p.Flush(100000)
}

const PayloadAndroidJSON = `{
		"anonymousId": null,
		"userId": 15461789,
		"event": "User Approved",
		"keyType": "CONSUMER",
		"integrations": {
			"appsflyer": true
		},
		"properties": {
			"os": "ios"
		},
		"context": null,
		"createdAt": "2020-11-13T09:00:42-0300",
		"uuid": "5221e763fd1310e535f2533436215ccf"
	}`

//const PayloadAndroidJSON = `{
//		"userId": "8453297",
//		"event": "Money Moved",
//		"keyType": "CONSUMER",
//		"integrations": {
//			"mixpanel": true,
//			"appsflyer": true
//		},
//		"properties": {
//			"amount": 120,
//			"currency": "BRL",
//			"external_id": "1986066232",
//			"from": "22484716",
//			"installments": "",
//			"is_mixed": false,
//			"is_qrcode_readed": "",
//			"latitude": "",
//			"longitude": "",
//			"message": "",
//			"privacy": "",
//			"signature": "",
//			"source": "Internal",
//			"state": "Approved",
//			"to": "Claro",
//			"transaction_id": "",
//			"type": "Digital Goods Mobile Recharged",
//			"wallet": null
//		},
//		"context": {
//			"ip": "127.0.0.1"
//		},
//		"timestamp": "2020-06-02T10:16:42-0300"
//	}`

