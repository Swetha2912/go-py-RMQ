package infrastructures

import (
	"github.com/gofiber/fiber"
)

type FiberConfig struct{
	ListenOn string
}

type FiberRouter struct{
	Engine *fiber.App
	ListenOn string
	Config FiberConfig
}

func (broker *FiberRouter) Init(fiberRouter interface{}) error {
	broker.Engine = fiberRouter.(*fiber.App)
	return nil
}

func (broker *FiberRouter) Listen(_fiberConfig interface{}) error{
	fiberConfig := _fiberConfig.(FiberConfig)
	broker.Engine.Listen(fiberConfig.ListenOn)
	return nil
}


func (broker *FiberRouter) RegisterEndpoints(endpoints interface{}) error {
	return nil
}

func (broker *FiberRouter) RegisterMiddleware(_middlewareHandlers interface{}) error {
	return nil
}

func (broker *FiberRouter) Invoke(pattern interface{}, payload interface{}) (bool, []byte) {
	return true, []byte("chi : default response")
}

func ProvideFiberRouter() *FiberRouter {
	return &FiberRouter{}
}