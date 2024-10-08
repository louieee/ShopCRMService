package tasks

import (
	"fmt"
	"github.com/joho/godotenv"
	"ShopService/core/rabbitmq"
	"log"
	"os"
	"strconv"
)

// this function receives 2 arguments each integers and adds them
func addTask(arguments ...string) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	if len(arguments) != 2 {
		panic("Invalid Number of arguments")
	}
	first, _ := strconv.Atoi(arguments[0])
	second, _ := strconv.Atoi(arguments[1])
	res := first + second
	server := rabbitmq.RabbitMQServer.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	var result = server.StoreMessage("task_queue.add", res)
	if result != nil {
		panic(fmt.Sprintf("Error storing task result to redis: %s", *result))
	}

}
