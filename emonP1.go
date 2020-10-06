package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"scm.t-m-m.be/emonP1/P1"
	"sort"
	"time"
)

func main() {

	message, err := ioutil.ReadFile("test/data")
	if err != nil {
		os.Exit(1)
	}
	telegram := P1.Parse(string(message))
	fmt.Fprintf(os.Stdout, "Device: %s\n", telegram.Device)
	fmt.Fprintf(os.Stdout, "version: %s\n", telegram.Version)
	fmt.Fprintf(os.Stdout, "Failures: %d\n", telegram.Failures)
	fmt.Fprintf(os.Stdout, "Objects: %d\n", telegram.Size())
	fmt.Fprintf(os.Stdout, "timestamp: %s\n", telegram.Timestamp.Format(time.RFC3339))
	ids := telegram.OBISIds()
	sort.Strings(ids)
	for _, k := range ids {
		o, _ := telegram.Get(k)
		fmt.Println(o)
	}

}
