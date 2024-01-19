package main

import (
	"github.com/toca/clipx/win32"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

type Signals interface {
	SetOnInterrupted(chan struct{})
	SetOnTerminated(chan struct{})
}

func SetOnInterrupted(interrupted chan struct{}) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT)
	go func() {
		<-sigChan
		interrupted <- struct{}{}
		<-interrupted
	}()
}

var sigOnce = sync.Once{}
var terminatedListeners = make([]chan struct{}, 0)

func SetOnTerminated(terminated chan struct{}) {
	sigOnce.Do(func() {
		res, _, err := win32.SetConsoleCtrlHandler.Call(syscall.NewCallback(ConsoleCtrlHandler), win32.TRUE)
		if res == win32.FALSE {
			panic(err)
		}
	})
	terminatedListeners = append(terminatedListeners, terminated)
}

func ConsoleCtrlHandler(event win32.DWORD) uintptr {
	switch event {
	case win32.CTRL_CLOSE_EVENT:
		log.Println(terminatedListeners)
		for i := range terminatedListeners {
			terminatedListeners[i] <- struct{}{}
			<-terminatedListeners[i]
		}
		return uintptr(win32.TRUE)
	default:
		return uintptr(win32.FALSE)
	}
}
