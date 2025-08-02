# Go recipes: practical examples

## Basics

### Numbers

```go
package main

import "fmt"

func main() {
 nums := []int{1, 2, 3}
 fmt.Println(Mean(nums))

 nums = []int{1, 2, 3, 4}
 fmt.Println(Mean(nums))
}

func Sum(nums []int) int {
 total := 0
 for _, n := range nums {
  total += n
 }

 return total
}

func Mean(nums []int) float64 {
 sum := float64(Sum(nums))
 if sum == 0 {
  return 0
 }

 return sum / float64(len(nums))
}
```

### Slices

```go

```
