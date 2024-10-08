package consumer

import (
	"encoding/json"
	"fmt"
	"ShopService/core/rabbitmq/helpers"
	"ShopService/core/rabbitmq/handlers"
	_ "github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"strings"
)

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
		var payload helpers.Payload
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

func (consumer Consumer) HandleMessage(message helpers.Payload) {
	data := strings.ReplaceAll(message.Data, `'`, `"`)
	switch consumer.queue {
	case "crm_queue":
		switch message.DataType {
		case "user":
			handlers.HandleUser(message.Action, message.Data)
		}
		println(fmt.Sprintf("[Y] %s", data))
	default:
		println(fmt.Sprintf("[X] %s", data))

	}
}

func (consumer Consumer) GetQueue() string {
	return consumer.queue
}

func (consumer *Consumer) SetQueue(queue string) {
    consumer.queue = queue
}
