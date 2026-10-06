//go:build windows

package admin

import (
	"syscall"
	"unsafe"
)

func platformLocale() string {
	proc := syscall.NewLazyDLL("kernel32.dll").NewProc("GetUserDefaultLocaleName")
	buffer := make([]uint16, 85)
	result, _, _ := proc.Call(uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)))
	if result == 0 {
		return ""
	}
	return syscall.UTF16ToString(buffer)
}
