package main

import (
	"os"
	"os/signal"
	"syscall"
)

func OnInterrupted(callback func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		callback()
	}()
}
