package Handlers

import "emonP1/P1"

type MsgBus interface {
	Debug(on bool)
	Close()
	IsConnected() bool
	Publish(cnt uint64, value *P1.Telegram) error
}
