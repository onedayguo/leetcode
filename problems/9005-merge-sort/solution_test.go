package solution

import (
	"sort"
	"testing"
)

func TestMergeSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{100}},
		{name: "already sorted", input: []int{1, 2, 3, 4, 5}},
		{name: "reverse sorted", input: []int{5, 4, 3, 2, 1}},
		{name: "duplicates", input: []int{4, 1, 4, 2, 2, 3}},
		{name: "with negatives", input: []int{-10, 3, 0, -1, 8, -5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			mergeSort(got)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("mergeSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

