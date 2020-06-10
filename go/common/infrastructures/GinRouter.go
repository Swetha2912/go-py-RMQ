package infrastructures

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"sample_rmq/common/interfaces"
)

var JwtSecret = "tericsoft"

type GinRouter struct {
	endpoints   map[string]func(req interfaces.IRouterRequest)
	Engine      *gin.Engine
	PortNumber  string
	Config      GinConfig
	Middlewares []func(request interfaces.IRouterRequest)
}

type GinConfig struct {
	EnableSockets bool
	PortNumber    string
}

func GinSocketPageHandler(req interfaces.IRouterRequest) {
	ginReq := req.(*GinRequest)
	ginReq.GinContext.HTML(200, "index.html", nil)
}

func GinSocketHandler(req interfaces.IRouterRequest) {
	var wsupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	ginReq := req.(*GinRequest)
	conn, err := wsupgrader.Upgrade(ginReq.GinContext.Writer, ginReq.GinContext.Request, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		_ = conn.WriteMessage(t, msg)
	}
}

func (broker *GinRouter) Init(_config interface{}) error {
	config := _config.(GinConfig)
	broker.Engine = gin.Default()
	broker.Config = config
	broker.PortNumber = config.PortNumber
	if config.EnableSockets {
		broker.Engine.LoadHTMLFiles("index.html")
	}
	return nil
}

func (broker *GinRouter) Listen(_conn interface{}) error {
	if _conn != nil {
		connString := _conn.(string)
		broker.PortNumber = connString
	}
	broker.Engine.Use(broker.endpointsMiddleware()) // inject custom middleware level and take over requests
	err := broker.Engine.Run(broker.PortNumber)
	return err
}

func (broker *GinRouter) RegisterEndpoints(_endpoints interface{}) error {
	endpoints := _endpoints.(map[string]interfaces.RequestHandler)
	broker.endpoints = endpoints
	if broker.Config.EnableSockets {
		broker.endpoints["GET/ws"] = GinSocketHandler
		broker.endpoints["GET/socketdemo"] = GinSocketPageHandler
	}
	return nil
}

func (broker *GinRouter) RegisterMiddleware(_middlewareHandlers interface{}) error {
	broker.Middlewares = _middlewareHandlers.([]func(request interfaces.IRouterRequest))
	//for _, middleware := range broker.Middleware {
	//	broker.Engine.Use(func(c *gin.Context) {
	//		request := GinRequest{GinContext: c}
	//		middleware(request)
	//	})
	//}
	return nil
}

func (broker *GinRouter) endpointsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := &GinRequest{GinContext: c}
		if c.Request.Method == "POST" {
			var payloadBody map[string]interface{}
			_ = c.BindJSON(&payloadBody)
			request.Body = payloadBody
		}
		for route, handler := range broker.endpoints {
			if route == c.Request.Method+c.Request.RequestURI {
				for _, middleware := range broker.Middlewares {
					if !request.IsEnded() {
						middleware(request)
					}
				}
				if !request.IsEnded() {
					handler(request)
				}
			}
		}
	}
}

func (broker *GinRouter) Call(pattern interface{}, payload interface{}) (bool, interfaces.RouterResponse) {
	return true, interfaces.RouterResponse{}
}

func ProvideGinRouter() *GinRouter {
	return &GinRouter{}
}


func (broker *GinRouter) Invoke(pattern interface{}, payload interface{}) (bool, []byte) {
	return true, []byte("chi : default response")
}
