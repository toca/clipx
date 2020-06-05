package models

import (
	"clipx/win32"
	"fmt"
	"syscall"
	"unsafe"
)

type KeyHooker interface {
	Start() error
	Stop() error
}

type WindowsKeyHooker struct {
	hooked   chan KeyInfo
	quit     chan bool
	neighbor win32.HHOOK
	window   win32.HWND
}

func NewKeyHooker(hooked chan KeyInfo, quit chan bool) KeyHooker {
	return &WindowsKeyHooker{hooked, quit, 0, 0}
}

// generic keyboard event info
type KeyAction uint

const KeyUp KeyAction = 1
const KeyDown KeyAction = 0

type KeyInfo struct {
	Action          KeyAction
	VirtualKeyCode  uint32
	ScanCode        uint32
	ModifierKeyFlag uint32
}

// end key info

func (this *WindowsKeyHooker) Start() error {
	//register window
	className := syscall.StringToUTF16Ptr("clipx_hook")
	windowClass := win32.WNDCLASSEXW{ClassName: className}
	windowClass.WndProc = syscall.NewCallback(this.windowProc)
	windowClass.Size = win32.UINT(unsafe.Sizeof(windowClass))
	res, lastErr, err := win32.RegisterClassExW.Call(uintptr(unsafe.Pointer(&windowClass)))
	if lastErr != 0 || res == win32.FALSE {
		return err
	}

	// create window
	hinstance, _, err := win32.GetModuleHandleW.Call(0)
	this.window, lastErr, err = win32.CreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(className)),
		0,
		win32.CW_USERDEFAULT, win32.CW_USERDEFAULT,
		win32.CW_USERDEFAULT, win32.CW_USERDEFAULT,
		uintptr(win32.HWND_MESSAGE),
		0,
		hinstance,
		0)
	if lastErr != 0 {
		return err
	}

	// set hook
	this.neighbor, lastErr, err = win32.SetWindowsHookExW.Call(win32.WH_KEYBOARD_LL, syscall.NewCallback(this.hookProc), hinstance, uintptr(0))
	if lastErr != 0 {
		return err
	}

	// message loop
	msg := win32.MSG{}
	for {
		res, lastErr, err := win32.GetMessageW.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		if res == 0 {
			break
		}
		if lastErr != 0 {
			fmt.Println(err)
		}
		win32.TranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		win32.DispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}
	return nil
}

func (this *WindowsKeyHooker) Stop() error {
	_, lastErr, err := win32.SendMessageW.Call(uintptr(this.window), uintptr(win32.WM_QUIT), uintptr(0), uintptr(0))
	if lastErr != 0 {
		return err
	}

	_, lastErr, err = win32.UnhookWindowsHookEx.Call(win32.WH_KEYBOARD_LL)
	this.quit <- true
	if lastErr != 0 {
		return err
	} else {
		return nil
	}
}

func (this *WindowsKeyHooker) hookProc(code int32, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	result, lastErr, err := win32.CallNextHookEx.Call(this.neighbor, uintptr(code), wParam, lParam)
	if lastErr != 0 {
		fmt.Println(err)
	}
	switch wParam {
	case win32.WM_KEYUP:
		kbdInfo := (*win32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		keyInfo := KeyInfo{Action: KeyUp, VirtualKeyCode: kbdInfo.VkCode, ScanCode: kbdInfo.ScanCode, ModifierKeyFlag: 0} // TODO impl modifirekey
		this.hooked <- keyInfo
	case win32.WM_KEYDOWN:
		kbdInfo := (*win32.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		keyInfo := KeyInfo{Action: KeyDown, VirtualKeyCode: kbdInfo.VkCode, ScanCode: kbdInfo.ScanCode, ModifierKeyFlag: 0} // TODO impl modifirekey
		this.hooked <- keyInfo
	}
	return result
}

func (this *WindowsKeyHooker) windowProc(window win32.HWND, message win32.UINT, wParam win32.WPARAM, lParam win32.LPARAM) win32.LRESULT {
	res, lastErr, err := win32.DefWindowProcW.Call(window, uintptr(message), wParam, lParam)
	if lastErr != 0 {
		fmt.Printf("WindowsMonitor.windowProc failed: %v\n", err)
	}
	return res
}
