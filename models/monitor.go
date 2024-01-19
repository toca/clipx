package models

import (
	"github.com/toca/clipx/win32"
	"log"
	"syscall"
	"time"
	"unsafe"
)

type Monitor interface {
	Monitoring() error
	Stop() error
}

func NewMonitor(written chan bool) Monitor {
	return &WindowsMonitor{written, 0, make(chan bool)}
}

type WindowsMonitor struct {
	written chan bool
	window  win32.HWND
	close   chan bool
}

func (this *WindowsMonitor) Monitoring() error {
	//register window
	className := syscall.StringToUTF16Ptr("clipx")
	windowClass := win32.WNDCLASSEXW{ClassName: className}
	windowClass.WndProc = syscall.NewCallback(this.windowProc)
	windowClass.Size = win32.UINT(unsafe.Sizeof(windowClass))
	res, _, err := win32.RegisterClassExW.Call(uintptr(unsafe.Pointer(&windowClass)))
	if res == 0 {
		log.Printf("WindowsMonitor.Monitoring RegisterClassEx Error: %v", err)
		return err
	}

	// create window
	window, _, err := win32.CreateWindowExW.Call(0, uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(className)), 0, 0, 0, 0, 0, uintptr(win32.HWND_MESSAGE), 0, 0, 0)
	if window == 0 {
		log.Printf("WindowsMonitor.Monitoring CreateWindowEx Error: %v", err)
		return err
	}
	this.window = window

	// register clipboard listener
	addRes, _, err := win32.AddClipboardFormatListener.Call(this.window)
	if addRes == win32.FALSE {
		log.Printf("WindowsMonitor.Monitoring AddClipboardFormatListener Error: %v", err)
		return err
	}
	defer win32.RemoveClipboardFormatListener.Call(this.window)
	defer win32.DestroyWindow.Call(this.window)

	// message loop
	msg := win32.MSG{}
	for {
		res, _, _ := win32.GetMessageW.Call(uintptr(unsafe.Pointer(&msg)), this.window, 0, 0)
		if res == 0 {
			log.Printf("WindowsMonitor GetMeessage end")
			break
		}
		win32.TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		win32.DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
	close(this.close)
	return nil
}

func (this WindowsMonitor) windowProc(window win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	res, _, _ := win32.DefWindowProcW.Call(window, uintptr(message), wParam, lParam)

	switch message {
	case win32.WM_CLIPBOARDUPDATE:
		this.written <- true
		time.Sleep(5 * time.Millisecond)
	case win32.WM_DESTROY:
		win32.PostQuitMessage.Call(0)
	}

	return res
}

func (this *WindowsMonitor) Stop() error {
	log.Printf("Monitor.Stop")
	close(this.written)
	lresult, _, _ := win32.SendMessageW.Call(uintptr(this.window), uintptr(win32.WM_DESTROY), uintptr(0), uintptr(0))
	log.Printf("Monitor.Stop lresult:%v", lresult)
	<-this.close
	return nil
}
