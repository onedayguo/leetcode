package solution

import (
	"sort"
	"testing"
)

func TestBucketSort(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
	}{
		{name: "empty", input: []int{}},
		{name: "single", input: []int{11}},
		{name: "already sorted", input: []int{0, 2, 4, 6}},
		{name: "reverse sorted", input: []int{10, 8, 6, 4, 2, 0}},
		{name: "duplicates", input: []int{5, 1, 5, 3, 1, 0}},
		{name: "sparse values", input: []int{100, 3, 250, 40, 7}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := append([]int(nil), tc.input...)
			want := append([]int(nil), tc.input...)
			sort.Ints(want)

			bucketSort(got)

			for i := range want {
				if got[i] != want[i] {
					t.Fatalf("bucketSort(%v) = %v, want %v", tc.input, got, want)
				}
			}
		})
	}
}

