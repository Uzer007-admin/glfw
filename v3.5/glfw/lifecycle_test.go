//go:build glfw_integration

package glfw

import (
	"runtime"
	"testing"
)

func TestTerminateWindowRegistry(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for iteration := 0; iteration < 2; iteration++ {
		if err := Init(); err != nil {
			t.Fatal(err)
		}
		WindowHint(Visible, False)
		if _, err := CreateWindow(8, 8, "registry", nil, nil); err != nil {
			Terminate()
			t.Fatal(err)
		}
		Terminate()
		windows.l.Lock()
		count := len(windows.m)
		windows.l.Unlock()
		if count != 0 {
			t.Fatalf("iteration %d: %d retained windows", iteration, count)
		}
	}
}

func TestUnknownScancodeName(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Terminate()
	if name := GetKeyName(KeyUnknown, -1); name != "" {
		t.Fatalf("invalid scancode: %q", name)
	}
}

func TestNativeVersion(t *testing.T) {
	major, minor, revision := GetVersion()
	if major != 3 || minor != 5 || revision != 1 {
		t.Fatalf("%d.%d.%d", major, minor, revision)
	}
}
