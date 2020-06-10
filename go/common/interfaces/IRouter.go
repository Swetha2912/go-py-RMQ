package interfaces

type IRouter interface {
	Init(interface{}) error
	RegisterEndpoints(interface{}) error
	RegisterMiddleware(interface{}) error
	Listen(interface{}) error
	Invoke(interface{}, interface{}) (bool, []byte)
}
