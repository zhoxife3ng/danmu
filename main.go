package main

import (
	"github.com/x554462/danmu/danmu"
	"os"
	"os/signal"
)

var interrupt = make(chan os.Signal, 1)

func main() {
	danmu.Run()
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt
}
