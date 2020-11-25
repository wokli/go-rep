package main

import "testing"

func TestDummy(t *testing.T) {
	res := true
	if !res {
		t.Errorf("Dummy test failed")
	}
}
