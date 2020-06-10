// +build mux

package infrastructures

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/mux"
	"net/http"
	"sample_rmq/common/interfaces"
	"time"
)

type MuxRouter struct {
	endpoints   []MuxEndpoint
	Engine      *mux.Router
	ListenOn    string
	Config      MuxConfig
	Middlewares []func(request interfaces.IRouterRequest)
}

type MuxConfig struct {
	EnableSockets bool
	ListenOn      string
}

func MuSocketPageHandler(req interfaces.IRouterRequest) {
	//ginReq := req.(*MuxRequest)
	//ginReq.MuxContext.HTML(200, "index.html", nil)
}

func MuxSocketHandler(req interfaces.IRouterRequest) {
	//var wsUpgrader = websocket.Upgrader{
	//	ReadBufferSize:  1024,
	//	WriteBufferSize: 1024,
	//}
	//ginReq := req.(*MuxRequest)
	//conn, err := wsUpgrader.Upgrade(ginReq.MuxContext.Writer, ginReq.MuxContext.HttpRequest, nil)
	//if err != nil {
	//	fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
	//	return
	//}
	//
	//for {
	//	t, msg, err := conn.ReadMessage()
	//	if err != nil {
	//		break
	//	}
	//	_ = conn.WriteMessage(t, msg)
	//}
}

func (broker *MuxRouter) Init(_config interface{}) error {
	config := _config.(MuxConfig)
	broker.Engine = mux.NewRouter()
	broker.Config = config
	broker.ListenOn = config.ListenOn
	if config.EnableSockets {
		//broker.Engine.PathPrefix("index.html").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	}
	return nil
}

func (broker *MuxRouter) Listen(_conn interface{}) error {
	fmt.Println("Starting Mux router ...")
	if _conn != nil {
		connString := _conn.(string)
		broker.ListenOn = connString
	}
	//broker.Engine.Use(broker.endpointsMiddleware()) // inject custom middleware level and take over requests
	//err := broker.Engine.Run(broker.ListenOn)

	srv := &http.Server{
		Handler: broker.Engine,
		Addr:    broker.ListenOn,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	return err
}

func (broker *MuxRouter) RegisterEndpoints(endpoints interface{}) error {
	broker.endpoints = endpoints.([]MuxEndpoint)
	for _, endpoint := range broker.endpoints {
		fmt.Println("url : ", endpoint.URL, " function() : ", endpoint.Handler)
		broker.Engine.Path(endpoint.URL).Methods(endpoint.Method).HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			endpoint.Handler(&MuxRequest{ResponseWriter: writer, Request: request})
			return
		})
	}
	return nil
}

func (broker *MuxRouter) RegisterMiddleware(_middlewareHandlers interface{}) error {
	broker.Middlewares = _middlewareHandlers.([]func(request interfaces.IRouterRequest))
	//for _, middleware := range broker.Middleware {
	//	broker.Engine.Use(func(c *gin.Context) {
	//		request := MuxRequest{MuxContext: c}
	//		middleware(request)
	//	})
	//}
	return nil
}

func (broker *MuxRouter) endpointsMiddleware() gin.HandlerFunc {
	//return func(c *gin.Context) {
	//	request := &MuxRequest{MuxContext: c}
	//	if c.HttpRequest.Method == "POST" {
	//		var payloadBody map[string]interface{}
	//		_ = c.BindJSON(&payloadBody)
	//		request.Body = payloadBody
	//	}
	//	for route, handler := range broker.endpoints {
	//		if route == c.HttpRequest.Method+c.HttpRequest.RequestURI {
	//			for _, middleware := range broker.Middleware {
	//				if !request.IsEnded() {
	//					middleware(request)
	//				}
	//			}
	//			if !request.IsEnded() {
	//				handler(request)
	//			}
	//		}
	//	}
	//}

	return nil
}

func (broker *MuxRouter) Call(pattern interface{}, payload interface{}) (bool, interfaces.RouterResponse) {
	return true, interfaces.RouterResponse{}
}

func ProvideMuxRouter() *MuxRouter {
	return &MuxRouter{}
}
