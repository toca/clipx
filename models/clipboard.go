package models

import (
	"github.com/toca/clipx/win32"
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

type Clipboard interface {
	IsStringable() (bool, error)
	GetAsString() (string, error)
	SetString(data *string) error
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
		res, _, _ := win32.IsClipboardFormatAvailable.Call(uintptr(clipboardFormats[i]))
		if res != win32.FALSE {
			return true, nil
		}
	}
	return false, nil
}

// 可能なら文字列として取得する
func (this *WindowsClipboard) GetAsString() (string, error) {
	switch getFormat() {
	case win32.CF_UNICODETEXT:
		return getStringData()
	case win32.CF_HDROP:
		return getPathData()
	default:
		return "", fmt.Errorf("unavailable clipboard format")
	}
}

// クリップボードにデータをセット
func (this *WindowsClipboard) SetString(rawData *string) error {
	data := syscall.StringToUTF16(*rawData)
	size := uintptr(len(data) * 2)
	// GlobalAlloc
	globalHandle, _, err := win32.GlobalAlloc.Call(win32.GHND, size)
	if globalHandle == 0 {
		log.Printf("Clipboard.SetString Error: %v", globalHandle)
		return err
	}
	// GlobalLock get pointer
	blockHandle, _, err := win32.GlobalLock.Call(globalHandle)
	if blockHandle == 0 {
		log.Printf("Clipboard.SetString Error: %v", blockHandle)
		win32.GlobalFree.Call(globalHandle)
		return err
	}
	// GlobalUnlock
	defer func() {
		res, _, err := win32.GlobalUnlock.Call(blockHandle)
		if res == 0 {
			lastError, _, _ := win32.GetLastError.Call()
			if lastError != 0 {
				log.Println(err)
			}
		}
	}()
	// can not detect error?
	_, _, _ = win32.CopyMemory.Call(blockHandle, uintptr(unsafe.Pointer(&data[0])), size)

	// OpenClipboard
	res, _, _ := win32.OpenClipboard.Call(0)
	// CloseClipboard
	defer func() {
		res, _, err := win32.CloseClipboard.Call(blockHandle)
		if res == win32.FALSE {
			log.Println(err)
		}
	}()
	if res == win32.FALSE {
		log.Printf("Clipboard.SetString Error: %v", err)
		return err
	}
	// EmptyClipboard
	res, _, err = win32.EmptyClipboard.Call()
	if res == win32.FALSE {
		log.Printf("Clipboard.SetString Error: %v", err)
		return err
	}
	// SetClipboardData
	res, _, err = win32.SetClipboardData.Call(win32.CF_UNICODETEXT, blockHandle)
	if res == 0 {
		log.Printf("Clipboard.SetString Error: %v", err)
		return err
	}
	return nil
}

func GetClipboardSequenceNumber() (uint32, error) {
	seq, _, err := win32.GetClipboardSequenceNumber.Call()
	if seq == 0 {
		return 0, err
	}
	return uint32(seq), nil
}

func getStringData() (string, error) {
	// open
	res, _, err := win32.OpenClipboard.Call(0)
	if res == win32.FALSE {
		lastErr, _, _ := win32.GetLastError.Call()
		log.Printf("Clipboard getStringData LastError:%v", lastErr)
		return "", err
	}
	defer win32.CloseClipboard.Call()

	// get handle
	resultHandle, _, err := win32.GetClipboardData.Call(win32.CF_UNICODETEXT)
	if resultHandle == 0 {
		return "", err
	}

	// lock
	resultPtr, _, err := win32.GlobalLock.Call(resultHandle)
	if resultPtr == 0 {
		return "", err
	}
	defer win32.GlobalUnlock.Call(resultHandle)
	str := win32.Utf16PtrTostring((*uint16)(unsafe.Pointer(resultPtr)))
	return str, nil
}

func getPathData() (string, error) {
	// open
	res, _, err := win32.OpenClipboard.Call(0)
	if res == win32.FALSE {
		lastErr, _, _ := win32.GetLastError.Call()
		log.Printf("Clipboard getPath LastError: %v", lastErr)
		return "", err
	}
	defer win32.CloseClipboard.Call()

	// get handle
	handle, _, err := win32.GetClipboardData.Call(win32.CF_HDROP)
	if handle == 0 {
		return "", err
	}
	count, _, err := win32.DragQueryFileW.Call(handle, 0xFFFFFFFF, uintptr(0), 0)
	if count == 0 {
		return "", err
	}
	result := ""
	var i uintptr = 0
	for ; i < count; i++ {
		size, _, err := win32.DragQueryFileW.Call(handle, i, 0, 0)
		if size != 0 {
			return "", err
		}
		buf := make([]uint16, size+1)
		res, _, err = win32.DragQueryFileW.Call(handle, uintptr(i), uintptr(unsafe.Pointer(&buf[0])), size+1)
		if res == 0 {
			return "", err
		}
		result += syscall.UTF16ToString(buf) + EOL
	}
	return result, nil
}

func getFormat() win32.UINT {
	for i := range clipboardFormats {
		res, _, _ := win32.IsClipboardFormatAvailable.Call(uintptr(clipboardFormats[i]))
		if res != win32.FALSE {
			return clipboardFormats[i]
		}
	}
	return 0
}
