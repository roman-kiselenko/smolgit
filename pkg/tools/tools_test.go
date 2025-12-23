package tools

import "testing"

func TestNew(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic: %v", r)
		}
	}()
	New("test")
}
