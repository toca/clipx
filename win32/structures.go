package win32

// vc++ int = int32

type WNDCLASSEXW struct {
	Size       UINT
	Style      UINT
	WndProc    WNDPROC
	ClsExtra   int32
	WndExtra   int32
	Instance   HINSTANCE
	Icon       HICON
	Cursor     HCURSOR
	Background HBRUSH
	MenuName   LPCWSTR
	ClassName  LPCWSTR
	IconSm     HICON
}

type POINT struct {
	X LONG
	Y LONG
}
type MSG struct {
	Hwnd    HWND
	Message UINT
	WParam  WPARAM
	LParam  LPARAM
	Time    DWORD
	Point   POINT
	Private DWORD
}

type KBDLLHOOKSTRUCT struct {
	VkCode    DWORD
	ScanCode  DWORD
	Flags     DWORD
	Time      DWORD
	ExtraInfo uintptr
}
