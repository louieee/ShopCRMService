package celery

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"time"
)
import "github.com/gocelery/gocelery"

func getCeleryClient() *gocelery.CeleryClient {
	//create broker and backend
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL("redis://localhost:6379")
			if err != nil {
				println(err.Error())
				return nil, err
			}
			return c, err
		},
	}

	celeryBroker := gocelery.NewRedisBroker(redisPool)
	celeryBackend := gocelery.NewRedisBackend(redisPool)

	//////use AMQP instead
	//celeryBroker := gocelery.NewAMQPCeleryBroker("amqp://")
	//celeryBackend := gocelery.NewAMQPCeleryBackend("amqp://")

	// Configure with 2 celery workers
	celeryClient, err := gocelery.NewCeleryClient(celeryBroker, celeryBackend, 4)
	if err != nil {
		println(err.Error())
	}
	return celeryClient
}

var celeryClient *gocelery.CeleryClient = getCeleryClient()

func RunCeleryWorker() {
	registerTasks(celeryClient)
	celeryClient.StartWorker()
	//go celeryClient.StopWorker()
	//SendTask("just_comment")
}

func registerTasks(celeryClient *gocelery.CeleryClient) {
	celeryClient.Register("add", add)
	celeryClient.Register("eat", eat)
}

func SendTask(taskName string, arguments ...interface{}) {
	asyncResult, err := celeryClient.Delay(taskName, arguments...)
	if err != nil {
		panic(err)
	}

	// check if result is ready
	isReady, _ := asyncResult.Ready()
	fmt.Printf("ready status %v\n", isReady)

	// get result with 5s timeout
	res, err := asyncResult.Get(5 * time.Second)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res)
	}
}
