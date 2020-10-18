package Handlers

import "emonP1/P1"

type MsgBus interface {
	Close()
	IsConnected() bool
	Publish(value *P1.Telegram) error
}
