package win32

import (
	"syscall"
)

func loadDll(name string) *syscall.DLL {
	dll, err := syscall.LoadDLL(name)
	if err != nil {
		panic(err)
	}
	return dll
}
func findProc(dll *syscall.DLL, name string) *syscall.Proc {
	proc, err := dll.FindProc(name)
	if err != nil {
		panic(err)
	}
	return proc
}

// dll
var kernel32Dll = loadDll("kernel32.dll")
var user32Dll = loadDll("user32.dll")
var shell32Dll = loadDll("shell32.dll")

// functions
var (
	GetLastError        = findProc(kernel32Dll, "GetLastError")
	
	RegisterClassExW    = findProc(user32Dll, "RegisterClassExW")
	DefWindowProcW      = findProc(user32Dll, "DefWindowProcW")
	CreateWindowExW     = findProc(user32Dll, "CreateWindowExW")
	DestroyWindow       = findProc(user32Dll, "DestroyWindow")
	GetForegroundWindow = findProc(user32Dll, "GetForegroundWindow")
	SetForegroundWindow = findProc(user32Dll, "SetForegroundWindow")
	SetFocus            = findProc(user32Dll, "SetFocus")
	SetActiveWindow     = findProc(user32Dll, "SetActiveWindow")

	AddClipboardFormatListener    = findProc(user32Dll, "AddClipboardFormatListener")
	RemoveClipboardFormatListener = findProc(user32Dll, "RemoveClipboardFormatListener")
	GetClipboardSequenceNumber    = findProc(user32Dll, "GetClipboardSequenceNumber")
	IsClipboardFormatAvailable    = findProc(user32Dll, "IsClipboardFormatAvailable")
	OpenClipboard                 = findProc(user32Dll, "OpenClipboard")
	CloseClipboard                = findProc(user32Dll, "CloseClipboard")
	GetClipboardData              = findProc(user32Dll, "GetClipboardData")
	SetClipboardData              = findProc(user32Dll, "SetClipboardData")
	EmptyClipboard                = findProc(user32Dll, "EmptyClipboard")

	DragQueryFileW = findProc(shell32Dll, "DragQueryFileW")

	SetWindowsHookExW   = findProc(user32Dll, "SetWindowsHookExW")
	UnhookWindowsHookEx = findProc(user32Dll, "UnhookWindowsHookEx")
	CallNextHookEx      = findProc(user32Dll, "CallNextHookEx")

	GetCurrentThreadId = findProc(kernel32Dll, "GetCurrentThreadId")

	GlobalAlloc           = findProc(kernel32Dll, "GlobalAlloc")
	GlobalFree            = findProc(kernel32Dll, "GlobalFree")
	GlobalLock            = findProc(kernel32Dll, "GlobalLock")
	GlobalUnlock          = findProc(kernel32Dll, "GlobalUnlock")
	GetModuleHandleW      = findProc(kernel32Dll, "GetModuleHandleW")
	GetConsoleWindow      = findProc(kernel32Dll, "GetConsoleWindow")
	GetStdHandle          = findProc(kernel32Dll, "GetStdHandle")
	SetConsoleWindowInfo  = findProc(kernel32Dll, "SetConsoleWindowInfo")
	SetConsoleCtrlHandler = findProc(kernel32Dll, "SetConsoleCtrlHandler")

	GetMessageW      = findProc(user32Dll, "GetMessageW")
	TranslateMessage = findProc(user32Dll, "TranslateMessage")
	DispatchMessageW = findProc(user32Dll, "DispatchMessageW")
	SendMessageW     = findProc(user32Dll, "SendMessageW")
	PostQuitMessage  = findProc(user32Dll, "PostQuitMessage")

	SendInput = findProc(user32Dll, "SendInput")

	CopyMemory = findProc(kernel32Dll, "RtlCopyMemory")
)
