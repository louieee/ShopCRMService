package rabbitmq

import (
	"encoding/json"
	"fmt"
	"ShopService/core/rabbitmq/helpers"
	"github.com/go-redis/redis"
	_ "github.com/joho/godotenv"
	"github.com/streadway/amqp"
)


type Server struct {
	url         string
	redisUrl    string
	connection  *amqp.Connection
	redisClient *redis.Client
	channel     *amqp.Channel
	exchange    helpers.Exchange
}

var RabbitMQServer = Server{url: "",
	redisUrl:  "",
	exchange:  helpers.Exchange{Name: "sales_app", Type: "fanout"},
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
	return server
}
func (server Server) Publish(queues []string, message helpers.Payload) {
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

func (server Server) Consume(consumers [] helpers.Consumer) {
	t := map[string]interface{}{"type": "quorum"}
	//server.consumers[1].Consume(server.channel)
	for _, consumer := range consumers {
		_, err := server.channel.QueueDeclarePassive(consumer.GetQueue(),
			true, false, false, false, t)
		if err != nil {
			panic(fmt.Sprintf("RabbitMQ: %s", err.Error()))
		}
		err = server.channel.QueueBind(consumer.GetQueue(), "", server.exchange.Name, false, nil)
		if err != nil {
			panic(fmt.Sprintf("RabbitMQ: %s", err.Error()))
		}
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

func (server Server) StoreMessage(key string, data interface{}) *string {
	err := server.redisClient.Set(key, data, 0).Err()
	if err != nil {
		var message = err.Error()
		return &message
	}
	return nil
}
