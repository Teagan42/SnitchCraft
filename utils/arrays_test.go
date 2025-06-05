package utils

import (
	"testing"
)

func TestMap(t *testing.T) {
	ints := []int{1, 2, 3}
	double := func(x int) int { return x * 2 }
	got := Map(ints, double)
	want := []int{2, 4, 6}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Map: got %v, want %v", got, want)
			break
		}
	}

	// Test with empty input
	gotEmpty := Map([]int{}, double)
	if len(gotEmpty) != 0 {
		t.Errorf("Map with empty input: got %v, want []", gotEmpty)
	}
}

func TestDo(t *testing.T) {
	ints := []int{1, 2, 3}
	sum := 0
	add := func(x int) { sum += x }
	out := Do(ints, add)
	if sum != 6 {
		t.Errorf("Do: sum = %d, want 6", sum)
	}
	if len(out) != len(ints) {
		t.Errorf("Do: output length = %d, want %d", len(out), len(ints))
	}
	for i := range ints {
		if out[i] != ints[i] {
			t.Errorf("Do: output = %v, want %v", out, ints)
			break
		}
	}
}

func TestFilter(t *testing.T) {
	ints := []int{1, 2, 3, 4, 5}
	isEven := func(x int) bool { return x%2 == 0 }
	got := Filter(ints, isEven)
	want := []int{2, 4}
	if len(got) != len(want) {
		t.Errorf("Filter: got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Filter: got %v, want %v", got, want)
			break
		}
	}

	// Test with no matches
	none := Filter(ints, func(x int) bool { return x > 10 })
	if len(none) != 0 {
		t.Errorf("Filter with no matches: got %v, want []", none)
	}

	// Test with empty input
	empty := Filter([]int{}, isEven)
	if len(empty) != 0 {
		t.Errorf("Filter with empty input: got %v, want []", empty)
	}
}
