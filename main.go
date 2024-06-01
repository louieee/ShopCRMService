package main

import (
	"ShopService/core"
	"ShopService/core/celery"
	"ShopService/core/rabbitMQ"
	_ "ShopService/docs"
	"ShopService/routers"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// @title           Shop Service
// @version         1.0
// @description     A shop service used to manage sales
// @termsOfService  https://tos.santoshk.dev

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
//@securityDefinitions.apikey	Bearer
//	@in							header
//	@name						Authorization
//	@description				This is used to authorize the authentication

// @host      localhost:8080
// @BasePath  /api/v1
// @Security Bearer
func main() {
	//appType := os.Getenv("APP_TYPE")
	var appType string
	// Register flags
	flag.StringVar(&appType, "appType", "", "Determines which server to run")
	flag.Parse()
	switch appType {
	case "server":
		// Run server-specific logic
		runGoServer()
	case "celery worker":
		// Run worker-specific logic
		celery.RunCeleryWorker()
	case "websocket":
		core.RunWebsocketServer()
	case "rabbitMQ":
		rabbitMQServer()
	default:
		//runGoServer()
		rabbitMQServer()
	}
}

func runGoServer() {
	db := core.GetDB()
	router := gin.Default()
	routers.Routes(router, db)
	defer func(db *gorm.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	// Start the server
	port := core.ServerConfig["PORT"]
	fmt.Printf("Server is running on :%d\n", port)
	err := router.Run(fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
}

func rabbitMQServer() {
	server := rabbitMQ.RabbitMQServer
	server = server.Connect()
	println("connected to rabbit mq")

	//task := rabbitMQ.Message{
	//	Queue: "task_queue", Payload: rabbitMQ.Payload{
	//		DataType: "Task",
	//		Task: rabbitMQ.Task{
	//			Function: "add",
	//			Args:     []string{"1", "9"},
	//		},
	//	},
	//}
	//message := rabbitMQ.Message{
	//	Queue: "message_queue", Payload: rabbitMQ.Payload{
	//		DataType: "Message",
	//		Msg:      "Hello Rabbit MQ Server",
	//	},
	//}
	product := rabbitMQ.Message{
		Queue: "product_queue", Payload: rabbitMQ.Payload{
			Msg:      "add",
			DataType: "Product",
			Product: rabbitMQ.Product{
				ID:       "#ERTY",
				Name:     "Box",
				Quantity: 5,
			},
		},
	}
	server.Publish(product)
	server.ConsumeTask()
}
