package main

import (
	"os"
	"os/signal"
	"syscall"

	"experiment/containertest/testdata"
)

func main() {
	// go run()
	// go run()
	// time.Sleep(1 * time.Minute)

	run()
}

func run() {
	downContainer, err := testdata.UpContainer(true)
	if err != nil {
		panic(err.Error())
	}

	osSig := make(chan os.Signal, 2)
	signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM)
	<-osSig
	err = downContainer()
	if err != nil {
		panic(err.Error())
	}
}
