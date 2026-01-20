package proc

import (
	"fmt"
	"testing"
)

type MockResolver struct {
	Response string
	Err      error
}

func (m *MockResolver) Resolve(pid uint32) (string, error) {
	return m.Response, m.Err
}

func TestGetProcessNameSuccess(t *testing.T) {
	mock := &MockResolver{Response: "my-app"}
	mgr := NewProcessManager(mock)

	pid := uint32(1234)
	name, ok := mgr.GetProcessName(&pid)

	if !ok || name != "my-app" {
		t.Errorf("Expected my-app, got %s", name)
	}
}

func TestGetProcessNameError(t *testing.T) {
	mock := &MockResolver{Response: "", Err: fmt.Errorf("disk read error")}
	mgr := NewProcessManager(mock)

	pid := uint32(1234)
	name, ok := mgr.GetProcessName(&pid)

	if ok || name != "" {
		t.Errorf("Expected empty app name and false, got %t", ok)
	}
}
