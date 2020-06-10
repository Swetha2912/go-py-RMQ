package main

import (
	"fmt"
	"os"
	"sample_rmq/common/infrastructures"
	"sample_rmq/common/interfaces"
	"sample_rmq/common/utilities"
	"sample_rmq/gateway/controllers"

	"github.com/gofiber/fiber"
	"github.com/joho/godotenv"
)

type Container struct {
	Test          interfaces.ITestController
	DefaultRouter interfaces.IRouter
	rmqBroker     interfaces.IRouter
	zmqBroker     interfaces.IBroker
}

var container Container

func ProvideContainer(test controllers.InferenceController, router *infrastructures.FiberRouter, pythonBroker *infrastructures.RMQRouter, zmqBroker *infrastructures.ZmqBroker) Container {
	return Container{Test: test, DefaultRouter: router, rmqBroker: pythonBroker, zmqBroker: zmqBroker}
}

func main() {

	fmt.Println("hello i am in main")
	//loading .env file
	err1 := godotenv.Load()
	utilities.ExitOnErr(err1, "cannot parse env file")

	container = InitializeContainer()

	// RMQ connection
	rmqConfig := infrastructures.RMQConfig{
		DefaultExchange: "sample_rmq",
		ConnString:      "amqp://" + os.Getenv("RMQ_USER") + ":" + os.Getenv("RMQ_PASS") + "@" + os.Getenv("RMQ_HOST"),
		DefaultQueue:    "",
		ListenPatterns:  []string{"*.*"},
	}
	err := container.rmqBroker.Init(rmqConfig)
	fmt.Println("rmq error", err)

	fiberRouter := fiber.New()
	fiberRouter.Post("/sample/rmq", infrastructures.FiberEndpointWrapper(container.Test.Hello))

	err2 := container.DefaultRouter.Init(fiberRouter)
	utilities.ExitOnErr(err2, "cannot initiate default router")

	fmt.Print("starting app : sample_rmq \n")
	err2 = container.DefaultRouter.Listen(infrastructures.FiberConfig{
		ListenOn: "3000",
	})
	utilities.ExitOnErr(err2, "cannot listen on default router")

}
