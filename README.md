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

```go
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
```

## Working with time

### Time arithmetic

```go
package main

import (
 "fmt"
 "time"
)

func isBusinessDay(date time.Time) bool {
 wday := date.Weekday()
 if wday == time.Saturday || wday == time.Sunday {
  return false
 }

 return true
}

func nextBusinessDay(date time.Time) time.Time {
 const day = 24 * time.Hour
 for {
  date = date.Add(day)
  if isBusinessDay(date) {
   break
  }
 }

 return date
}

func main() {
 date := time.Date(2021, time.December, 31, 0, 0, 0, 0, time.UTC)
 fmt.Println(date, date.Weekday()) // 2021-12-31 00:00:00 +0000 UTC Friday
 nbd := nextBusinessDay(date)
 fmt.Println(nbd, nbd.Weekday()) // 2022-01-03 00:00:00 +0000 UTC Monday

 date = time.Date(2022, time.January, 4, 0, 0, 0, 0, time.UTC)
 fmt.Println(date, date.Weekday()) // 2022-01-04 00:00:00 +0000 UTC Tuesday
 nbd = nextBusinessDay(date)
 fmt.Println(nbd, nbd.Weekday()) // 2022-01-05 00:00:00 +0000 UTC Wednesday
}
```

### Time measuring

```go
package main

import (
 "fmt"
 "log"
 "time"
)

func timeit(name string) func() {
 start := time.Now()

 return func() {
  duration := time.Since(start)
  log.Printf("%s took %s", name, duration)
 }
}

func dot(v1, v2 []float64) (float64, error) {
 defer timeit("dot")()

 if len(v1) != len(v2) {
  return 0, fmt.Errorf("dot of different size vectors")
 }

 d := 0.0
 for i, val1 := range v1 {
  val2 := v2[i]
  d += val1 * val2
 }

 return d, nil
}

func main() {
 v := []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
 fmt.Println(dot(v, v))
}
```

### Formatting times

```go
package main

import (
 "fmt"
 "time"
)

func main() {
 lennon := time.Date(1940, time.October, 9, 18, 30, 0, 0, time.UTC)
 fmt.Println(lennon) // 1940-10-09 18:30:00 +0000 UTC

 fmt.Println(lennon.Format("2006-01-02"))  // 1940-10-09
 fmt.Println(lennon.Format("Mon, Jan 02")) // Wed, Oct 09

 fmt.Println(lennon.Format(time.RFC3339Nano)) // 1940-10-09T18:30:00Z

 d := 3500 * time.Millisecond
 fmt.Println(d) // 3.5s
}
```

### Parsing times

```go
package main

import (
 "fmt"
 "time"
)

func main() {
 ts := "June 18, 1942"

 t, err := time.Parse("January 02, 2006", ts)
 if err != nil {
  fmt.Printf("error: %s\n", err)
 } else {
  fmt.Println(t) // 1942-06-18 00:00:00 +0000 UTC
 }

 ds := "2700ms"
 d, err := time.ParseDuration(ds)
 if err != nil {
  fmt.Printf("error: %s\n", err)
 } else {
  fmt.Println(d) // 2.7s
 }
}
```

### Time zones

```go
package main

import (
 "fmt"
 "time"
)

func main() {
 chi, err := time.LoadLocation("America/Chicago")
 if err != nil {
  fmt.Printf("error: %s", err)
  return
 }

 chiTime := time.Date(2021, time.February, 28, 19, 30, 0, 0, chi)
 fmt.Println("Chicago:", chiTime) // Chicago: 2021-02-28 19:30:00 -0600 CST

 nyc, err := time.LoadLocation("America/New_York")
 if err != nil {
  fmt.Printf("error: %s", err)
  return
 }

 nycTime := chiTime.In(nyc)
 fmt.Println("NYC:", nycTime) // NYC: 2021-02-28 20:30:00 -0500 EST
}
```

### Challenge: time zones

```go
package main

import (
 "fmt"
 "time"
)

// tsConvert convert time stamp in "YYYY-MM-DDTHH:MM" format from one time zone to another
func tsConvert(ts, from, to string) (string, error) {
 fromTz, err := time.LoadLocation(from)
 if err != nil {
  return "", err
 }

 toTz, err := time.LoadLocation(to)
 if err != nil {
  return "", err
 }

 const format = "2006-01-02T15:04"
 fromTime, err := time.ParseInLocation(format, ts, fromTz)
 if err != nil {
  return "", err
 }

 toTime := fromTime.In(toTz)
 return toTime.Format(format), nil
}

func main() {
 ts := "2021-03-08T19:12"
 out, err := tsConvert(ts, "America/Los_Angeles", "Asia/Jerusalem")
 if err != nil {
  fmt.Printf("error: %s", err)
  return
 }

 fmt.Println(out) // 2021-03-09T05:12
}
```

## Strings

### String formatting

```go
package main

import (
 "fmt"
 "io"
 "log"
 "os"
)

// Trade represents a trade
type Trade struct {
 Symbol string
 Volume int
 Price  float64
}

// genReport generates a fixed with report in the format
// Symbol: 10 chars, left padded
// Volume: 4 digits, 0 padded
// Price: 2 digits after the decimal
func genReport(w io.Writer, trades []Trade) {
 for i, t := range trades {
  log.Printf("%d: %#v", i, t)
  // ... 2: main.Trade{Symbol:"BRK-A", Volume:1, Price:399100}
  fmt.Fprintf(w, "%-10s %04d %.2f\n", t.Symbol, t.Volume, t.Price)
  // MSFT       0231 234.57
 }
}

func main() {
 log.SetPrefix("LOG: ")

 trades := []Trade{
  {"MSFT", 231, 234.57},
  {"TSLA", 123, 686.75},
  {"BRK-A", 1, 399100},
 }
 genReport(os.Stdout, trades)
}
```

### Unicode

```go
package main

import (
 "fmt"
 "unicode/utf8"
)

func lineLength(words []string) int {
 total := 0
 for _, word := range words {
  total += utf8.RuneCountInString(word)
 }

 numSpaces := len(words) - 1
 return total + numSpaces
}

func main() {
 words := []string{"«", "Don't", "Panic", "»"}
 fmt.Println(lineLength(words)) // 15
}
```

### Case insensitive

```go
package main

import (
 "fmt"
 "strings"
)

// Letter in Greek
type Letter struct {
 Symbol  string
 English string
}

var letters = []Letter{
 {"Σ", "Sigma"},
 // TODO
}

// englishFor return the English name for a greek letter
func englishFor(greek string) (string, error) {
 for _, letter := range letters {
  if strings.EqualFold(greek, letter.Symbol) {
   return letter.English, nil
  }
 }

 return "", fmt.Errorf("unknown greek letter: %#v", greek)
}

func main() {
 fmt.Println(englishFor("Σ"))
 fmt.Println(englishFor("σ"))
 fmt.Println(englishFor("ς"))
}
```

### Regular expressions

```go
package main

import (
 "fmt"
 "log"
 "regexp"
 "strconv"
)

/*
12 shares of MSFT for $234.57
10 shares of TSLA for $692.4
*/
var transRe = regexp.MustCompile(`(\d+) shares of ([A-Z]+) for \$(\d+(\.\d+)?)`)

// Transaction is a b
type Transaction struct {
 Symbol string
 Volume int
 Price  float64
}

func parseLine(line string) (Transaction, error) {
 matches := transRe.FindStringSubmatch(line)
 if matches == nil {
  return Transaction{}, fmt.Errorf("bad line: %q", line)
 }
 var t Transaction
 t.Symbol = matches[2]
 t.Volume, _ = strconv.Atoi(matches[1])
 t.Price, _ = strconv.ParseFloat(matches[3], 64)
 return t, nil
}

func main() {
 line := "12 shares of MSFT for $234.57"
 t, err := parseLine(line)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Printf("%+v\n", t) // {Symbol:MSFT Volume:12 Price:234.57}
}
```

### Reading text files

```go
package main

import (
 "bufio"
 "fmt"
 "io"
 "log"
 "os"
 "strings"
)

// grep returns lines in r that contain term
func grep(r io.Reader, term string) ([]string, error) {
 var matches []string
 s := bufio.NewScanner(r)
 for s.Scan() {
  if strings.Contains(s.Text(), term) {
   matches = append(matches, s.Text())
  }
 }

 if err := s.Err(); err != nil {
  return nil, err
 }

 return matches, nil
}

func main() {
 file, err := os.Open("journal.txt")
 if err != nil {
  log.Fatal(err)
 }
 defer file.Close()

 matches, err := grep(file, "System is rebooting")
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("%d reboots\n", len(matches))
}
```

### Challenge: text

```go
package main

import (
 "bufio"
 "fmt"
 "log"
 "os"
 "regexp"
)

var cmdRe = regexp.MustCompile(`;go ([a-z]+)`)

// cmdFreq returns the frequency of "go" subcommand usage in ZSH history
func cmdFreq(fileName string) (map[string]int, error) {
 file, err := os.Open(fileName)
 if err != nil {
  return nil, err
 }
 defer file.Close()

 freqs := make(map[string]int)
 s := bufio.NewScanner(file)
 for s.Scan() {
  matches := cmdRe.FindStringSubmatch(s.Text())
  if len(matches) == 0 {
   continue
  }
  cmd := matches[1]
  freqs[cmd]++
 }

 if err := s.Err(); err != nil {
  return nil, err
 }

 return freqs, nil
}

func main() {
 freqs, err := cmdFreq("./zsh_history")
 if err != nil {
  log.Fatal(err)
 }

 for cmd, count := range freqs {
  fmt.Printf("%s -> %d\n", cmd, count)
 }
}
```
