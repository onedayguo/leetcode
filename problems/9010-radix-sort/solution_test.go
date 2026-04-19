package solution

import (
	"sort"
	"testing"
)

func TestRadixSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{0}},
		{name: "already sorted", input: []int{1, 2, 3, 10, 100}},
		{name: "reverse sorted", input: []int{1000, 100, 10, 1, 0}},
		{name: "duplicates", input: []int{170, 45, 75, 45, 24, 2, 66}},
		{name: "mixed digits", input: []int{5, 123, 9, 10000, 56, 1}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			radixSort(got)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("radixSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

