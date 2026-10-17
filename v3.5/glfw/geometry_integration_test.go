//go:build glfw_integration

package glfw

import (
	"runtime"
	"testing"
)

func benchmarkGeometry(b *testing.B, legacy bool) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if err := Init(); err != nil {
		b.Fatal(err)
	}
	defer Terminate()
	WindowHint(Visible, False)
	window, err := CreateWindow(128, 96, "geometry benchmark", nil, nil)
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
		} else {
			window.GetSurfaceGeometry()
		}
	}
}
func BenchmarkLegacySurfaceGeometry(b *testing.B) { benchmarkGeometry(b, true) }
func BenchmarkSurfaceGeometry(b *testing.B)       { benchmarkGeometry(b, false) }
