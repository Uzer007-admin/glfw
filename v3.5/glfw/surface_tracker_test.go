package glfw

import (
	"math"
	"testing"
)

func TestSurfaceTrackerGenerationsAndScheduling(t *testing.T) {
	var tracker SurfaceTracker
	m := SurfaceMetrics{Width: 80, Height: 40, FramebufferWidth: 100, FramebufferHeight: 60, ContentScaleX: 1.25, ContentScaleY: 1.5, Visible: true}
	s, c := tracker.Update(m)
	if s.Generation != 1 || s.ResourceGeneration != 1 || !s.Drawable || c&SurfaceGeometryChanged == 0 {
		t.Fatal(s, c)
	}
	if s, c = tracker.Update(m); c != 0 || s.Generation != 1 {
		t.Fatal("static surface", s, c)
	}
	m.Focused = true
	s, c = tracker.Update(m)
	if c != SurfaceFocusChanged || s.ResourceGeneration != 1 {
		t.Fatal("focus invalidated pixels", s, c)
	}
	m.Iconified = true
	s, c = tracker.Update(m)
	if s.Drawable || c != SurfacePresentationChanged || s.ResourceGeneration != 1 {
		t.Fatal(s, c)
	}
	m.Iconified = false
	m.FramebufferHeight = 0
	s, _ = tracker.Update(m)
	if s.Drawable || s.ResourceGeneration != 2 {
		t.Fatal(s)
	}
	m.FramebufferHeight = 60
	m.ContentScaleX = 1.75
	s, c = tracker.Update(m)
	if !s.Drawable || c&SurfaceScaleChanged == 0 || s.ResourceGeneration != 3 {
		t.Fatal(s, c)
	}
	before := s.Generation
	s = tracker.InvalidateContext()
	if s.Generation != before+1 || s.ResourceGeneration != 4 {
		t.Fatal(s)
	}
	if _, c = tracker.Update(m); c != 0 {
		t.Fatal(c)
	}
}

func TestSurfacePixelRegionsFractionalClipping(t *testing.T) {
	m := SurfaceMetrics{Width: 80, Height: 40, FramebufferWidth: 100, FramebufferHeight: 60, ContentScaleX: 2, ContentScaleY: 2}
	r, ok := m.FramebufferRegion(1.1, 2.1, 10.1, 8.1)
	if !ok || r != (FramebufferRegion{1, 3, 13, 13}) {
		t.Fatal(r, ok)
	}
	if y, ok := r.BottomLeftY(60); !ok || y != 44 {
		t.Fatal(y, ok)
	}
	r, ok = m.FramebufferRegion(-3, -2, 100, 100)
	if !ok || r != (FramebufferRegion{0, 0, 100, 60}) {
		t.Fatal(r, ok)
	}
	for _, rect := range [][4]float64{{math.NaN(), 0, 1, 1}, {0, 0, math.Inf(1), 1}, {80, 0, 1, 1}, {0, 0, -1, 1}} {
		if _, ok = m.FramebufferRegion(rect[0], rect[1], rect[2], rect[3]); ok {
			t.Fatal(rect)
		}
	}
	m.FramebufferHeight = 0
	if _, ok = m.FramebufferRegion(0, 0, 1, 1); ok {
		t.Fatal("zero framebuffer accepted")
	}
	var w *Window
	if _, _, err := w.UpdateSurface(&SurfaceTracker{}); err != ErrWindowDestroyed {
		t.Fatal(err)
	}
}

func BenchmarkSurfaceTrackerStatic(b *testing.B) {
	var tracker SurfaceTracker
	m := SurfaceMetrics{Width: 800, Height: 600, FramebufferWidth: 1200, FramebufferHeight: 900, Visible: true, ContentScaleX: 1.5, ContentScaleY: 1.5}
	tracker.Update(m)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		tracker.Update(m)
	}
}
