package win32

// vc++ int = int32

type WNDCLASSEXW struct {
	Size UINT
	Style UINT
	WndProc WNDPROC
	ClsExtra int32
	WndExtra int32
	Instance HINSTANCE
	Icon HICON
	Cursor HCURSOR
	Background HBRUSH
	MenuName LPCWSTR
	ClassName LPCWSTR
	IconSm HICON
}

type POINT struct {
	X LONG
	Y LONG
}
type MSG struct {
	Hwnd HWND
    Message UINT
    WParam WPARAM
    LParam LPARAM
    Time DWORD
    Point POINT
    Private DWORD
}

type STROKE int32

const (
	KEY_DOWN STROKE = iota
	KEY_UP
	SYSKEY_DOWN
	SYSKEY_UP
	UNKNOWN
)

/// <summary>
/// キーボードの状態の構造体
/// </summary>
// type KEY_STATE struct {
// 	Stroke STOROKE
// 	Key KEY
// 	ScanCode uint32
// 	public uint Flags;
// 	public uint Time;
// 	public System.IntPtr ExtraInfo;
// }