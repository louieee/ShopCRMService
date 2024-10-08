package helpers

import "github.com/streadway/amqp"



type Consumer interface {
		Consume(ch *amqp.Channel)
		HandleMessage(message Payload)
		GetQueue() string
}



type Payload struct {
	Action   string `json:"action"`
	DataType string `json:"data_type"`
	Data     string `chanHandlerjson:"data"`
}

type Exchange struct {
	Name string
	Type string
}
