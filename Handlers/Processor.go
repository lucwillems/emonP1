package Handlers

import (
	"bufio"
	"emonP1/P1"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type P1Processor struct {
	FrameCnt uint64
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
				if p.verbose {
					fmt.Printf("Ignoring garbage character: %c\n", b)
				}
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
	(*p.output).Debug(b)
}

func (p *P1Processor) _process() error {
	if msg, err := p.readFrame(); err != nil {
		return err
	} else {
		if p.verbose {
			fmt.Fprintf(os.Stdout, "%s", msg)
		}

		if telegram, err := P1.Parse(msg, p.verbose); err != nil {
			return err
		} else {
			if p.verbose {
				fmt.Fprintf(os.Stdout, "Device: %s\n", telegram.Device)
			}
			/* publish telegram */
			if err := (*p.output).Publish(p.FrameCnt, telegram); err != nil {
				fmt.Fprint(os.Stderr, err)
			}
			p.FrameCnt++
		}
	}
	return nil
}

func (p *P1Processor) Process() error {
	fmt.Fprintf(os.Stdout, "verbose: %v\n", p.verbose)
	var signalEvent = make(chan os.Signal, 1)
	signal.Notify(signalEvent, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	for {
		select {
		case <-signalEvent:
			//fmt.Fprintf(os.Stdout,"%n processed\n",p.FrameCnt)
			return nil
		default:
			err := p._process()
			if err != nil {
				return err
			}
		}
	}
}

func (p *P1Processor) Close() {
	(*p.output).Close()
}
