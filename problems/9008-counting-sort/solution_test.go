package solution

import (
	"sort"
	"testing"
)

func TestCountingSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{7}},
		{name: "already sorted", input: []int{0, 1, 2, 3, 4}},
		{name: "reverse sorted", input: []int{9, 7, 5, 3, 1}},
		{name: "duplicates", input: []int{4, 2, 4, 2, 1, 1}},
		{name: "includes zero", input: []int{0, 5, 0, 2, 3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			countingSort(got)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("countingSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

