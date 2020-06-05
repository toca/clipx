package win32

// "syscall"

// vc++ int = int32
type (
	BOOL      = int32
	DWORD     = uint32
	HBRUSH    = uintptr
	HCURSOR   = uintptr
	HDROP     = uintptr
	HICON     = uintptr
	HINSTANCE = uintptr
	HHOOK     = uintptr
	HWND      = uintptr
	LONG      = int32
	LPARAM    = uintptr
	LPCWSTR   = *uint16
	LRESULT   = uintptr
	SHORT     = int16
	UINT      = uint32
	ULONG     = uint32
	USHORT    = uint16
	WNDPROC   = uintptr
	WORD      = uint16
	WPARAM    = uintptr
)

const (
	FALSE = 0
	TRUE  = 1
)

// create window
const CW_USERDEFAULT = 0x80000000

// clipboard format
const (
	CF_UNICODETEXT = 13
	CF_HDROP       = 15
)

// windows hook
const (
	WH_KEYBOARD_LL = 13
)

const HWND_MESSAGE HWND = ^HWND(2) // -3

// WM_XXXX window message
const (
	WM_CLIPBOARDUPDATE = 0x0000031D
	// HWND_MESSAGE       = ^uint32(2)
	WM_QUIT       = 0x00000012
	WM_KEYDOWN    = 0x0100
	WM_KEYUP      = 0x0101
	WM_SYSKEYDOWN = 0x0104
	WM_SYSKEYUP   = 0x0105
)
