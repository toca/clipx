package main

import (
	"clipx/controllers"
	"clipx/models"
	"clipx/views"
	"flag"
	"log"
	"sync"
	"time"
)

const DATA_LENGTH = 16

// clipboard
var cb models.Clipboard

// history buffer
var list models.List
var storagePath *string

// globak keyboard hook
var hooked chan models.KeyInfo
var hookErr chan error
var hook models.KeyHooker

// monitoring clipboard
var written chan bool
var monitorErr chan error
var monitor models.Monitor

// cursor
var cursor models.Cursor

// controller
var ctrl controllers.Controller

// view
var view *views.View
var viewClosed chan struct{}

func main() {
	initialize()
	// show ui
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

	// start clipboard monitoring
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
	// なんかクリップボードのchan詰まっている気がする
	// view のループとmain loop 同期しないからclipboard の更新が終わった保証ないのでは？
	// list does not collect same content
	// mouse select
	list.Dump()
	log.Println("[process finished]")
}

func initialize() {
	// clipboard
	cb = models.NewClipboard()
	// history buffer
	storagePath = flag.String("s", "", "path to save and load clipboard data")
	flag.Parse()
	if 0 < len(*storagePath) {
		storage := models.NewStorage()
		data, err := storage.Load(storagePath)
		if err != nil {
			panic(err)
		}
		list = models.NewListWithData(DATA_LENGTH, data)
	} else {
		list = models.NewList(DATA_LENGTH)
	}

	// globak keyboard hook
	hooked = make(chan models.KeyInfo, 64)
	hookErr = make(chan error, 1)
	hook = models.NewKeyHooker(hooked)

	// monitoring clipboard
	written = make(chan bool, 64)
	monitorErr = make(chan error, 1)
	monitor = models.NewMonitor(written)

	// cursor
	cursor = models.NewCursor(DATA_LENGTH)

	// controller
	ctrl = controllers.NewController(cursor, cb, list)

	// view
	view = views.NewView(ctrl, list, cursor)
	viewClosed = view.GetClosedNotify()
}

var once sync.Once // for cleanup()

func cleanup() {
	log.Println("main:enter creanup")
	once.Do(func() {
		log.Println("main: begen creanup")
		storage := models.NewStorage()
		err := storage.Save(storagePath, list.GetData())
		if err != nil {
			log.Printf("Storage.Save failed: %v", err)
		}
		err = monitor.Stop()
		if err != nil {
			log.Printf("Monitor.Stop failed: %v\n", err)
		}
		err = hook.Stop()
		if err != nil {
			log.Printf("Hooker.Stop failed: %v\n", err)
		}
		log.Println("main: end creanup")
	})
	log.Println("main:reave creanup")
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
	list.Push(str)
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
