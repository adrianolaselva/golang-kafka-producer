package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

const (
	Iterations = 2
	Address = "127.0.0.1:11300"
	Queue = "track-events"
)

func main() {
	conn, _ := beanstalk.Dial("tcp", Address)

	tube := beanstalk.NewTube(conn, Queue)

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

		_, _ = tube.Put(value, uint32(i), 0, 0)
	}
}
