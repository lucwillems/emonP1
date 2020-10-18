package main

import (
	"bufio"
	"emonP1/Handlers"
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"os"
	"time"
)

var terminate = make(chan bool, 1)
var VERSION = "1.0"
var frameCnt int64 = 0
var (
	deviceFlag  = flag.String("device", "/dev/ttyUSB0", "Serial device to read P1 data from.")
	fileFlag    = flag.String("file", "", "use file as input.")
	baudFlag    = flag.Int("baud", 115200, "Baud rate (speed) to use.")
	timeoutFlag = flag.Int("timeout", 2000, "read timeout in msec.")
	bitsFlag    = flag.Int("bits", 8, "Number of databits.")
	verboseFlag = flag.Bool("verbose", false, "verbose")
	mqttUrl     = flag.String("mqtt", "tcp://127.0.0.1:1883", "send over mqtt url")
	parityFlag  = flag.String("parity", "none", "Parity the use (none/odd/even/mark/space).")
)

func openSerial() (*bufio.Reader, error) {
	var parity serial.Parity
	switch *parityFlag {
	case "none":
		parity = serial.ParityNone
	case "odd":
		parity = serial.ParityOdd
	case "even":
		parity = serial.ParityEven
	case "mark":
		parity = serial.ParityMark
	case "space":
		parity = serial.ParitySpace
	default:
		log.Fatal("Invalid parity setting")
	}

	c := &serial.Config{
		Name:        *deviceFlag,
		Baud:        *baudFlag,
		Size:        byte(*bitsFlag),
		Parity:      parity,
		ReadTimeout: time.Millisecond * time.Duration(*timeoutFlag),
	}
	fmt.Fprintf(os.Stdout, "open %s (%v,%v,%s)\n", c.Name, c.Baud, c.Size, *parityFlag)
	p, err := serial.OpenPort(c)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(p)
	return reader, nil
}

func openFile() (*bufio.Reader, error) {
	fmt.Fprintf(os.Stdout, "opening file %s\n", *fileFlag)
	data, err := os.Open(*fileFlag)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(data)
	return reader, nil
}

func MqttBus() (*Handlers.MqttBus, error) {
	mqtt, err := Handlers.NewMqttBus("emonP1", "emonP1", "emonP1", *mqttUrl)
	return mqtt, err
}
func StdoutBus() (*Handlers.WriteBus, error) {
	out, err := Handlers.NewWriterBus(os.Stdout)
	return out, err
}

func Reader() *bufio.Reader {
	var reader *bufio.Reader
	var err error
	if *fileFlag != "" {
		if reader, err = openFile(); err != nil {
			handleFatal(err)
		}
	} else {
		if reader, err = openSerial(); err != nil {
			handleFatal(err)
		}
	}
	return reader
}

func MessageBus() *Handlers.MsgBus {
	var msgBus Handlers.MsgBus
	if *mqttUrl != "" {
		if mqtt, err := MqttBus(); err != nil {
			handleFatal(err)
		} else {
			msgBus = Handlers.MsgBus(mqtt)
		}
	} else {
		if out, err := StdoutBus(); err != nil {
			handleFatal(err)
		} else {
			msgBus = Handlers.MsgBus(out)
		}
	}
	return &msgBus
}

func process(processor *Handlers.P1Processor) {
	if err := processor.Process(); err != nil {
		handleFatal(err)
	}
	processor.Close()
	terminate <- true
}

func main() {
	fmt.Printf("emonP1 (%s)\n", VERSION)
	flag.Parse()
	fmt.Printf("running...\n")
	processor := Handlers.NewP1Processor(Reader(), MessageBus())
	processor.Debug(*verboseFlag)
	go process(processor)
	<-terminate
	fmt.Fprintf(os.Stdout, "%d frames processed\n", processor.FrameCnt)
	fmt.Printf("done")
}

func handleFatal(err error) {
	if err != nil {
		if err == io.EOF {
			return
		}
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
}
