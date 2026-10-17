package glfw

import "testing"

func TestSurfacePixelScaleIsIndependentOfContentScale(t *testing.T) {
	state := SurfaceMetrics{Width: 800, Height: 600, FramebufferWidth: 1200, FramebufferHeight: 1200, ContentScaleX: 1.25, ContentScaleY: 1.25}
	x, y, valid := state.PixelScale()
	if !valid || x != 1.5 || y != 2 || state.ContentScaleX != 1.25 {
		t.Fatalf("pixel/content scale conflated: %+v", state)
	}
	state.FramebufferWidth = 0
	if x, y, valid := state.PixelScale(); valid || x != 0 || y != 0 {
		t.Fatal("zero framebuffer reported as usable")
	}
}
