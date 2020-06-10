package infrastructures

import (
	"github.com/gin-gonic/gin"
	"sample_rmq/common/interfaces"
)

type GinRequest struct {
	GinContext   *gin.Context
	Body         map[string]interface{}
	MatchedRoute string
	RequestEnded bool
}

func (request *GinRequest) GetID() string {
	return ""
}

func (request *GinRequest) GetBody() map[string]interface{} {
	return request.Body
}

func (request *GinRequest) IsEnded() bool {
	return request.RequestEnded
}

func (request *GinRequest) ReplyBack(responseCode int, responseBody interfaces.RouterResponse) bool {
	request.GinContext.JSON(responseCode, responseBody)
	request.GinContext.Abort()
	request.RequestEnded = true
	return true
}

func (request *GinRequest) Next() {
}

func (request *GinRequest) SetContextData(key string, data interface{}) {
	request.GinContext.Set(key, data)
}

func (request *GinRequest) GetContextData(key string) interface{} {
	data, _ := request.GinContext.Get(key)
	return data
}

func (request *GinRequest) GetEndpoint() string {
	return ""
}

func (request *GinRequest) GetParam(key string) interface{} {
	return nil
}



func (request *GinRequest) GetFile(key string) (interface{}, interface{}) {
	return nil, nil
}

func (request *GinRequest) GetHeaders(key string) interface{} {
	return nil
}
