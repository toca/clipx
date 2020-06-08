package models

import (
	"clipx/win32"
	"log"
	"syscall"
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
	res, lastErr, err := win32.RegisterClassExW.Call(uintptr(unsafe.Pointer(&windowClass)))
	if lastErr != 0 || res == win32.FALSE {
		return err
	}

	// create window
	window, lastErr, err := win32.CreateWindowExW.Call(0, uintptr(unsafe.Pointer(className)), uintptr(unsafe.Pointer(className)), 0, 0, 0, 0, 0, uintptr(win32.HWND_MESSAGE), 0, 0, 0)
	if lastErr != 0 {
		return err
	}
	this.window = window

	// register clipboard listener
	addRes, lastErr, err := win32.AddClipboardFormatListener.Call(this.window)
	if lastErr != 0 || addRes == win32.FALSE {
		return err
	}
	defer win32.RemoveClipboardFormatListener.Call(this.window)
	defer win32.DestroyWindow.Call(this.window)

	// message loop
	msg := win32.MSG{}
	for {
		res, lastErr, err := win32.GetMessageW.Call(uintptr(unsafe.Pointer(&msg)), this.window, 0, 0)
		if res == 0 {
			break
		}
		if lastErr != 0 {
			return err
		}
		win32.TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		win32.DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
	close(this.close)
	return nil
}

func (this WindowsMonitor) windowProc(window win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	res, lastErr, err := win32.DefWindowProcW.Call(window, uintptr(message), wParam, lParam)
	if lastErr != 0 {
		log.Printf("WindowsMonitor.windowProc failed: %v\n", err)
	}

	switch message {
	case win32.WM_CLIPBOARDUPDATE:
		this.written <- true
	case win32.WM_DESTROY:
		win32.PostQuitMessage.Call(0)
	}

	return res
}

func (this *WindowsMonitor) Stop() error {
	close(this.written)
	_, lastErr, err := win32.SendMessageW.Call(uintptr(this.window), uintptr(win32.WM_DESTROY), uintptr(0), uintptr(0))
	if lastErr != 0 {
		return err
	}
	<-this.close
	return nil
}
