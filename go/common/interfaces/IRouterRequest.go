package interfaces

type IRouterRequest interface {
	GetBody() map[string]interface{} // gets payload from the request object
	ReplyBack(int, RouterResponse) bool
	GetEndpoint() string
	SetContextData(string, interface{})
	GetContextData(string) interface{}
	GetParam(string) interface{}
	IsEnded() bool
	Next()
	GetFile(string) (interface{}, interface{})
	GetHeaders(string) interface{}
	GetID() string
}

type RouterResponse struct {
	StatusCode int         `bson:"http_code" json:"http_code"`
	Msg        string      `bson:"msg" json:"msg"`
	Data       interface{} `bson:"data" json:"data"`
	Error      interface{} `bson:"error" json:"error"`
	Errors     interface{} `bson:"errors" json:"errors"`
}

type RequestHandler = func(req IRouterRequest)
