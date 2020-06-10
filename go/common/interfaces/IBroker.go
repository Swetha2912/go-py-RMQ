package interfaces

type IBroker interface {
	Connect(interface{}) error
	GetData(interface{}) (error, []byte)
	PutData(interface{}, []byte) (error)
}
