package glfw

import (
	"reflect"
	"testing"
	"unsafe"
)

func TestNativeArrayCopiesOwnDataAndReuseCapacity(t *testing.T) {
	source := []float32{-1, 0, 0.5, 1}
	dst := make([]float32, 0, 8)
	first := unsafe.SliceData(dst[:cap(dst)])
	result := copyNativeSlice[float32](dst, unsafe.Pointer(unsafe.SliceData(source)), len(source))
	if unsafe.SliceData(result) != first || !reflect.DeepEqual(result, source) {
		t.Fatalf("copy/reuse: %v", result)
	}
	source[0] = 0.25
	if result[0] != -1 {
		t.Fatal("result borrows native storage")
	}
	if copyNativeSlice[float32](dst, nil, 0) != nil {
		t.Fatal("absent native array must be nil")
	}
	if empty := copyNativeSlice[float32](nil, unsafe.Pointer(unsafe.SliceData(source)), 0); empty == nil || len(empty) != 0 {
		t.Fatal("present empty native array must be non-nil")
	}
	bytes := []byte{0, 1, 255}
	actions := copyNativeBytes[Action](nil, unsafe.Pointer(unsafe.SliceData(bytes)), len(bytes))
	if !reflect.DeepEqual(actions, []Action{0, 1, 255}) {
		t.Fatal(actions)
	}
	bytes[1] = 0
	if actions[1] != 1 {
		t.Fatal("button array borrows native storage")
	}
	gamma := []uint16{0, 32768, 65535}
	copy := copyNativeSlice[uint16](nil, unsafe.Pointer(unsafe.SliceData(gamma)), len(gamma))
	gamma[1] = 1
	if copy[1] != 32768 {
		t.Fatal("gamma array borrows native storage")
	}
}

func TestInvalidGammaRampPanicsBeforeNativeCall(t *testing.T) {
	for _, ramp := range []*GammaRamp{nil, {}, {Red: []uint16{1}}, {Red: []uint16{1}, Green: []uint16{1}, Blue: []uint16{1, 2}}} {
		func() {
			defer func() {
				failure, ok := recover().(*Error)
				if !ok || failure.Code != invalidValue {
					t.Fatalf("invalid ramp must panic with InvalidValue; got %v", failure)
				}
			}()
			new(Monitor).SetGammaRamp(ramp)
		}()
	}
	if size := validateGammaRamp(&GammaRamp{Red: []uint16{1}, Green: []uint16{2}, Blue: []uint16{3}}); size != 1 {
		t.Fatal(size)
	}
}

func BenchmarkNativeArrayCopy(b *testing.B) {
	source := make([]float32, 64)
	destination := make([]float32, 0, 64)
	pointer := unsafe.Pointer(unsafe.SliceData(source))
	b.ReportAllocs()
	for b.Loop() {
		destination = copyNativeSlice[float32](destination, pointer, len(source))
	}
}
