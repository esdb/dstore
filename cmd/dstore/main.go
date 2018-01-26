package main

import (
	"github.com/v2pro/plz"
	"github.com/esdb/dstore/endpoint/http"
	"os"
	"os/signal"
	"fmt"
)

func main() {
	plz.PlugAndPlay()
	http.StartHttpEndpoints()
	fmt.Println("ctrl-c to quit")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for sig := range c {
		if sig.String() == "interrupt" {
			return
		}
	}
}
