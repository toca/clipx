package win32

import (
	// "syscall"
)
// vc++ int = int32
type (
	BOOL     = int32
	DWORD    = uint32
	LONG     = int32
	SHORT    = int16
	UINT     = uint32
	ULONG    = uint32
	USHORT   = uint16
	WORD     = uint16
	WNDPROC  = uintptr
	HWND    = uintptr
	WPARAM = uintptr
	LPARAM = uintptr
	HINSTANCE = uintptr
	HICON = uintptr
	HCURSOR = uintptr
	HBRUSH = uintptr
	LPCWSTR = *uint16
	LRESULT  = uintptr
)
const HWND_MESSAGE HWND = ^HWND(2)
const (
	FALSE = 0
	TRUE = 1
)

const (
	CF_UNICODETEXT = 13
)
// WM_XXXX window message
const (
	WM_CLIPBOARDUPDATE = 0x0000031D
	// HWND_MESSAGE       = ^uint32(2)
	WM_QUIT            = 0x00000012
)