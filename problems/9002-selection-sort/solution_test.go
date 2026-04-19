package solution

import (
	"sort"
	"testing"
)

func TestSelectionSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{10}},
		{name: "already sorted", input: []int{1, 2, 3, 4}},
		{name: "reverse sorted", input: []int{9, 7, 5, 3, 1}},
		{name: "duplicates", input: []int{3, 1, 2, 1, 3}},
		{name: "with negatives", input: []int{-2, 5, 0, -9, 4}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			selectionSort(got)

			if len(got) != len(want) {
				t.Fatalf("selectionSort(%v) length = %d, want %d", tc.input, len(got), len(want))
			}
			for i := range got {
				if got[i] != want[i] {
					t.Fatalf("selectionSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

