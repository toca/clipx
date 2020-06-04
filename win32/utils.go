package win32

import (
	"unsafe"
	"unicode/utf16"
)
// https://golang.org/src/internal/syscall/windows/syscall_windows.go
func Utf16PtrTostring(ptr *uint16) string {
	if ptr == nil {
		return ""
	}
	max := (1 << 30) - 1
	n := 0
	end := unsafe.Pointer(ptr)
	for *(*uint16)(end) != 0 && n < max {
		end = unsafe.Pointer(uintptr(end) + unsafe.Sizeof(*ptr))
		n++
	}
	buf := (*[(1 << 30) -1]uint16)(unsafe.Pointer(ptr))[:n:n] // three index
	return string(utf16.Decode(buf))
}