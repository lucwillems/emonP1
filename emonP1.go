package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"os"
	"scm.t-m-m.be/emonP1/P1"
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
	parityFlag  = flag.String("parity", "none", "Parity the use (none/odd/even/mark/space).")
)

func processSerial() error {
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
		return err
	}
	defer p.Close()
	reader := bufio.NewReader(p)
	for {
		if err := process(reader); err != nil {
			return err
		}
	}
	return nil
}

func processFile() error {

	data, err := os.Open(*fileFlag)
	if err != nil {
		return err
	}
	defer data.Close()
	reader := bufio.NewReader(data)

	for {
		if err := process(reader); err != nil {
			if err == io.EOF {
				return nil
			}
			break
		}
		time.Sleep(1)
	}
	return nil
}

func process(br *bufio.Reader) error {

	if msg, err := readFrame(br); err == nil {
		if *verboseFlag {
			fmt.Fprintf(os.Stdout, "%s", msg)
		}
		if telegram, err := P1.Parse(msg, *verboseFlag); err != nil {
			return err
		} else {
			if *verboseFlag {
				fmt.Fprintf(os.Stdout, "Device: %s\n", telegram.Device)
			}
		}
		return nil
	} else {
		return err
	}
}

func readFrame(br *bufio.Reader) (string, error) {
	for {
		if b, err := br.Peek(1); err == nil {
			if string(b) != "/" {
				fmt.Printf("Ignoring garbage character: %c\n", b)
				br.ReadByte()
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
		frame, err := br.ReadBytes('!')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		bcrc, err := br.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		msg := string(frame) + "\n" + "!" + string(bcrc)
		frameCnt++
		return msg, nil
	}
}

func main() {

	fmt.Printf("emonP1 (%s)\n", VERSION)
	flag.Parse()
	fmt.Printf("running...\n")
	var err error

	if *fileFlag != "" {
		err = processFile()
	} else {
		err = processSerial()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, "%d frames processed", frameCnt)
	os.Exit(0)
}
