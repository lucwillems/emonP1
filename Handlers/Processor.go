package Handlers

import (
	"bufio"
	"emonP1/P1"
	"fmt"
	"io"
	"os"
	"sort"
	"time"
)

type P1Processor struct {
	FrameCnt int
	verbose  bool
	input    *bufio.Reader
	output   *MsgBus
}

func NewP1Processor(input *bufio.Reader, channel *MsgBus) *P1Processor {
	var p P1Processor
	p.output = channel
	p.input = input
	return &p
}

func (p *P1Processor) readFrame() (string, error) {
	for {
		if b, err := p.input.Peek(1); err == nil {
			if string(b) != "/" {
				fmt.Printf("Ignoring garbage character: %c\n", b)
				p.input.ReadByte()
				continue
			}
		} else {
			if err == io.EOF {
				return "", err
			}
			//wait for new character
			time.Sleep(1)
			continue
		}
		frame, err := p.input.ReadBytes('!')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		bcrc, err := p.input.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		msg := string(frame) + "\n" + "!" + string(bcrc)
		return msg, nil
	}
}

func (p *P1Processor) Debug(b bool) {
	p.verbose = b
}

func (p *P1Processor) Process() error {
	if msg, err := p.readFrame(); err == nil {
		if p.verbose {
			fmt.Fprintf(os.Stdout, "%s", msg)
		}

		if telegram, err := P1.Parse(msg, p.verbose); err != nil {
			return err
		} else {
			if p.verbose {
				fmt.Fprintf(os.Stdout, "Device: %s\n", telegram.Device)
			}
			/* publish each telegram instance */
			ids := telegram.OBISIds()
			sort.Strings(ids)
			for _, k := range ids {
				o, _ := telegram.Get(k)
				if err := (*p.output).Publish(o.Id, o.Value); err != nil {
					fmt.Fprint(os.Stderr, err)
				}
				if p.verbose {
					fmt.Println(o)
				}
			}
			p.FrameCnt++
		}

	}
	return nil
}
