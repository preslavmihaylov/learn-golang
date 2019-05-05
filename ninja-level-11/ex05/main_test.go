package main

import (
	"testing"
)

func TestSqrtWithNine(t *testing.T) {
	res, err := Sqrt(9)
	if res != 3 && err != nil {
		t.Errorf("expected 3 and no error, got %v, err:%v", res, err)
	}
}

func TestSqrtWithNegativeNumber(t *testing.T) {
	res, err := Sqrt(-10)
	if err == nil {
		t.Errorf("expected error, got %v, err:%v", res, err)
	}
}
