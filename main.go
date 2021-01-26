package main

import (
	"fmt"
	"strings"
)

const (
	Iterations = 1000
)

func main() {
	topics := strings.Split("topic1,topic2", ",")

	fmt.Println(topics)

	//rand.Seed(42)
	//answers := []string{
	//	"It is certain",
	//	"It is decidedly so",
	//	"Without a doubt",
	//	"Yes definitely",
	//	"You may rely on it",
	//	"As I see it yes",
	//	"Most likely",
	//	"Outlook good",
	//	"Yes",
	//	"Signs point to yes",
	//	"Reply hazy try again",
	//	"Ask again later",
	//	"Better not tell you now",
	//	"Cannot predict now",
	//	"Concentrate and ask again",
	//	"Don't count on it",
	//	"My reply is no",
	//	"My sources say no",
	//	"Outlook not so good",
	//	"Very doubtful",
	//}
	//
	////var price float32 = 100
	////var percentage float32 = 10
	////price -= price * (percentage / 100)
	////
	////fmt.Println(price)
	//
	//data := []map[string]interface{}{}
	//for i := 0; i < Iterations;i++ {
	//	data = append(data, map[string]interface{}{
	//		"uuid": fmt.Sprintf("%x", md5.Sum([]byte(string(time.Now().Format(time.RFC3339Nano))))),
	//		"detail": answers[rand.Intn(len(answers))],
	//		"detail1": answers[rand.Intn(len(answers))],
	//		"detail2": answers[rand.Intn(len(answers))],
	//		"detail3": answers[rand.Intn(len(answers))],
	//		"detail4": answers[rand.Intn(len(answers))],
	//		"detail5": answers[rand.Intn(len(answers))],
	//	})
	//}
	//
	//fmt.Printf("size: %v\n", len(data))
	//fmt.Printf("size: %v\n", unsafe.Sizeof(data))
	//
	//result, _ := yaml.Marshal(data)
	//
	//fmt.Printf("size: %v\n", len(result))
	//
	//bytes, _ := mapLength(data)
	//fmt.Printf("size: %v\n", len(bytes))

}

func mapLength(v interface{}) ([]byte, error) {
	if b, ok := v.(*[]byte); ok {
		fmt.Println(b)
		return *b, nil
	}
	return nil, nil
}


