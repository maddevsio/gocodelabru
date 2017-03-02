package lru

import "testing"

// Test that Add returns true/false if an eviction occurred
func TestLRU_Add(t *testing.T) {

	l, err := New(1)
	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if l.Add(1, 1) == true {
		t.Errorf("should not have an eviction")
	}
	if l.Add(2, 2) == false {
		t.Errorf("should have an eviction")
	}
}
