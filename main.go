package main

import (
	"ShopService/core"
	"ShopService/core/celery"
	"ShopService/core/rabbitmq"
	"ShopService/core/rabbitmq/consumer"
	"ShopService/core/rabbitmq/helpers"
	_ "ShopService/docs"
	"ShopService/routers"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"log"
	"os"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
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
		rabbitMQServer()
		runGoServer()

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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	server := rabbitmq.RabbitMQServer
	server = server.Connect(os.Getenv("RABBIT_MQ_HOST"), os.Getenv("REDIS_URL"))
	println("connected to rabbit mq")
	report_consumer :=  consumer.Consumer{}
	report_consumer.SetQueue("report_queue")
	fmt.Printf("Consumer queue: %s \n", report_consumer.GetQueue())
	server.Consume([]helpers.Consumer{report_consumer})
}
