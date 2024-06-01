package rabbitMQ

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

type PayloadString struct {
	Message string
}
type Task struct {
	Function string
	Args     []string
}

type Product struct {
	ID       string
	Name     string
	Quantity uint
}

type Payload struct {
	Msg      string
	Task     Task
	Product  Product
	DataType string
}

type Message struct {
	Queue   string
	Payload Payload
}
type Server struct {
	url         string
	redisUrl    string
	connection  *amqp.Connection
	redisClient *redis.Client
}

var RabbitMQServer = Server{url: "amqp://guest:guest@localhost:5672/", redisUrl: "localhost:6379"}

func (server Server) Connect() Server {
	server.redisClient = redis.NewClient(&redis.Options{
		Addr: server.redisUrl,
	})

	rabbitMQConnection, err := amqp.Dial(server.url)
	if err != nil {
		panic("Could not connect to rabbit MQ")
	}
	server.connection = rabbitMQConnection
	return server
}
func (server Server) Publish(message Message) {
	ch, _ := server.connection.Channel()

	q, _ := ch.QueueDeclare(
		message.Queue, // name
		true,          // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	bytesToSend, _ := json.Marshal(message.Payload)

	_ = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        bytesToSend})

}

func (server Server) Consume(queue string) {
	ch, _ := server.connection.Channel()
	q, _ := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	for {
		msg := <-msgs
		var payload Payload
		err1 := json.Unmarshal(msg.Body, &payload)
		if err1 != nil {
			panic(fmt.Sprintf("Error unwrapping message: %s", err1.Error()))
		}
		// process task and store its result
		server.processPayload(queue, payload)
		println("received: ", payload.Msg)
		err2 := msg.Ack(false)
		if err2 != nil {
			panic("Error Acknowledging task")
		}
	}
}
func (server Server) ConsumeTask() {
	ch, _ := server.connection.Channel()
	err := ch.ExchangeDeclarePassive(
		"business",
		"fanout",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
		return
	}
	q, _ := ch.QueueDeclarePassive(
		"task", // name
		true,   // durable
		false,  // delete when unused
		false,  // exclusive
		false,  // no-wait
		nil,    // arguments
	)
	err = ch.QueueBind(
		"task",
		"",
		"business",
		false, nil)
	if err != nil {
		panic(err)
		return
	}
	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	for {
		msg := <-msgs
		var payload PayloadString
		err1 := json.Unmarshal(msg.Body, &payload)
		if err1 != nil {
			//msg.Ack(false)
			panic(fmt.Sprintf("Error unwrapping message: %s", err1.Error()))
			//return
		}
		println("received: ", payload.Message)
		// process task and store its result
		err2 := msg.Ack(false)
		if err2 != nil {
			panic("Error Acknowledging task")
		}
	}
}
func (server Server) processPayload(queue string, msg Payload) {
	switch msg.DataType {
	case "Task":
		msgData := msg.Task
		function := taskFunctions()[msgData.Function]
		function(msgData.Args...)
		res := server.RetrieveMessage(fmt.Sprintf("%s.%s", queue, msgData.Function))
		println("The result is ", res)

	case "Message":
		msgData := msg.Msg
		println("Msg is ", msgData)
		return

	case "Product":
		msgData := msg.Product
		function := productFunctions()[msg.Msg]
		function(msgData)
		res := server.RetrieveMessage(fmt.Sprintf("product.%s", msgData.ID))
		println("The result is ", res)
	default:
		println("The message", msg.Msg)
	}

}

type taskFunctionType func(arguments ...string)
type productFunctionType func(product Product) *Product

func taskFunctions() map[string]taskFunctionType {
	queues := make(map[string]taskFunctionType)
	queues["add"] = addTask
	return queues
}
func productFunctions() map[string]productFunctionType {
	queues := make(map[string]productFunctionType)
	queues["add"] = addProduct
	return queues
}

func (server Server) RetrieveMessage(key string) string {
	result, err := server.redisClient.Get(key).Result()
	if err != nil {
		panic("Error retrieving message result")
	}
	println("Result From Redis: ", result)
	return result
}
