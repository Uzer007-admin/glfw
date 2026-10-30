//go:build windows && glfw_integration

package glfw

import (
	"runtime"
	"syscall"
	"testing"
	"unsafe"
)

// Deliver a notification to GLFW's own helper on this OS thread. No display
// mode is changed and no message is broadcast to other applications.
func TestMonitorConfigurationWin32DisplayChange(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Terminate()
	user32 := syscall.NewLazyDLL("user32.dll")
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getClass := user32.NewProc("GetClassNameW")
	var helper uintptr
	enumerate := syscall.NewCallback(func(hwnd, _ uintptr) uintptr {
		var name [64]uint16
		if n, _, _ := getClass.Call(hwnd, uintptr(unsafe.Pointer(&name[0])), uintptr(len(name))); n != 0 && syscall.UTF16ToString(name[:n]) == "GLFW3 Helper" {
			helper = hwnd
		}
		return 1
	})
	thread, _, _ := kernel32.NewProc("GetCurrentThreadId").Call()
	user32.NewProc("EnumThreadWindows").Call(thread, enumerate, 0)
	if helper == 0 {
		t.Fatal("GLFW helper on the current thread was not found")
	}
	send := func() { user32.NewProc("SendMessageW").Call(helper, 0x007e, 0, 0) }
	first, second := 0, 0
	SetMonitorConfigurationCallback(func() { first++ })
	send()
	if first != 1 {
		t.Fatalf("WM_DISPLAYCHANGE notifications without hotplug: %d", first)
	}
	if previous := SetMonitorConfigurationCallback(func() { second++ }); previous == nil {
		t.Fatal("previous callback was not returned")
	}
	send()
	if first != 1 || second != 1 {
		t.Fatalf("callback replacement: first=%d second=%d", first, second)
	}
	SetMonitorConfigurationCallback(nil)
	send()
	if second != 1 {
		t.Fatal("disabled callback received WM_DISPLAYCHANGE")
	}
	SetMonitorConfigurationCallback(func() {})
	Terminate()
	if fMonitorConfigurationHolder != nil {
		t.Fatal("Terminate retained the configuration callback")
	}
}
