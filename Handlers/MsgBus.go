package Handlers

type MsgBus interface {
	Close()
	IsConnected() bool
	Publish(id string, value interface{}) error
}
