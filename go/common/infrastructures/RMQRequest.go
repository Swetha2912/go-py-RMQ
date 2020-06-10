package infrastructures

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"sample_rmq/common/interfaces"
)

type RMQRequest struct {
	RMQCh         *amqp.Channel
	RMQConn       *amqp.Connection
	Body          map[string]interface{}
	CorrelationID string
	MatchedRoute  string
	ReplyTo       string
	RequestEnded  bool
}

func (request *RMQRequest) GetID() string {
	return ""
}

func (request *RMQRequest) IsEnded() bool {
	return request.RequestEnded
}

func (request *RMQRequest) Next() {
}

func (request *RMQRequest) GetBody() map[string]interface{} {
	return request.Body
}

// ReplyBack responds back to the request
func (request *RMQRequest) ReplyBack(responseCode int, responseBody interfaces.RouterResponse) bool {
	bytes, _ := json.Marshal(responseBody)
	request.RequestEnded = true
	_ = request.RMQCh.Publish(
		"",              // Exchange global
		request.ReplyTo, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			CorrelationId: request.CorrelationID,
			Body:          bytes,
		})
	return true
}

func (request *RMQRequest) GetHeaders(key string) interface{} {
	return nil
}

func (request *RMQRequest) GetEndpoint() string {
	return ""
}

func (request *RMQRequest) SetContextData(key string, data interface{}) {
}

func (request *RMQRequest) GetContextData(key string) interface{} {
	return nil
}

func (request *RMQRequest) GetParam(key string) interface{} {
	return nil
}

func (request *RMQRequest) GetFile(key string) (interface{}, interface{}) {
	return nil, nil
}
