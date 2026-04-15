package project

import (
	"strings"
	"testing"
)

func TestProblemFolder(t *testing.T) {
	tests := []struct {
		id   string
		slug string
		want string
	}{
		{id: "1", slug: "Two Sum", want: "1-two-sum"},
		{id: "15", slug: "3Sum", want: "15-3sum"},
		{id: "", slug: "A_B", want: "a-b"},
	}

	for _, tt := range tests {
		got := ProblemFolder(tt.id, tt.slug)
		if got != tt.want {
			t.Fatalf("ProblemFolder(%q,%q) = %q, want %q", tt.id, tt.slug, got, tt.want)
		}
	}
}

func TestExtractSubmissionCode(t *testing.T) {
	input := `package main

import "fmt"

func twoSum(nums []int, target int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			if nums[i]+nums[j] == target {
				return []int{i, j}
			}
		}
	}
	return nil
}

func main() {
	// 本地测试入口（可选）
	result := twoSum([]int{2, 7, 11, 15}, 9)
	fmt.Println("Result:", result)
}
`

	output := ExtractSubmissionCode(input)

	// Should NOT contain package or import
	if strings.Contains(output, "package") {
		t.Fatal("output should not contain 'package'")
	}
	if strings.Contains(output, "import") {
		t.Fatal("output should not contain 'import'")
	}
	if strings.Contains(output, "func main()") {
		t.Fatal("output should not contain 'func main()'")
	}

	// Should contain twoSum function
	if !strings.Contains(output, "func twoSum") {
		t.Fatal("output should contain 'func twoSum'")
	}
	if !strings.Contains(output, "return []int{i, j}") {
		t.Fatal("output should contain twoSum logic")
	}

	t.Log("Extracted code:")
	t.Log(output)
}
