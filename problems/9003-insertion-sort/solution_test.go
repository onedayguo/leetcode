package solution

import (
	"sort"
	"testing"
)

func TestInsertionSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{3}},
		{name: "already sorted", input: []int{1, 2, 3, 4}},
		{name: "reverse sorted", input: []int{6, 5, 4, 3, 2, 1}},
		{name: "duplicates", input: []int{2, 2, 1, 3, 1}},
		{name: "with negatives", input: []int{-4, 7, 0, -1, 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			insertionSort(got)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("insertionSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

