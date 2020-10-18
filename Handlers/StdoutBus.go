package Handlers

import (
	"emonP1/P1"
	"fmt"
	"io"
	"os"
)

type WriteBus struct {
	output io.Writer
}

func NewWriterBus(writer io.Writer) (*WriteBus, error) {
	var bus WriteBus
	bus.output = writer
	return &bus, nil
}

/* implement MsgBus interface for WriterBus */

func (c *WriteBus) Close() {
}

func (c *WriteBus) IsConnected() bool {
	return true
}

func (c *WriteBus) Publish(telegram *P1.Telegram) error {
	if c.output != nil {
		fmt.Fprintf(os.Stdout, "%s : %s\n", telegram.Device, telegram.Timestamp)
	}
	return nil
}
