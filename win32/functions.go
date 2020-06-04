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
var kernel32Dll  = loadDll("kernel32.dll")
var user32Dll    = loadDll("user32.dll")
var shell32Dll   = loadDll("shell32.dll")

// functions
var (
	RegisterClassExW     = findProc(user32Dll, "RegisterClassExW")
	DefWindowProcW       = findProc(user32Dll, "DefWindowProcW")
	CreateWindowExW      = findProc(user32Dll, "CreateWindowExW")
	DestroyWindow        = findProc(user32Dll, "DestroyWindow")
	
	AddClipboardFormatListener    = findProc(user32Dll, "AddClipboardFormatListener")
	RemoveClipboardFormatListener = findProc(user32Dll, "RemoveClipboardFormatListener")
	GetClipboardSequenceNumber    = findProc(user32Dll, "GetClipboardSequenceNumber")
	IsClipboardFormatAvailable    = findProc(user32Dll, "IsClipboardFormatAvailable")
	OpenClipboard                 = findProc(user32Dll, "OpenClipboard")
	CloseClipboard                = findProc(user32Dll, "CloseClipboard")
	GetClipboardData              = findProc(user32Dll, "GetClipboardData")

	DragQueryFileW                = findProc(shell32Dll, "DragQueryFileW")

	GlobalLock    = findProc(kernel32Dll, "GlobalLock")	
	GlobalUnlock  = findProc(kernel32Dll, "GlobalUnlock")

	GetMessageW           = findProc(user32Dll, "GetMessageW")
	TranslateMessage      = findProc(user32Dll, "TranslateMessage")
	DispatchMessageW      = findProc(user32Dll, "DispatchMessageW")
	SendMessageW          = findProc(user32Dll, "SendMessageW")
)