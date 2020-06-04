package win32

import (
	// "syscall"
)
// vc++ int = int32
type (
	BOOL     = int32
	DWORD    = uint32
	HBRUSH = uintptr
	HCURSOR = uintptr
	HDROP = uintptr
	HICON = uintptr
	HINSTANCE = uintptr
	HWND    = uintptr
	LONG     = int32
	LPARAM = uintptr
	LPCWSTR = *uint16
	LRESULT  = uintptr
	SHORT    = int16
	UINT     = uint32
	ULONG    = uint32
	USHORT   = uint16
	WNDPROC  = uintptr
	WORD     = uint16
	WPARAM = uintptr
)

const (
	FALSE = 0
	TRUE = 1
)

// clipboard format
const (
	CF_UNICODETEXT = 13
	CF_HDROP       = 15
)

const HWND_MESSAGE HWND = ^HWND(2) // -3

// WM_XXXX window message
const (
	WM_CLIPBOARDUPDATE = 0x0000031D
	// HWND_MESSAGE       = ^uint32(2)
	WM_QUIT            = 0x00000012
)