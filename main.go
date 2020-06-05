package main

import (
	"clipx/models"
	"fmt"
	"log"
	"time"
)

// clipboard
var cb = models.NewClipboard()

// history buffer
var list = models.NewList(4)

func main() {
	log.Println("std err")
	// the view
	// var view = views.NewView(list)

	// globak keyboard hook
	hooked := make(chan models.KeyInfo, 64)
	hookErr := make(chan error, 1)
	hook := models.NewKeyHooker(hooked)
	// start hooking
	go func() {
		err := hook.Start()
		if err != nil {
			hookErr <- err
		}
	}()
	defer hook.Stop()

	// monitoring clipboard
	written := make(chan bool, 16)
	monitorErr := make(chan error, 1)
	monitor := models.NewMonitor(written)
	// start monitoring
	go func() {
		fmt.Println("[begin monitoring]")
		err := monitor.Monitoring()
		if err != nil {
			monitorErr <- err
		}
	}()
	defer monitor.Stop()

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
		// view.Close()
	})

	// view
	// go func() {
	// 	view.Show()
	// }()
	// main loop
loop:
	for {
		select {
		case _, ok := <-written:
			if ok {
				onClipboardWritten()
			} else {
				break loop
			}
		case keyInfo, ok := <-hooked:
			if ok {
				onHooked(&keyInfo)
			} else {
				break loop
			}
		case err := <-monitorErr:
			fmt.Printf("Monitoring Error: %v\n", err)
		case err := <-hookErr:
			fmt.Printf("Hooker Error: %v\n", err)
		}
	}

	// wait channel close
	<-written
	<-hooked
	// change quit -> close
	// paste
	// save & load file
	list.Dump()
	fmt.Println("[process finished]")
	// todo defer

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

const VK_CTRL = 17
const VK_LCONTROL = 162
const VK_RCONTROL = 163
const ThresholdMilli = 600

var lastKeyDown = time.Now()

func onHooked(keyInfo *models.KeyInfo) {
	if keyInfo.Action != models.KeyUp {
		return
	}
	if keyInfo.VirtualKeyCode != VK_CTRL && keyInfo.VirtualKeyCode != VK_LCONTROL && keyInfo.VirtualKeyCode != VK_RCONTROL {
		return
	}
	now := time.Now()
	if now.Sub(lastKeyDown).Milliseconds() <= ThresholdMilli {
		fmt.Printf("to be selection mode")
	}
	lastKeyDown = now
}
