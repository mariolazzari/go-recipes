package main

import "fmt"

func filter(pred func(int) bool, values []int) []int {
	var filtered []int
	for _, v := range values {
		if pred(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

func isOdd(n int) bool {
	return n%2 == 1
}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8}
	filtered := filter(isOdd, nums)
	fmt.Println("All items:", nums)
	fmt.Println("Filtered :", filtered)
}
