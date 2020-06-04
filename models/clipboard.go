package models

import (
	"fmt"
	"clipx/win32"
	"unsafe"
	"syscall"
)

type Clipboard interface {
	IsStringable() (bool, error)
	GetAsString() (string, error)
}

type WindowsClipboard struct {

}

func NewClipboard() Clipboard {
	return &WindowsClipboard{}
}

const EOL = "\r\n"

var clipboardFormats = [2]win32.UINT{win32.CF_HDROP, win32.CF_UNICODETEXT}
// 文字列っぽいか否か
func (this *WindowsClipboard) IsStringable() (bool, error) {

	for i := range clipboardFormats {
		res, lastErr, err := win32.IsClipboardFormatAvailable.Call(uintptr(clipboardFormats[i]));
		if lastErr != 0 {
			return false, err
		}
		if res != win32.FALSE {
			return true, nil
		}
	}
	return false, nil
}

// 可能なら文字列として取得する
func (this *WindowsClipboard) GetAsString() (string, error) {
	switch (getFormat()) {
	case win32.CF_UNICODETEXT:
		return getStringData()
	case win32.CF_HDROP:
		return getPathData()
	default:
		return "", fmt.Errorf("unavailable clipboard format")
	}
}


func GetClipboardSequenceNumber() (uint32, error) {
	seq, lastErr, err := win32.GetClipboardSequenceNumber.Call()
	if lastErr != 0 {
		return 0, err
	}
	return uint32(seq), nil
}


func getStringData() (string, error) {
	// open
	res, lastErr, err := win32.OpenClipboard.Call(0)
	if lastErr != 0 {
		return "", err
	}
	if res == win32.FALSE {
		return "", err
	}
	defer win32.CloseClipboard.Call()

	// get handle
	resultHandle, lastErr, err := win32.GetClipboardData.Call(win32.CF_UNICODETEXT)
	if lastErr != 0 {
		return "", err
	}

	// lock
	resultPtr, lastErr, err := win32.GlobalLock.Call(resultHandle)
	if lastErr != 0 {
		return "", err
	}
	defer win32.GlobalUnlock.Call(resultHandle)
	str := win32.Utf16PtrTostring((*uint16)(unsafe.Pointer(resultPtr)))
	return str, nil
}

func getPathData() (string, error) {
	// open
	res, lastErr, err := win32.OpenClipboard.Call(0)
	if lastErr != 0 {
		return "", err
	}
	if res == win32.FALSE {
		return "", err
	}
	defer win32.CloseClipboard.Call()

	// get handle
	handle, lastErr, err := win32.GetClipboardData.Call(win32.CF_HDROP)
	if lastErr != 0 {
		return "", err
	}
	count, lastErr, err := win32.DragQueryFileW.Call(handle, 0xFFFFFFFF, uintptr(0), 0)
	if lastErr != 0 {
		return "", err
	}
	result := ""
	var i uintptr = 0
	for ; i < count; i++ {
		size, lastErr, err := win32.DragQueryFileW.Call(handle, i, 0, 0)
		if lastErr != 0 {
			return "", err
		}
		buf := make([]uint16, size+1)
		_, lastErr, err = win32.DragQueryFileW.Call(handle, uintptr(i), uintptr(unsafe.Pointer(&buf[0])), size+1)
		if lastErr != 0 {
			return "", err
		}
		result += syscall.UTF16ToString(buf) + EOL
	}
	return result, nil
}

func getFormat() win32.UINT {
	for i := range clipboardFormats {
		res, _, _ := win32.IsClipboardFormatAvailable.Call(uintptr(clipboardFormats[i]));
		if res != win32.FALSE {
			return clipboardFormats[i]
		}
	}
	return 0
}