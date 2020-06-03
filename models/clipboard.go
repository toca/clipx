package models

import (
	"clipx/win32"
)

type Clipboard interface {
	IsStringable() (bool, error)
}

type WindowsClipboard struct {

}
func NewClipboard() Clipboard {
	return &WindowsClipboard{}
}
func (this *WindowsClipboard) IsStringable() (bool, error){
	res, lastErr, err := win32.IsClipboardFormatAvailable.Call(win32.CF_UNICODETEXT);
	if lastErr == 0 {
		return false, err
	}
	if res != win32.FALSE {
		return true, nil
	}
	return false, nil
}

func GetClipboardSequenceNumber() (uint32, error) {
	seq, lastErr, err := win32.GetClipboardSequenceNumber.Call()
	if lastErr != 0 {
		return 0, err
	}
	return uint32(seq), nil
}

// get type?
// get text