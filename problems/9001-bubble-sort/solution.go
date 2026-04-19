package solution

func bubbleSort(num []int) {
	// TODO: implement bubble sort.
	for i := 0; i < len(num); i++ {
		var isSort = false
		for j := 0; j < len(num)-i-1; j++ {
			if num[j] > num[j+1] {
				isSort = true
				num[j], num[j+1] = num[j+1], num[j]
			}
		}
		if !isSort {
			return
		}
	}
}
