package main

import (
	"clipx/models"
	"fmt"
)

// clipboard
var cb = models.NewClipboard()

// history buffer
var list = models.NewList(4)

func main() {

	// globak keyboard hook
	hookQuit := make(chan bool, 1)
	hooked := make(chan bool, 64)
	hookErr := make(chan error, 1)
	hook := models.NewKeyHooker(hooked, hookQuit)
	// start hooking
	go func() {
		err := hook.Start()
		if err != nil {
			hookErr <- err
		}
	}()

	// monitoring clipboard
	written := make(chan bool, 16)
	cbQuit := make(chan bool, 1)
	monitorErr := make(chan error, 1)
	monitor := models.NewMonitor(written, cbQuit)
	// start monitoring
	go func() {
		fmt.Println("[begin monitoring]")
		err := monitor.Monitoring()
		if err != nil {
			monitorErr <- err
		}
	}()

	// stop
	// signals
	OnInterrupted(func() {
		fmt.Println("[interrupted]")
		err := monitor.Stop()
		if err != nil {
			fmt.Printf("Monitor.Stop failed: %v\n", err)
		}
		err = hook.Stop()
		if err != nil {
			fmt.Printf("Hooker.Stop failed: %v\n", err)
		}
	})

	// main loop
loop:
	for {
		select {
		case <-written:
			onClipboardWritten()
		case <-hooked:
			onHooked()
		case <-cbQuit:
			<-hookQuit
			fmt.Println("[quit]")
			break loop
		case err := <-monitorErr:
			fmt.Printf("Monitoring Error: %v\n", err)
		case err := <-hookErr:
			fmt.Printf("Hooker Error: %v\n", err)
		}
	}

	// key hook
	// buffer
	// ui
	// paste
	list.Dump()
	fmt.Println("[process finished]")
}

func onClipboardWritten() {
	fmt.Printf("[written]")
	stringable, err := cb.IsStringable()
	if err != nil {
		fmt.Println(err)
		return
	}
	if !stringable {
		return
	}
	str, err := cb.GetAsString()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(str)
	}
	list.Add(str)
}

func onHooked() {
	fmt.Println("HOOKED!")
}
