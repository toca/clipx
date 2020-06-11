package models

import (
	"clipx/win32"
	"log"
	"unsafe"
)

type Window interface {
	Show()
	Hide()
	SendPasteCommand() error
	ResizeWindow(w int16, i int16) error
}

func NewWindow() Window {
	windowHandle, lastErr, err := win32.GetConsoleWindow.Call()
	if lastErr != 0 {
		log.Panic(err)
	}
	return &MsWinWindow{windowHandle}
}

type MsWinWindow struct {
	windowHandle win32.HWND
}

func (this *MsWinWindow) Show() {
	this.Hide()
	res, _, err := win32.SendMessageW.Call(this.windowHandle, win32.WM_SYSCOMMAND, win32.SC_RESTORE, 0)
	if res != 0 {
		log.Println(err)
	}
	// does not work
	// res, lastErr, err := win32.SetActiveWindow.Call(this.windowHandle)
	// log.Printf("@MsWinWindow: %v:%v:%v", res, lastErr, err)
	// if lastErr != 0 {
	// 	log.Println(err)
	// }
}

func (this *MsWinWindow) Hide() {
	_, lastErr, err := win32.SendMessageW.Call(this.windowHandle, win32.WM_SYSCOMMAND, win32.SC_MINIMIZE, 0)
	if lastErr != 0 {
		log.Println(err)
	}
	_, lastErr, err = win32.SendMessageW.Call(this.windowHandle, win32.WM_SYSCOMMAND, win32.SC_PREVWINDOW, 0)
	if lastErr != 0 {
		log.Println(err)
	}
}

const VK_CONTROL = 17
const VK_V = 86

func (this *MsWinWindow) SendPasteCommand() error {
	const numInput = 4
	var keyInputs [4]win32.INPUT
	ctrl := win32.KEYBDINPUT{
		VK:        VK_CONTROL,
		Scan:      0,
		Flags:     win32.KEYEVENTF_EXTENDEDKEY,
		Time:      0,
		ExtraInfo: 0,
	}
	v := win32.KEYBDINPUT{
		VK:        VK_V,
		Scan:      0,
		Flags:     win32.KEYEVENTF_EXTENDEDKEY,
		Time:      0,
		ExtraInfo: 0,
	}
	ctrlUp := win32.KEYBDINPUT{
		VK:        VK_CONTROL,
		Scan:      0,
		Flags:     win32.KEYEVENTF_EXTENDEDKEY | win32.KEYEVENTF_KEYUP,
		Time:      0,
		ExtraInfo: 0,
	}
	vUp := win32.KEYBDINPUT{
		VK:        VK_V,
		Scan:      0,
		Flags:     win32.KEYEVENTF_EXTENDEDKEY | win32.KEYEVENTF_KEYUP,
		Time:      0,
		ExtraInfo: 0,
	}

	keyInputs[0].Type = win32.INPUT_KEYBOARD
	p1 := (*[32]byte)(unsafe.Pointer(&ctrl))
	keyInputs[0].Union = *p1
	keyInputs[1].Type = win32.INPUT_KEYBOARD
	p2 := (*[32]byte)(unsafe.Pointer(&v))
	keyInputs[1].Union = *p2
	keyInputs[2].Type = win32.INPUT_KEYBOARD
	p3 := (*[32]byte)(unsafe.Pointer(&ctrlUp))
	keyInputs[2].Union = *p3
	keyInputs[3].Type = win32.INPUT_KEYBOARD
	p4 := (*[32]byte)(unsafe.Pointer(&vUp))
	keyInputs[3].Union = *p4

	res, _, err := win32.SendInput.Call(numInput, uintptr(unsafe.Pointer(&keyInputs[0])), unsafe.Sizeof(keyInputs[0]))
	if res != numInput {
		return err
	} else {
		return nil
	}
}

func (this MsWinWindow) ResizeWindow(w int16, h int16) error {
	stdOutHandle, lastErr, err := win32.GetStdHandle.Call(win32.STD_OUTPUT_HANDLE)
	if lastErr != 0 {
		return err
	}
	newSize := win32.SMALL_RECT{0, 0, w, h}
	_, lastErr, err = win32.SetConsoleWindowInfo.Call(stdOutHandle, win32.TRUE, uintptr(unsafe.Pointer(&newSize)))
	if lastErr != 0 {
		return err
	}
	return nil
}
