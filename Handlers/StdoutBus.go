package Handlers

import (
	"emonP1/P1"
	"fmt"
	"io"
	"os"
	"sort"
)

type WriteBus struct {
	debug  bool
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

func (c *WriteBus) Debug(debug bool) {
	c.debug = debug
}

func (c *WriteBus) IsConnected() bool {
	return true
}

func (c *WriteBus) Publish(cnt uint64, telegram *P1.Telegram) error {
	if c.output != nil {
		fmt.Fprintf(os.Stdout, "%d %s\n", cnt, telegram.Device)
		ids := telegram.OBISIds()
		sort.Strings(ids)
		for _, k := range ids {
			o, _ := telegram.Get(k)
			fmt.Println(o)
		}
	}
	return nil
}
