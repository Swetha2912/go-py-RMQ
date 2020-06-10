// +build mux

package infrastructures

import (
	"encoding/json"
	"net/http"
	"sample_rmq/common/interfaces"
)

type MuxEndpoint struct {
	Method  string
	URL     string
	Handler func(request interfaces.IRouterRequest)
}

type MuxRequest struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	Body           map[string]interface{}
	MatchedRoute   string
	RequestEnded   bool
}

func (request *MuxRequest) GetBody() map[string]interface{} {
	return request.Body
}

func (request *MuxRequest) IsEnded() bool {
	return request.RequestEnded
}

func (request *MuxRequest) ReplyBack(responseCode int, responseBody interfaces.RouterResponse) bool {
	bytesResponse, _ := json.Marshal(responseBody)
	request.ResponseWriter.WriteHeader(responseCode)
	_, _ = request.ResponseWriter.Write(bytesResponse)
	request.RequestEnded = true
	return true
}

func (request *MuxRequest) SetData(key string, data interface{}) {
	//request.MuxContext.Set(key, data)
}

func (request *MuxRequest) GetData(key string) interface{} {
	//data, _ := request.MuxContext.Get(key)
	//return data
	return nil
}

func (request *MuxRequest) GetEndpoint() string {
	return ""
}

func (request *MuxRequest) GetFiles(key string) (interface{}, interface{}) {
	return nil, nil
}

func (request *MuxRequest) GetHeaders() map[string]interface{} {
	return nil
}
