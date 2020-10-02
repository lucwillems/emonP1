package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"scm.t-m-m.be/emonP1/P1"
	"strings"
	"time"
)

func main() {

	message, err := ioutil.ReadFile("test/data")
	if err != nil {
		os.Exit(1)
	}
	lines := strings.Split(string(message), "\n")
	telegram := P1.ParseTelegram(lines)
	fmt.Fprintf(os.Stdout, "Device: %s\n", telegram.Device)
	fmt.Fprintf(os.Stdout, "version: %s\n", telegram.Version)
	fmt.Fprintf(os.Stdout, "Objects: %d\n", telegram.Size())
	fmt.Fprintf(os.Stdout, "timestamp: %s\n", telegram.Timestamp.Format(time.RFC3339))
	for _, k := range telegram.SortedIds() {
		o := telegram.Get(P1.OBISId(k))
		fmt.Println(o.ToString())
	}

}
