//go:build glfw_integration

package glfw

import (
	"errors"
	"runtime"
	"testing"
)

func TestWindowVisualNativeLifecycle(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Terminate()
	DefaultWindowHints()
	WindowHint(Visible, False)
	WindowHint(TransparentFramebuffer, True)
	w, err := CreateWindow(96, 64, "native visuals test", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if w.data != nil {
			w.Destroy()
		}
	}()
	state, err := w.QueryWindowVisuals()
	if runtime.GOOS != "windows" {
		if !errors.Is(err, ErrVisualUnsupported) {
			t.Fatal(err)
		}
	} else {
		if err != nil || state.OSBuild == 0 || state.HighContrastKnown == false {
			t.Fatalf("state %+v, %v", state, err)
		}
		t.Logf("created surface: %+v", state)
		if state.OSAttributes&(VisualDarkCaption|VisualCorners) == VisualDarkCaption|VisualCorners {
			result, err := w.ApplyWindowVisuals(WindowVisualOptions{Attributes: VisualDarkCaption | VisualCorners, DarkCaption: true, Corners: CornersRound})
			if err != nil || result.Accepted != VisualDarkCaption|VisualCorners {
				t.Fatalf("caption/corners: %+v %v", result, err)
			}
		}
		if state.OSAttributes&VisualBackdrop != 0 {
			for _, mode := range []SystemBackdrop{BackdropMainWindow, BackdropTransientWindow, BackdropTabbedWindow, BackdropNone} {
				result, err := w.ApplyWindowVisuals(WindowVisualOptions{Attributes: VisualBackdrop, Backdrop: mode})
				if err != nil || result.Accepted != VisualBackdrop {
					t.Fatalf("mode %v: %+v %v", mode, result, err)
				}
				current, err := w.QueryWindowVisuals()
				if err != nil || !current.BackdropKnown || current.Backdrop != mode {
					t.Fatalf("backdrop query: %+v %v", current, err)
				}
			}
		}
		if state.OSAttributes&VisualPremultipliedAlpha != 0 && state.WGLAlphaAvailable() {
			for _, enabled := range []bool{true, false} {
				result, err := w.ApplyWindowVisuals(WindowVisualOptions{Attributes: VisualPremultipliedAlpha, PremultipliedAlpha: enabled})
				if err != nil || result.Accepted != VisualPremultipliedAlpha {
					t.Fatalf("alpha: %+v %v", result, err)
				}
			}
		}
		if _, err := w.ApplyWindowVisuals(WindowVisualOptions{Attributes: VisualBackdrop, Backdrop: SystemBackdrop(50)}); !errors.Is(err, ErrVisualOptions) {
			t.Fatal(err)
		}
		if state.OSAttributes&(VisualBackdrop|VisualPremultipliedAlpha) == VisualBackdrop|VisualPremultipliedAlpha {
			WindowHint(TransparentFramebuffer, False)
			opaque, err := CreateWindow(64, 48, "opaque preflight", nil, nil)
			if err != nil {
				t.Fatal(err)
			}
			defer opaque.Destroy()
			before, err := opaque.QueryWindowVisuals()
			if err != nil {
				t.Fatal(err)
			}
			result, err := opaque.ApplyWindowVisuals(WindowVisualOptions{Attributes: VisualBackdrop | VisualPremultipliedAlpha, Backdrop: BackdropMainWindow, PremultipliedAlpha: true})
			if !errors.Is(err, ErrSurfaceAlpha) || result.Accepted != 0 {
				t.Fatalf("opaque alpha enabled: %+v %v", result, err)
			}
			after, err := opaque.QueryWindowVisuals()
			if err != nil || before.Backdrop != after.Backdrop {
				t.Fatal("failed preflight changed backdrop")
			}
		}
		threadResult := make(chan error, 1)
		go func() {
			runtime.LockOSThread()
			defer runtime.UnlockOSThread()
			_, err := w.QueryWindowVisuals()
			threadResult <- err
		}()
		if err := <-threadResult; !errors.Is(err, ErrVisualThread) {
			t.Fatal("wrong thread accepted", err)
		}
	}
	w.Destroy()
	if _, err := w.QueryWindowVisuals(); !errors.Is(err, ErrWindowDestroyed) {
		t.Fatal("stale native handle queried", err)
	}
}
