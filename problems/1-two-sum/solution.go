package solution

import "fmt"

func twoSum(nums []int, target int) []int {
	var existMap = make(map[int]int)
	for i := 0; i < len(nums); i++ {
		if index, exist := existMap[target-nums[i]]; exist {
			return []int{index, i}
		}
		existMap[nums[i]] = i
	}
	return []int{}
}

func main() {
	// 本地测试入口（可选）
	result := twoSum([]int{2, 7, 11, 15}, 9)
	fmt.Println("Result:", result)
}
