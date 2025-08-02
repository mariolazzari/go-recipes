# Go recipes: practical examples

## Basics

### Numbers

Compute average

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

Compute median

```go
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
```

### Maps

Word frequency

```go
package main

import (
 "fmt"
 "strings"
)

func main() {
 msg := "To be or not to be"
 fmt.Println("Frequency:", frequency(msg))
}

func frequency(str string) map[string]int {
 freq := make(map[string]int)

 words := strings.SplitSeq(strings.ToLower(str), " ")
 for word := range words {
  freq[word] += 1
 }

 return freq
}
```
