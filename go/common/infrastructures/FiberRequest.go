package infrastructures

import (
	"sample_rmq/common/interfaces"
	"fmt"

	"github.com/gofiber/fiber"
)

type FiberRequest struct {
	FiberCtx *fiber.Ctx
}

func (request *FiberRequest) GetID() string {
	return ""
}

func (request *FiberRequest) GetBody() map[string]interface{} {
	var reqBody map[string]interface{}
	value := request.FiberCtx.FormValue("key2")
	// value2 := request.FiberCtx.Body("key2")
	fmt.Print("ping pong value exist : ", value)
	_ = request.FiberCtx.BodyParser(&reqBody)
	return reqBody
}

func (request *FiberRequest) GetHeaders(key string) interface{} {
	return request.FiberCtx.Get(key)
}

func (request *FiberRequest) IsEnded() bool {
	return false
}

func (request *FiberRequest) ReplyBack(responseCode int, responseBody interfaces.RouterResponse) bool {
	responseBody.StatusCode = responseCode
	request.FiberCtx.JSON(responseBody)
	return true
}

func (request *FiberRequest) Next() {
	request.FiberCtx.Next()
}

func (request *FiberRequest) SetContextData(key string, data interface{}) {
	request.FiberCtx.Locals(key, data)
}

func (request *FiberRequest) GetContextData(key string) interface{} {
	return request.FiberCtx.Locals(key)
}

func (request *FiberRequest) GetParam(key string) interface{} {
	data := request.FiberCtx.Query(key)
	return data
}

func (request *FiberRequest) GetEndpoint() string {
	endPoint := request.FiberCtx.OriginalURL()
	return endPoint
}

func (request *FiberRequest) GetFile(key string) (interface{}, interface{}) {
	file, err := request.FiberCtx.FormFile(key)
	if err == nil {
		// Save file to root directory:
		fmt.Println("%%%%% size", file.Size, file.Filename)
		request.FiberCtx.SaveFile(file, "/go/src/sample_rmq/gateway/uploads/"+file.Filename)
	}
	return file, err
}
