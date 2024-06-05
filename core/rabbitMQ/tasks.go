package rabbitMQ

import (
	"fmt"
	"github.com/joho/godotenv"
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
	server := RabbitMQServer.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	err = server.redisClient.Set("task_queue.add", res, 0).Err()
	if err != nil {
		panic(fmt.Sprintf("Error storing task result to redis: %s", err.Error()))
	}

}
