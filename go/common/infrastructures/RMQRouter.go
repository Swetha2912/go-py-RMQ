package infrastructures

import (
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/streadway/amqp"
	"sync"
	"sample_rmq/common/interfaces"
	"sample_rmq/common/utilities"
	"time"
)

type RMQRouter struct {
	RMQCh       *amqp.Channel
	RMQConn     *amqp.Connection
	endpoints   map[string]func(req interfaces.IRouterRequest)
	config      RMQConfig
	Replies     map[string]chan []byte
	Middlewares []func(request interfaces.IRouterRequest)
}

var Mutex sync.RWMutex

type RMQConfig struct {
	DefaultExchange string
	ConnString      string
	DefaultQueue    string
	ReplyQueue      string
	ListenPatterns  []string
}

func (broker *RMQRouter) Init(config interface{}) error {
	broker.config = config.(RMQConfig)
	conn, err := amqp.Dial(broker.config.ConnString)
	if err != nil {
		return err
	}

	ch, err := conn.Channel()
	if err != nil {
		return err
	}

	broker.RMQConn = conn
	broker.RMQCh = ch
	broker.Replies = make(map[string]chan []byte)

	err = broker.GetReplies()
	if err != nil {
		return err
	}

	err = broker.RMQCh.ExchangeDeclare(
		broker.config.DefaultExchange,
		"topic",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	return nil
}

func (broker *RMQRouter) GetReplies() error {

	q, err := broker.RMQCh.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	broker.config.ReplyQueue = q.Name

	messages, err := broker.RMQCh.Consume(
		q.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	go func() {
		defer utilities.RecoverPanic()
		for d := range messages {

			//docRply := &specs.CreateDocumentReply{}
			//err = proto.Unmarshal(d.Body, docRply)

			//resp := interfaces.RouterResponse{
			//	StatusCode: 0,
			//	Msg:        "",
			//	Data:       docRply,
			//	Error:      nil,
			//	Errors:     nil,
			//}
			Mutex.Lock()
			if broker.Replies[d.CorrelationId] != nil {
				replyChannel := broker.Replies[d.CorrelationId]
				replyChannel <- d.Body
				delete(broker.Replies, d.CorrelationId)
			}
			Mutex.Unlock()
		}
	}()

	return nil
}

func (broker *RMQRouter) Listen(_patterns interface{}) error {

	fmt.Println("Starting RMQ router ...")

	q, err := broker.RMQCh.QueueDeclare(
		broker.config.DefaultQueue,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for _, pattern := range broker.config.ListenPatterns {
		err = broker.RMQCh.QueueBind(
			q.Name,                        // queue name
			pattern,                       // routing key
			broker.config.DefaultExchange, // exchange
			false,
			nil)
		if err != nil {
			return err
		}
	}

	messages, err := broker.RMQCh.Consume(
		q.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return err
	}

	requests := make(chan *RMQRequest)
	//forever := make(chan bool)

	go func() {
		defer utilities.RecoverPanic()
		for d := range messages {
			var reqBody utilities.Payload
			err := utilities.ConvertMap(d.Body, &reqBody)

			if err != nil {
				return
			}

			// new request
			request := &RMQRequest{
				MatchedRoute:  d.RoutingKey,
				Body:          reqBody,
				CorrelationID: d.CorrelationId,
				RMQCh:         broker.RMQCh,
				RMQConn:       broker.RMQConn,
				ReplyTo:       d.ReplyTo,
			}
			requests <- request
		}
	}()

	go func() {
		for req := range requests {
			if broker.endpoints[req.MatchedRoute] != nil {
				fmt.Println("incoming request --> " + req.MatchedRoute)
				handler := broker.endpoints[req.MatchedRoute]
				for _, middleware := range broker.Middlewares {
					if !req.IsEnded() {
						middleware(req)
					}
				}
				if !req.IsEnded() {
					handler(req)
				}
			}
		}
	}()

	fmt.Println("listening on RabbitMQ ..")
	//<-forever
	return nil
}

func (broker *RMQRouter) RegisterEndpoints(endpoints interface{}) error {
	broker.endpoints = endpoints.(map[string]interfaces.RequestHandler)
	return nil
}

func (broker *RMQRouter) RegisterMiddleware(_middlewareHandlers interface{}) error {
	broker.Middlewares = _middlewareHandlers.([]func(request interfaces.IRouterRequest))
	return nil
}

func GetUUID() string {
	u2, _ := uuid.NewV4()
	return u2.String()
}

func (broker *RMQRouter) Invoke(_pattern interface{}, _body interface{}) (bool, []byte) {
	pattern := _pattern.(string)

	//body := _body.(specs.CreateDocumentMessage)
	//bytes, _ := proto.Marshal(body)

	bytes := _body.([]byte)

	correlationID := GetUUID()
	replyWaiter := make(chan []byte)

	Mutex.Lock()
	broker.Replies[correlationID] = replyWaiter
	Mutex.Unlock()

	_ = broker.RMQCh.Publish(
		broker.config.DefaultExchange,
		pattern,
		false,
		false,
		amqp.Publishing{
			CorrelationId: correlationID,
			Body:          bytes,
			ReplyTo:       broker.config.ReplyQueue,
		})

	select {
	case reply := <-replyWaiter:
		return true, reply

	case <-time.After(1000 * time.Second):
		Mutex.Lock()
		close(broker.Replies[correlationID])
		delete(broker.Replies, correlationID)
		Mutex.Unlock()
		return false, []byte("timeout")
	}

}

func ProvideRMQRouter() *RMQRouter {
	return &RMQRouter{}
}
