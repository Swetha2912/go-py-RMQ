//+build !mux
//+build !gin

package infrastructures

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sample_rmq/common/interfaces"
)

type ChiRequest struct {
	HttpResponse http.ResponseWriter
	HttpRequest  *http.Request
	Body         map[string]interface{}
	NextFunc     http.HandlerFunc
}

func (request *ChiRequest) GetID() string {
	return "1234"
}

func (request *ChiRequest) GetBody() map[string]interface{} {
	if request.Body["*"] == nil {
		request.Body = map[string]interface{}{}

		var jsonPayload map[string]interface{}
		bodyBytes, _ := ioutil.ReadAll(request.HttpRequest.Body)
		err := json.Unmarshal(bodyBytes, &jsonPayload)
		fmt.Println("error json payload", err)
		_ = request.HttpRequest.Body.Close() //  must close
		request.HttpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		for key, val := range jsonPayload {
			request.Body[key] = val
		}

		err = request.HttpRequest.ParseForm()
		if err == nil {
			for key, values := range request.HttpRequest.PostForm {
				request.Body[key] = values[0]
			}
		} else {
			fmt.Print("error", err)
		}
		request.Body["*"] = "parse_complete"
	}

	return request.Body
}

func (request *ChiRequest) GetHeaders(key string) interface{} {
	return request.HttpRequest.Header.Get(key)
	//headers := make(map[string]interface{})
	//requestData := request.HttpRequest
	//requestHeaders := requestData.Header
	//for key, value := range requestHeaders {
	//	headers[key] = value[0]
	//}
	//return headers
}

func (request *ChiRequest) IsEnded() bool {
	return false
}

func (request *ChiRequest) ReplyBack(responseCode int, responseBody interfaces.RouterResponse) bool {
	bytesResponse, _ := json.Marshal(responseBody)
	request.HttpResponse.WriteHeader(responseCode)
	_, _ = request.HttpResponse.Write(bytesResponse)
	return true
}

func (request *ChiRequest) Next() {
	request.NextFunc(request.HttpResponse, request.HttpRequest)
}

func (request *ChiRequest) SetContextData(key string, data interface{}) {
	ctx := context.WithValue(request.HttpRequest.Context(), key, data)
	request.HttpRequest = request.HttpRequest.WithContext(ctx)
}

func (request *ChiRequest) GetContextData(key string) interface{} {
	data := request.HttpRequest.Context().Value(key)

	return data
}

func (request *ChiRequest) GetParam(key string) interface{} {
	return request.HttpRequest.URL.Query().Get(key)
}

func (request *ChiRequest) GetEndpoint() string {
	return ""
}

func (request *ChiRequest) GetFile(key string) (interface{}, interface{}) {
	file, fileHeader, _ := request.HttpRequest.FormFile(key)
	return file, fileHeader
}
