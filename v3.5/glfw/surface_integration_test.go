//go:build glfw_integration

package glfw

import (
	"runtime"
	"testing"
)

func TestSurfaceMetricsMatchesNativeQueries(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	defer Terminate()
	WindowHint(Visible, False)
	WindowHint(TransparentFramebuffer, True)
	window, err := CreateWindow(128, 96, "surface metrics", nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer window.Destroy()
	state := window.GetSurfaceMetrics()
	geometry := window.GetSurfaceGeometry()
	if geometry.Width != state.Width || geometry.Height != state.Height || geometry.FramebufferWidth != state.FramebufferWidth || geometry.FramebufferHeight != state.FramebufferHeight || geometry.ContentScaleX != state.ContentScaleX || geometry.ContentScaleY != state.ContentScaleY {
		t.Fatalf("geometry differs: %+v vs %+v", geometry, state)
	}
	if allocations := testing.AllocsPerRun(20, func() { window.GetSurfaceGeometry() }); allocations != 0 {
		t.Fatalf("geometry allocated: %f", allocations)
	}
	width, height := window.GetSize()
	framebufferWidth, framebufferHeight := window.GetFramebufferSize()
	scaleX, scaleY := window.GetContentScale()
	if state.Width != width || state.Height != height || state.FramebufferWidth != framebufferWidth || state.FramebufferHeight != framebufferHeight ||
		state.ContentScaleX != scaleX || state.ContentScaleY != scaleY || state.Visible || state.TransparentFramebuffer != (window.GetAttrib(TransparentFramebuffer) != 0) {
		t.Fatalf("native metrics mismatch: %+v", state)
	}
	if _, _, valid := state.PixelScale(); !valid {
		t.Fatal("positive window lacks pixel scale")
	}
	if allocations := testing.AllocsPerRun(20, func() { window.GetSurfaceMetrics() }); allocations != 0 {
		t.Fatalf("surface query allocated: %f", allocations)
	}
	var tracker SurfaceTracker
	tracked, changes, err := window.UpdateSurface(&tracker)
	if err != nil || tracked.Metrics != state || tracked.Drawable || tracked.ResourceGeneration != 1 || changes&SurfaceGeometryChanged == 0 {
		t.Fatalf("tracker differs from actual hidden surface: %+v %v %v", tracked, changes, err)
	}
	if allocations := testing.AllocsPerRun(20, func() {
		if _, changes, err := window.UpdateSurface(&tracker); err != nil || changes != 0 {
			t.Fatal(changes, err)
		}
	}); allocations != 0 {
		t.Fatalf("native tracker allocated: %f", allocations)
	}
	window.Destroy()
	if _, _, err := window.UpdateSurface(&tracker); err != ErrWindowDestroyed {
		t.Fatal("destroyed tracker native access", err)
	}
}

func benchmarkSurface(b *testing.B, legacy bool) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := Init(); err != nil {
		b.Fatal(err)
	}
	defer Terminate()
	WindowHint(Visible, False)
	window, err := CreateWindow(128, 96, "surface benchmark", nil, nil)
	if err != nil {
		b.Fatal(err)
	}
	defer window.Destroy()
	b.ReportAllocs()
	for b.Loop() {
		if legacy {
			window.GetSize()
			window.GetFramebufferSize()
			window.GetContentScale()
			for _, attribute := range []Hint{Focused, Iconified, Visible, Maximized, TransparentFramebuffer} {
				window.GetAttrib(attribute)
			}
		} else {
			window.GetSurfaceMetrics()
		}
	}
}

func BenchmarkLegacySurfaceMetrics(b *testing.B) { benchmarkSurface(b, true) }
func BenchmarkSurfaceMetrics(b *testing.B)       { benchmarkSurface(b, false) }
