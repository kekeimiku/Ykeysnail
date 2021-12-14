//go:build windows
// +build windows

package windows

import (
	"log"
	"syscall"
	"unsafe"
	"ykeysnail/window"

	"golang.org/x/sys/windows"
)

type Windows struct{}

var (
	user32         = windows.NewLazyDLL("user32.dll")
	GetClassNameW  = user32.NewProc("GetClassNameW")
	GetWindowTextW = user32.NewProc("GetWindowTextW")
	proc           = user32.NewProc("GetForegroundWindow")
)

func (w *Windows) Window() *window.WindowInfo {

	hwnd, _, _ := proc.Call()

	class := make([]uint16, 256)

	title := make([]uint16, 256)

	r0, _, e1 := syscall.Syscall(GetClassNameW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(&class[0])), uintptr(len(class)))
	if r0 == 0 {
		if e1 != 0 {
			log.Print(e1)
		} else {
			log.Print(syscall.EINVAL)
		}
	}

	r0, _, e1 = syscall.Syscall(GetWindowTextW.Addr(), 3, uintptr(hwnd), uintptr(unsafe.Pointer(&title[0])), uintptr(len(title)))
	if r0 == 0 {
		if e1 != 0 {
			log.Print(e1)
		} else {
			log.Print(syscall.EINVAL)
		}
	}

	return &window.WindowInfo{Title: syscall.UTF16ToString(title), Class: syscall.UTF16ToString(class)}
}
