package solution

import (
	"sort"
	"testing"
)

func TestHeapSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{1}},
		{name: "already sorted", input: []int{1, 2, 3, 4, 5}},
		{name: "reverse sorted", input: []int{9, 8, 7, 6, 5}},
		{name: "duplicates", input: []int{2, 9, 2, 1, 7, 1}},
		{name: "with negatives", input: []int{-6, 3, -1, 8, 0}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			heapSort(got)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("heapSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

