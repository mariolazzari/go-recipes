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

### Errors handling

```go
package main

import (
 "errors"
 "fmt"
 "io/fs"
 "log"
 "os"
)

func stopServer(pidFile string) error {
 file, err := os.Open(pidFile)
 if err != nil {
  return err
 }
 defer file.Close()

 var pid int
 if _, err := fmt.Fscanf(file, "%d", &pid); err != nil {
  return fmt.Errorf("invalid pid in %q - %w", pidFile, err)
 }

 if pid <= 0 {
  return fmt.Errorf("bad pid in %q - %d", pidFile, pid)
 }

 defer os.Remove(pidFile)

 log.Printf("stopping server with PID %d", pid)
 proc, err := os.FindProcess(pid)
 if err != nil {
  return fmt.Errorf("can't find process %d - %w", pid, err)
 }
 if err := proc.Kill(); err != nil {
  return fmt.Errorf("can't kill process %d - %w", pid, err)
 }

 return nil
}

func main() {
 if err := stopServer("httpd.pid"); err != nil {
  if errors.Is(err, fs.ErrNotExist) {
   fmt.Println("server not running")
  } else {
   log.Fatalf("error: %s", err)
  }
 }
}
```

### Defer

```go
package main

import (
 "encoding/csv"
 "log"
 "os"
)

func main() {
 items := []Item{
  {"m183x", "Magic Wand"},
  {"m184y", "Invisibility Cape"},
  {"m185z", "Levitation Spell"},
 }

 if err := writeItems("items.csv", items); err != nil {
  log.Fatal(err)
 }
}

type Item struct {
 SKU  string
 Name string
}

func writeItems(fileName string, items []Item) error {
 file, err := os.Create(fileName)
 if err != nil {
  return err
 }
 defer file.Close()

 row := []string{"sku", "name"}

 wtr := csv.NewWriter(file)
 defer wtr.Flush()

 if err := wtr.Write(row); err != nil {
  return err
 }

 for _, item := range items {
  row[0] = item.SKU
  row[1] = item.Name
  if err := wtr.Write(row); err != nil {
   return err
  }
 }

 return wtr.Error()
}
```

### Panic

```go
package main

import (
 "fmt"
)

func main() {
 nums := []int{1}
 // fmt.Println(secondToLast(nums)) // will panic
 fmt.Println(safeSecondToLast(nums))
}

func safeSecondToLast(nums []int) (i int, err error) {
 defer func() {
  if e := recover(); e != nil { // e is interface{}
   err = fmt.Errorf("%v", e)
  }
 }()

 return secondToLast(nums), nil
}

func secondToLast(nums []int) int {
 return nums[len(nums)-2]
}
```

### Challenge: write a filter

´´´go

```
