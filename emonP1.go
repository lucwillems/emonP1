package main

import (
	"fmt"
	"github.com/skoef/gop1"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func main() {

	message, err := ioutil.ReadFile("test/data")
	if err != nil {
		os.Exit(1)
	}
	lines := strings.Split(string(message), "\n")
	telegram := gop1.ParseTelegram(lines)
	fmt.Fprintf(os.Stdout, "Device: %s\n", telegram.Device)
	fmt.Fprintf(os.Stdout, "version: %s\n", telegram.Version)
	fmt.Fprintf(os.Stdout, "Objects: %d\n", telegram.Size())
	fmt.Fprintf(os.Stdout, "timestamp: %s\n", telegram.Timestamp.Format(time.RFC3339))

	for _, v := range telegram.Objects {
		fmt.Println(v.ToString())
	}

	fmt.Fprintf(os.Stdout, "%s\n", telegram.Get(gop1.OBISTypeElectricityDelivered).ToString())
	if f, err := telegram.Get(gop1.OBISTypeInstantaneousVoltageL1).AsFloat(); err == nil {
		fmt.Fprintf(os.Stdout, "float: %.2f\n", f)

	} else {
		fmt.Fprint(os.Stderr, err)
	}
}
