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

var VERSION = "1.0"
var frameCnt int64 = 0
var (
	deviceFlag  = flag.String("device", "/dev/ttyUSB0", "Serial device to read P1 data from.")
	fileFlag    = flag.String("file", "", "use file as input.")
	baudFlag    = flag.Int("baud", 115200, "Baud rate (speed) to use.")
	timeoutFlag = flag.Int("timeout", 2000, "read timeout in msec.")
	bitsFlag    = flag.Int("bits", 8, "Number of databits.")
	verboseFlag = flag.Bool("verbose", false, "verbose")
	mqttUrl     = flag.String("mqtt", "", "send over mqtt url")
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
	mqtt, err := Handlers.NewMqttBus("emonP1", "", "", *mqttUrl)
	return mqtt, err
}
func StdoutBus() (*Handlers.WriteBus, error) {
	out, err := Handlers.NewWriterBus(os.Stdout)
	return out, err
}

func MessageBus() (*Handlers.MsgBus, error) {
	var bus *Handlers.MsgBus = nil
	var err error

	if *mqttUrl != "" {
		if mqtt, err := MqttBus(); err == nil {
			x := Handlers.MsgBus(mqtt)
			return &x, nil
		} else {
			return nil, err
		}
	} else {
		if out, err := StdoutBus(); err == nil {
			x := Handlers.MsgBus(out)
			return &x, nil
		} else {
			return nil, err
		}
	}
	return bus, err

}
func main() {

	fmt.Printf("emonP1 (%s)\n", VERSION)
	flag.Parse()
	fmt.Printf("running...\n")
	var err error
	var channel *Handlers.MsgBus
	var reader *bufio.Reader

	if channel, err = MessageBus(); err == nil {
		if *fileFlag != "" {
			reader, err = openFile()
		} else {
			reader, err = openSerial()
		}
		if err != nil {
			handleFatal(err)
		}
		processor := Handlers.NewP1Processor(reader, channel)
		err = processor.Process()
		fmt.Fprintf(os.Stdout, "%d frames processed", processor.FrameCnt)
		handleFatal(err)
	}
	handleFatal(err)
}

func handleFatal(err error) {
	if err != nil {
		if err == io.EOF {
			os.Exit(0)
		}
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	os.Exit(0)
}
