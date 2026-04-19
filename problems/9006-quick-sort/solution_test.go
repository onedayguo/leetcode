package solution

import (
	"sort"
	"testing"
)

func TestQuickSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{42}},
		{name: "already sorted", input: []int{-2, -1, 0, 1, 2}},
		{name: "reverse sorted", input: []int{7, 6, 5, 4, 3, 2, 1}},
		{name: "duplicates", input: []int{3, 5, 3, 2, 8, 2}},
		{name: "with negatives", input: []int{10, -3, 0, 4, -8, 7}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			quickSort(got)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("quickSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

