/// +build zmq

package infrastructures

import (
	"fmt"
	"github.com/zeromq/goczmq"
)

type ZmqBroker struct {
	Socket *goczmq.Sock
	Config ZMQConfig
}

type ZMQConfig struct {
	SocketType       string
	ConnectionString string
}

func (broker *ZmqBroker) Connect(_config interface{}) error {
	broker.Config = _config.(ZMQConfig)
	if broker.Config.SocketType == "PULL" {
		broker.Socket, _ = goczmq.NewPush(broker.Config.ConnectionString)
	} else if broker.Config.SocketType == "PUSH" {
		broker.Socket, _ = goczmq.NewPush(broker.Config.ConnectionString)
	} else if broker.Config.SocketType == "REQ" {
		broker.Socket, _ = goczmq.NewReq(broker.Config.ConnectionString)
	} else if broker.Config.SocketType == "REP" {
		broker.Socket, _ = goczmq.NewRep(broker.Config.ConnectionString)
	}
	return nil
}

func (broker *ZmqBroker) GetData(_config interface{}) (error, []byte) {
	bytes, _, _ := broker.Socket.RecvFrame()
	fmt.Println("$$$$$$$4",bytes)
	return nil, bytes
}

func (broker *ZmqBroker) PutData(_config interface{}, data []byte) error {
	err := broker.Socket.SendFrame(data, goczmq.FlagNone)
	return err
}

func ProvideZmqBroker() *ZmqBroker {
	return &ZmqBroker{}
}
