//go:build glfw_bench

package glfw

import (
	"testing"
	"unsafe"
)

func BenchmarkLegacyAxisCopy(b *testing.B) {
	source := make([]float32, 64)
	destination := make([]float32, 64)
	pointer := unsafe.Pointer(unsafe.SliceData(source))
	b.ReportAllocs()
	for b.Loop() {
		legacyAxisCopy(destination, pointer)
	}
}
