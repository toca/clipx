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
	ULONG_PTR = uint64
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

// sys command
const (
	SC_MINIMIZE   = 0xF020
	SC_RESTORE    = 0xF120
	SC_PREVWINDOW = 0xF050
)
const HWND_MESSAGE HWND = ^HWND(2) // -3

// WM_XXXX window message
const (
	WM_NULL            = 0x0000
	WM_DESTROY         = 0x0002
	WM_CLIPBOARDUPDATE = 0x0000031D
	// HWND_MESSAGE       = ^uint32(2)
	WM_QUIT        = 0x00000012
	WM_KEYDOWN     = 0x0100
	WM_KEYUP       = 0x0101
	WM_SYSKEYDOWN  = 0x0104
	WM_SYSKEYUP    = 0x0105
	WM_SYSCOMMAND  = 0x0112
	WM_IME_KEYDOWN = 0x0290
	WM_IME_KEYUP   = 0x0291
)

// KEYEVENTF
const (
	KEYEVENTF_EXTENDEDKEY = 0x0001
	KEYEVENTF_KEYUP       = 0x0002
	KEYEVENTF_UNICODE     = 0x0004
	KEYEVENTF_SCANCODE    = 0x0008
)

// INPUT
const (
	INPUT_MOUSE    = 0
	INPUT_KEYBOARD = 1
	INPUT_HARDWARE = 2
)

// GlobalAlloc
const (
	GHND = 0x0042
)

// STD handle
const (
	STD_OUTPUT_HANDLE = 0xFFFFFFF5
)

const (
	CTRL_C_EVENT        = 0
	CTRL_BREAK_EVENT    = 1
	CTRL_CLOSE_EVENT    = 2
	CTRL_LOGOFF_EVENT   = 5
	CTRL_SHUTDOWN_EVENT = 6
)
