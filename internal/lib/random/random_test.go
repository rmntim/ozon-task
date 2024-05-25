package random

import "testing"

func TestNewRandomString(t *testing.T) {
	s := NewRandomString(10)
	if len(s) != 10 {
		t.Error("length should be 10")
	}
}
