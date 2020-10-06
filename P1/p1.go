package P1

import (
	"bufio"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// P1 allows you to easily read from a P1-compatible serial device. The output is
// parsed into structured data
type P1 struct {
	serialDevice *serial.Port
	stopped      chan (interface{})
	Incoming     chan (*Telegram)
}

// P1Config is the configuration to create a new P1 object with
type P1Config struct {
	USBDevice string
	Baudrate  int
	Timeout   int // in milliseconds
}

// New returns a P1 object with given configuration or error when something went
// wrong initialising the serial object
func New(config P1Config) (*P1, error) {
	if config.Baudrate <= 0 {
		config.Baudrate = 115200
	}

	if config.Timeout <= 0 {
		config.Timeout = 500
	}

	serialConfig := &serial.Config{
		Name:        config.USBDevice,
		Baud:        config.Baudrate,
		ReadTimeout: time.Millisecond * time.Duration(config.Timeout),
	}

	serialDevice, err := serial.OpenPort(serialConfig)
	if err != nil {
		return nil, err
	}

	return &P1{
		serialDevice: serialDevice,
		Incoming:     make(chan *Telegram),
	}, nil
}

// Start makes P1 start reading data from the serial device
func (p *P1) Start() {
	go p.readData()
}

func (p *P1) readData() {
	// open reader to serial device
	reader := bufio.NewReader(p.serialDevice)

	for {
		// telegram data is suffixed with a CRC code
		// this CRC code starts with ! so let's read until we receive that char
		message, err := reader.ReadString('\x21') // hex char code for !
		if err != nil {
			continue
		}

		/*
			// we're currently not processing the CRC itself, but let's read that as well
			// from the CRC, we assume the first carriage return is the end of the CRC
			_, err = reader.ReadString('\x0d') // hex char for CR
			if err != nil {
				continue
			}
		*/

		lines := strings.Split(message, "\n")
		p.Incoming <- ParseTelegram(lines)
	}
}
