package main

import (
	"clipx/controllers"
	"clipx/models"
	"clipx/views"
	"log"
	"sync"
	"time"
)

// clipboard
var cb = models.NewClipboard()

// history buffer
var list = models.NewList(16)

// globak keyboard hook
var hooked = make(chan models.KeyInfo, 64)
var hookErr = make(chan error, 1)
var hook = models.NewKeyHooker(hooked)

// monitoring clipboard
var written = make(chan bool, 16)
var monitorErr = make(chan error, 1)
var monitor = models.NewMonitor(written)

var cursor = models.NewCursor(list.Size())

// controller
var ctrl = controllers.NewController(cursor, cb, list)

// view
var viewClosed = make(chan struct{})
var view = views.NewView(ctrl, list, cursor, viewClosed)

var once sync.Once

func main() {
	// show
	go func() {
		view.Show()
	}()

	// start hooking
	go func() {
		err := hook.Start()
		if err != nil {
			hookErr <- err
		}
	}()

	// start monitoring
	go func() {
		log.Println("[begin monitoring]")
		err := monitor.Monitoring()
		if err != nil {
			monitorErr <- err
		}
	}()

	// cleanup just in case
	defer func() {
		cleanup()
	}()

	// stop by signal
	interrupted := make(chan bool)
	OnInterrupted(func() {
		log.Println("[interrupted]")
		interrupted <- true
	})

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
		case <-viewClosed:
			log.Printf("main:loop viewClosed")
			cleanup()
			log.Printf("main:loop viewClosed fin")
			break loop
		case <-interrupted:
			log.Printf("main:loop interrupted")
			cleanup()
			log.Printf("main:loop interrupted fin")
			break loop
		case err := <-monitorErr:
			log.Printf("Monitoring Error: %v\n", err)
		case err := <-hookErr:
			log.Printf("Hooker Error: %v\n", err)
		}
	} // end main loop

	// wait channel closed
	<-written
	<-hooked

	// TODO
	// save & load file
	// help
	// max lines?
	// なんかクリップボードのchan詰まっている気がする
	// view のループとmain loop 同期しないからclipboard の更新が終わった保証ないのでは？
	list.Dump()
	log.Println("[process finished]")
}

func cleanup() {
	once.Do(func() {
		err := monitor.Stop()
		if err != nil {
			log.Printf("Monitor.Stop failed: %v\n", err)
		}
		err = hook.Stop()
		if err != nil {
			log.Printf("Hooker.Stop failed: %v\n", err)
		}
	})
}

func onClipboardWritten() {
	log.Printf("[written]")
	stringable, err := cb.IsStringable()
	if err != nil {
		log.Println(err)
		return
	}
	if !stringable {
		return
	}
	str, err := cb.GetAsString()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(str)
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
		log.Printf("to be selection mode")
		ctrl.Appear()
	}
	lastKeyDown = now
}
