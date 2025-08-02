package main

import (
	"fmt"
	"sort"
)

func main() {
	nums := []float64{2, 4, 5, 6, 7, 8, 9}
	median := Median(nums)
	fmt.Println(median)

	nums = []float64{1, 2, 3, 4, 5, 6, 8, 9}
	median = Median(nums)
	fmt.Println(median)
}

func Median(nums []float64) float64 {
	// slice size
	count := len(nums)
	// copy into a sorted slice
	vals := make([]float64, count)
	copy(vals, nums)
	sort.Float64s(vals)

	if count%2 == 0 {
		// even elements
		return (nums[count/2-1] + nums[count/2]) / 2
	}
	// odd elements
	return nums[count/2]

}
