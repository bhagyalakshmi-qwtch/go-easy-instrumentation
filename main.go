package main

import (
	"log"

	"github.com/bhagyalakshmi-qwtch/go-easy-instrumentation/cmd"
)

func main() {
	log.Default().SetFlags(0)
	cmd.Execute()
}
