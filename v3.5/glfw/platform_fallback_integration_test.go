//go:build linux && glfw_integration && !wayland && !x11

package glfw

import (
	"os"
	"runtime"
	"testing"
)

// Запускать в Xvfb: отсутствующий Wayland не блокирует автоматический X11 fallback.
func TestAutoPlatformFallbackWithStaleWayland(t *testing.T) {
	if os.Getenv("DISPLAY") == "" {
		t.Skip("нужен X11 display")
	}
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if !PlatformSupported(PlatformX11) || !PlatformSupported(PlatformWayland) {
		t.Fatal("для проверки нужны оба backend в одном исполняемом файле")
	}
	t.Setenv("XDG_SESSION_TYPE", "wayland")
	t.Setenv("WAYLAND_DISPLAY", "fyne-nonexistent-wayland-display")
	InitHint(PlatformHint, int(AnyPlatform))
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	if GetPlatform() != PlatformX11 {
		t.Fatal(GetPlatform())
	}
	Terminate()
	InitHint(PlatformHint, int(PlatformWayland))
	if err := Init(); err == nil {
		Terminate()
		t.Fatal("явный Wayland заменён X11")
	}
	InitHint(PlatformHint, int(AnyPlatform))
}
