package jso

import (
	"strconv"
	"testing"
)

func TestFilterMap_IntToInt(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	result := filterMap(input, func(i int) (int, bool) {
		if i%2 == 0 {
			return i * 2, true
		}
		return 0, false
	})
	expected := []int{4, 8}
	if len(result) != len(expected) {
		t.Fatalf("expected length %d, got %d", len(expected), len(result))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected, result)
		}
	}
}
func TestFilterMap_EmptyInput(t *testing.T) {
	input := make([]int, 0)
	result := filterMap(input, func(i int) (int, bool) { return i, true })
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}
func TestFilterMap_AllFalse(t *testing.T) {
	input := []string{"a", "b", "c"}
	result := filterMap(input, func(s string) (string, bool) { return s, false })
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}
func TestFilterMap_AllTrue(t *testing.T) {
	input := []string{"a", "b"}
	result := filterMap(input, func(s string) (string, bool) { return s + "!", true })
	expected := []string{"a!", "b!"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected, result)
		}
	}
}
func TestFilterMap_IntToString(t *testing.T) {
	input := []int{1, 2, 3}
	result := filterMap(input, func(i int) (string, bool) {
		if i%2 == 1 {
			return strconv.Itoa(i), true
		}
		return "", false
	})
	expected := []string{"1", "3"}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("expected %v, got %v", expected, result)
		}
	}
}
