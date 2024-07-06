package rabbitMQ

import (
	"ShopService/core/rabbitMQ/consumers"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	_ "github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"strings"
)

type Payload struct {
	Action   string `json:"action"`
	DataType string `json:"data_type"`
	Data     string `json:"data"`
}

type Exchange struct {
	Name string
	Type string
}

type Consumer struct {
	queue string
}

func (consumer Consumer) Consume(ch *amqp.Channel) {
	msgs, _ := ch.Consume(
		consumer.queue, // queue
		"",             // consumer
		false,          // auto ack
		false,          // exclusive
		false,          // no local
		false,          // no wait
		nil,            // args
	)

	for {
		msg := <-msgs
		var payload Payload
		err1 := json.Unmarshal(msg.Body, &payload)
		if err1 != nil {
			panic(fmt.Sprintf("Error unwrapping message: %s", err1.Error()))
		}
		// process task and store its result
		consumer.HandleMessage(payload)
		println("received: ", payload.Data)
		err2 := msg.Ack(false)
		if err2 != nil {
			panic("Error Acknowledging task")
		}
	}
}

func (consumer Consumer) HandleMessage(message Payload) {
	data := strings.ReplaceAll(message.Data, `'`, `"`)
	switch consumer.queue {
	case "crm_queue":
		switch message.DataType {
		case "lead":
			consumers.HandleLead(message.Action, message.Data)
		}
		println(fmt.Sprintf("[Y] %s", data))
	default:
		println(fmt.Sprintf("[X] %s", data))

	}
}

type Server struct {
	url         string
	redisUrl    string
	connection  *amqp.Connection
	redisClient *redis.Client
	channel     *amqp.Channel
	consumers   []Consumer
	exchange    Exchange
}

var RabbitMQServer = Server{url: "",
	redisUrl:  "",
	exchange:  Exchange{"sales_app", "fanout"},
	consumers: []Consumer{{queue: "crm_queue"}, {queue: "chat_queue"}},
}

func (server Server) Connect(url string, redisUrl string) Server {
	server.url = url
	server.redisUrl = redisUrl
	server.redisClient = redis.NewClient(&redis.Options{
		Addr: server.redisUrl,
	})

	rabbitMQConnection, err := amqp.Dial(server.url)
	if err != nil {
		panic("Could not connect to rabbit MQ")
	}
	server.connection = rabbitMQConnection
	server.channel, _ = server.connection.Channel()
	err = server.channel.ExchangeDeclarePassive(server.exchange.Name,
		server.exchange.Type, true, false, false, false, nil)
	if err != nil {
		panic(fmt.Sprintf("RabbitMQ: %s", err.Error()))
	}
	t := map[string]interface{}{"type": "quorum"}
	for _, consumer := range server.consumers {
		_, err := server.channel.QueueDeclarePassive(consumer.queue,
			true, false, false, false, t)
		if err != nil {
			panic(fmt.Sprintf("RabbitMQ: %s", err.Error()))
		}
		err = server.channel.QueueBind(consumer.queue, "", server.exchange.Name, false, nil)
		if err != nil {
			panic(fmt.Sprintf("RabbitMQ: %s", err.Error()))
		}
	}
	return server
}
func (server Server) Publish(queues []string, message Payload) {
	ch, _ := server.connection.Channel()
	bytesToSend, _ := json.Marshal(message)
	_ = ch.Publish(
		server.exchange.Name, // exchange
		"",                   // routing key
		true,                 // mandatory
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytesToSend})

}

func (server Server) Consume() {

	//server.consumers[1].Consume(server.channel)
	for _, consumer := range server.consumers {
		go consumer.Consume(server.channel)
	}
}
func (server Server) RetrieveMessage(key string) string {
	result, err := server.redisClient.Get(key).Result()
	if err != nil {
		panic("Error retrieving message result")
	}
	println("Result From Redis: ", result)
	return result
}
