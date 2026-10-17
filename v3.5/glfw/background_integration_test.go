//go:build linux && glfw_integration

package glfw

import (
	"errors"
	"os"
	"runtime"
	"testing"
)

func TestBackgroundNativeLifecycle(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Terminate()
	DefaultWindowHints()
	WindowHint(Visible, False)
	WindowHint(TransparentFramebuffer, True)
	w, err := CreateWindow(96, 64, "GLFW background", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Destroy()
	s, err := w.QueryBackgroundEffects()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("actual surface: %+v", s)
	if s.Platform == BackgroundUnsupported {
		t.Fatal("native platform not detected")
	}
	if os.Getenv("FYNE_EXPECT_TRANSPARENT") == "1" && (!s.AlphaSurface || !s.TransparencyAvailable) {
		t.Fatal("alpha format/composition missing")
	}
	if os.Getenv("FYNE_EXPECT_SYSTEM_BLUR") == "1" && !s.BlurAvailable {
		t.Fatal("blur capability missing")
	}
	result := make(chan error, 1)
	go func() { _, err := w.QueryBackgroundEffects(); result <- err }()
	if err := <-result; !errors.Is(err, ErrVisualThread) {
		t.Fatalf("wrong-thread query: %v", err)
	}
	go func() { result <- w.SetBackgroundBlur(true) }()
	if err := <-result; !errors.Is(err, ErrVisualThread) {
		t.Fatalf("wrong-thread mutation: %v", err)
	}
	if s.AlphaSurface && s.Composition && s.BlurAvailable {
		for _, value := range []bool{true, false, true, false} {
			if err := w.SetBackgroundBlur(value); err != nil {
				t.Fatal(err)
			}
		}
	} else if err := w.SetBackgroundBlur(true); err == nil {
		t.Fatal("unsupported blur accepted")
	}
	if err := w.SetBackgroundBlur(false); err != nil {
		t.Fatal(err)
	}
	WindowHint(TransparentFramebuffer, False)
	opaque, err := CreateWindow(64, 48, "unreserved", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer opaque.Destroy()
	if err := opaque.SetBackgroundBlur(true); !errors.Is(err, ErrBackgroundAlpha) {
		t.Fatalf("opaque preflight: %v", err)
	}
	w.Destroy()
	if _, err := w.QueryBackgroundEffects(); !errors.Is(err, ErrWindowDestroyed) {
		t.Fatal(err)
	}
	if err := w.SetBackgroundBlur(false); !errors.Is(err, ErrWindowDestroyed) {
		t.Fatal(err)
	}
}
