package glfw

import (
	"errors"
	"testing"
)

func TestDestroyedWindowVisuals(t *testing.T) {
	for _, w := range []*Window{nil, {}} {
		if _, err := w.QueryWindowVisuals(); !errors.Is(err, ErrWindowDestroyed) {
			t.Fatal(err)
		}
		if _, err := w.ApplyWindowVisuals(WindowVisualOptions{}); !errors.Is(err, ErrWindowDestroyed) {
			t.Fatal(err)
		}
	}
}

func TestCreatedSurfaceAlpha(t *testing.T) {
	if (WindowVisualState{TransparentFramebuffer: true}).WGLAlphaAvailable() {
		t.Fatal("hint treated as real alpha")
	}
	if (WindowVisualState{PixelFormat: 1, AlphaBits: 8}).WGLAlphaAvailable() {
		t.Fatal("opaque surface treated as transparent")
	}
	if !(WindowVisualState{TransparentFramebuffer: true, PixelFormat: 1, AlphaBits: 8}).WGLAlphaAvailable() {
		t.Fatal("missing actual alpha")
	}
}
