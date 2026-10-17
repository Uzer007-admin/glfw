//go:build glfw_integration

package glfw

import (
	"runtime"
	"testing"
	"unsafe"
)

func TestGammaRampNullPlatformRoundTrip(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	InitHint(PlatformHint, int(PlatformNull))
	defer func() {
		Terminate()
		InitHint(PlatformHint, int(AnyPlatform))
	}()
	if err := Init(); err != nil {
		t.Fatal(err)
	}
	monitor := GetPrimaryMonitor()
	if monitor == nil {
		t.Fatal("null platform did not create a monitor")
	}
	ramp := monitor.GetGammaRamp()
	if ramp == nil || len(ramp.Red) != 256 {
		t.Fatal("null platform did not create a 256-entry ramp")
	}
	ramp.Red[42], ramp.Green[42], ramp.Blue[42] = 123, 456, 789
	monitor.SetGammaRamp(ramp)
	ramp.Red[42], ramp.Green[42], ramp.Blue[42] = 1, 2, 3
	dst := &GammaRamp{Red: make([]uint16, 0, 256), Green: make([]uint16, 0, 256), Blue: make([]uint16, 0, 256)}
	red := unsafe.SliceData(dst.Red[:cap(dst.Red)])
	result := monitor.GetGammaRampInto(dst)
	if result != dst || unsafe.SliceData(dst.Red) != red || dst.Red[42] != 123 || dst.Green[42] != 456 || dst.Blue[42] != 789 {
		t.Fatalf("copy/reuse failed: %v", dst)
	}
	if allocations := testing.AllocsPerRun(10, func() { monitor.GetGammaRampInto(dst) }); allocations != 0 {
		t.Fatalf("reused gamma buffers allocated: %f", allocations)
	}
	for joystick := Joystick1; joystick <= JoystickLast; joystick++ {
		if joystick.Present() {
			continue
		}
		if joystick.GetAxesInto(make([]float32, 0, 8)) != nil ||
			joystick.GetButtonsInto(make([]Action, 0, 16)) != nil ||
			joystick.GetHatsInto(make([]JoystickHatState, 0, 4)) != nil {
			t.Fatal("absent joystick returned non-nil state")
		}
		state := GamepadState{}
		state.Axes[0] = 0.75
		if joystick.GetGamepadStateInto(&state) || state.Axes[0] != 0.75 {
			t.Fatal("absent joystick changed gamepad output")
		}
		axes, buttons, hats := make([]float32, 0, 8), make([]Action, 0, 16), make([]JoystickHatState, 0, 4)
		if allocations := testing.AllocsPerRun(10, func() {
			joystick.GetAxesInto(axes)
			joystick.GetButtonsInto(buttons)
			joystick.GetHatsInto(hats)
			joystick.GetGamepadStateInto(&state)
		}); allocations != 0 {
			t.Fatalf("reusable joystick queries allocated: %f", allocations)
		}
		break
	}
}
