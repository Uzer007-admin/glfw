package glfw

import "testing"

func TestMonitorConfigurationDispatch(t *testing.T) {
	previous := fMonitorConfigurationHolder
	t.Cleanup(func() { fMonitorConfigurationHolder = previous })
	fMonitorConfigurationHolder = nil
	goMonitorConfigurationCB()
	calls := 0
	fMonitorConfigurationHolder = func() { calls++ }
	goMonitorConfigurationCB()
	if calls != 1 {
		t.Fatalf("configuration notification calls: %d", calls)
	}
	fMonitorConfigurationHolder = nil
	goMonitorConfigurationCB()
	if calls != 1 {
		t.Fatal("disabled callback was invoked")
	}
}
