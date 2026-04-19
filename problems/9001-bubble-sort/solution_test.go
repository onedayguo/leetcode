package solution

import (
	"sort"
	"testing"
)

func TestBubbleSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{7}},
		{name: "already sorted", input: []int{1, 2, 3, 4}},
		{name: "reverse sorted", input: []int{5, 4, 3, 2, 1}},
		{name: "duplicates", input: []int{4, 2, 4, 1, 2}},
		{name: "with negatives", input: []int{3, -1, 0, -5, 8}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			bubbleSort(got)

			if !equalSlices(got, want) {
				t.Fatalf("bubbleSort(%v) = %v, want %v", tc.input, got, want)
			}
		})
	}
}

func equalSlices(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}


