//go:build windows && glfw_integration

package glfw

import (
	"runtime"
	"syscall"
	"testing"
	"unsafe"
)

func TestCompositionSurfaceCreationHint(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Terminate()
	user := syscall.NewLazyDLL("user32.dll")
	get := "GetWindowLongPtrW"
	if unsafe.Sizeof(uintptr(0)) == 4 {
		get = "GetWindowLongW"
	}
	for _, scenario := range []struct{ hint, want int }{{True, 0x00200000}, {False, 0}, {True, 0}} {
		DefaultWindowHints()
		WindowHint(Visible, False)
		WindowHint(Win32NoRedirectionBitmap, scenario.hint)
		if scenario.hint == True && scenario.want == 0 {
			DefaultWindowHints() // Сброс действительно возвращает обычный HWND.
			WindowHint(Visible, False)
		}
		w, err := CreateWindow(96, 64, "Проверка поверхности DirectComposition", nil, nil)
		if err != nil {
			t.Fatal(err)
		}
		style, _, _ := user.NewProc(get).Call(uintptr(unsafe.Pointer(w.GetWin32Window())), ^uintptr(19))
		w.MakeContextCurrent()
		if GetCurrentContext() != w {
			t.Fatal("OpenGL context was not made current")
		}
		DetachCurrentContext()
		w.Destroy()
		if int(style)&0x00200000 != scenario.want {
			t.Fatalf("native style 0x%x, expected flag 0x%x", style, scenario.want)
		}
	}
}
