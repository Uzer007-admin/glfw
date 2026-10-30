//go:build glfw_integration

package glfw

import (
	"runtime"
	"testing"
)

func TestMonitorConfigurationLifecycle(t *testing.T) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	for i := 0; i < 2; i++ {
		if err := Init(); err != nil {
			t.Fatal(err)
		}
		SetMonitorConfigurationCallback(func() {})
		if previous := SetMonitorConfigurationCallback(nil); previous == nil {
			Terminate()
			t.Fatal("callback replacement lost the previous function")
		}
		SetMonitorConfigurationCallback(func() {})
		Terminate()
		if fMonitorConfigurationHolder != nil {
			t.Fatal("Terminate retained a monitor configuration callback")
		}
	}
}
