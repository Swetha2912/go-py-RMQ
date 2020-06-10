package infrastructures

import (
	"sample_rmq/common/interfaces"

	"github.com/gofiber/fiber"
)

func FiberEndpointWrapper(handler func(request interfaces.IRouterRequest)) func(*fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		fRequest := FiberRequest{
			FiberCtx: ctx,
		}
		handler(&fRequest)
	}
}
