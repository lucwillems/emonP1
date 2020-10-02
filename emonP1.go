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
	for _, k := range telegram.SortedIds() {
		fmt.Println(telegram.Get(gop1.OBISId(k)).ToString())
	}

}
