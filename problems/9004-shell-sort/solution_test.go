package solution

import (
	"sort"
	"testing"
)

func TestShellSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{99}},
		{name: "already sorted", input: []int{-3, -1, 0, 2, 4}},
		{name: "reverse sorted", input: []int{8, 7, 6, 5, 4, 3, 2, 1}},
		{name: "duplicates", input: []int{5, 3, 5, 3, 1}},
		{name: "with negatives", input: []int{9, -2, 4, -7, 0}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			shellSort(got)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("shellSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

