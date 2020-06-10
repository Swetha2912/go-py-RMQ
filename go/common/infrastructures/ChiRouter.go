//+build !mux
//+build !mux

package infrastructures

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sample_rmq/common/interfaces"
	"time"
)

type ChiEndpoint struct {
	Method     string
	URL        string
	Handler    func(request interfaces.IRouterRequest)
	Middleware []func(req interfaces.IRouterRequest)
}

type ChiRouter struct {
	endpoints  []ChiEndpoint
	Engine     *chi.Mux
	ListenOn   string
	Config     ChiConfig
	Middleware []func(request interfaces.IRouterRequest)
}

type ChiConfig struct {
	EnableSockets bool
	ListenOn      string
}

var requests map[string]ChiRequest

func ChiSocketHandler(w http.ResponseWriter, r *http.Request) {
	var wsUpgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, err := wsUpgrade.Upgrade(w, r, nil)
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

func (broker *ChiRouter) Init(_config interface{}) error {
	config := _config.(ChiConfig)
	broker.Engine = chi.NewRouter()
	broker.Config = config
	broker.ListenOn = config.ListenOn
	if config.EnableSockets {
		//broker.Engine.PathPrefix("index.html").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(dir))))
	}
	// inject reqID generator
	broker.Engine.Use(InjectUUID)
	return nil
}

func InjectUUID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(r.Context()))
	})
}

func (broker *ChiRouter) Listen(_conn interface{}) error {
	fmt.Println("Starting Chi router ...")
	if _conn != nil {
		connString := _conn.(string)
		broker.ListenOn = connString
	}
	srv := &http.Server{
		Handler:      broker.Engine,
		Addr:         broker.ListenOn,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	err := srv.ListenAndServe()
	return err
}

func GetWrappedRequest(endpoint ChiEndpoint) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cRequest := ChiRequest{
			HttpResponse: w,
			HttpRequest:  r,
		}
		endpoint.Handler(&cRequest)
	}
}

func MiddlewareWrapper(handler func(interfaces.IRouterRequest)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cRequest := ChiRequest{
				HttpResponse: w,
				HttpRequest:  r,
				NextFunc:     next.ServeHTTP,
			}
			handler(&cRequest)
		})
	}
}

func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (broker *ChiRouter) RegisterEndpoints(endpoints interface{}) error {
	broker.endpoints = endpoints.([]ChiEndpoint)

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "static"))
	FileServer(broker.Engine, "/static", filesDir)

	if broker.Config.EnableSockets {
		broker.Engine.MethodFunc("GET", "/ws", ChiSocketHandler)
	}
	for _, endpoint := range broker.endpoints {
		middleWares := make([]func(handler http.Handler) http.Handler, 0)
		for _, middleware := range endpoint.Middleware {
			middleWares = append(middleWares, MiddlewareWrapper(middleware))
		}
		broker.Engine.With(middleWares...).MethodFunc(endpoint.Method, endpoint.URL, GetWrappedRequest(endpoint))
	}
	return nil
}

func (broker *ChiRouter) RegisterMiddleware(_middlewareHandlers interface{}) error {
	// register websocket handler conditionally
	broker.Middleware = _middlewareHandlers.([]func(request interfaces.IRouterRequest))
	for _, middleware := range broker.Middleware {
		broker.Engine.Use(MiddlewareWrapper(middleware))
	}

	return nil
}

func (broker *ChiRouter) Invoke(pattern interface{}, payload interface{}) (bool, []byte) {
	return true, []byte("chi : default response")
}

func ProvideChiRouter() *ChiRouter {
	return &ChiRouter{}
}
