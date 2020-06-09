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

type INPUT struct {
	Type    DWORD
	padding [4]byte
	Union   [32]byte
	// {
	//   MOUSEINPUT    mi;
	//   KEYBDINPUT    ki;
	//   HARDWAREINPUT hi;
	// } DUMMYUNIONNAME;
}
type MOUSEINPUT struct {
	DX        LONG
	DY        LONG
	MouseData DWORD
	Flags     DWORD
	Time      DWORD
	ExtraInfo ULONG_PTR
}
type KEYBDINPUT struct {
	VK        WORD
	Scan      WORD
	Flags     DWORD
	Time      DWORD
	ExtraInfo ULONG_PTR
}
type HARDWAREINPUT struct {
	Msg    DWORD
	ParamL WORD
	ParamH WORD
}
