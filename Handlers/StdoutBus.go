package Handlers

import (
	"fmt"
	"io"
)

type WriteBus struct {
	output io.Writer
}

func NewWriterBus(writer io.Writer) (*WriteBus, error) {
	var bus WriteBus
	return &bus, nil
}

/* implement MsgBus interface for WriterBus */

func (c *WriteBus) Close() {
}

func (c *WriteBus) IsConnected() bool {
	return true
}

func (c *WriteBus) Publish(id string, value interface{}) error {
	if c.output != nil {
		fmt.Fprintf(c.output, "%s : %s", id, fmt.Sprint(value))
	}
	return nil
}
