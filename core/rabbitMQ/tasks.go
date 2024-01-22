package rabbitMQ

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// this function receives 2 arguments each integers and adds them
func addTask(arguments ...string) {
	if len(arguments) != 2 {
		panic("Invalid Number of arguments")
	}
	first, _ := strconv.Atoi(arguments[0])
	second, _ := strconv.Atoi(arguments[1])
	res := first + second
	server := RabbitMQServer.Connect()
	err := server.redisClient.Set("task_queue.add", res, 0).Err()
	if err != nil {
		panic(fmt.Sprintf("Error storing task result to redis: %s", err.Error()))
	}

}

func addProduct(product Product) *Product {
	res, _ := json.Marshal(product)
	server := RabbitMQServer.Connect()
	err := server.redisClient.Set(fmt.Sprintf("product.%s", product.ID), string(res), 0).Err()
	if err != nil {
		panic(fmt.Sprintf("Error storing product to redis: %s", err.Error()))
	}
	println("added product")
	return nil

}
