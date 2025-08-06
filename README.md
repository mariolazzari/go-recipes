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

## Structs

### Go structs

```go
package main

import (
 "fmt"
 "log"
 "time"
)

type Event struct {
 ID   string
 Time time.Time
}

type DoorEvent struct {
 Event
 Action string // open, close
}

type TemperatureEvent struct {
 Event
 Value float64
}

func NewDoorEvent(id string, time time.Time, action string) (*DoorEvent, error) {
 if id == "" {
  return nil, fmt.Errorf("empty id")
 }

 evt := DoorEvent{
  Event:  Event{id, time},
  Action: action,
 }
 return &evt, nil
}

func main() {
 evt, err := NewDoorEvent("front door", time.Now(), "open")
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("%+v\n", evt)
 // &{Event:{ID:front door Time:2021-04-30 14:47:40.31038 +0300 IDT m=+0.000170354} Action:open}
}
```

### Go methods

```go
package main

import (
 "fmt"
)

// A Thermostat measures and controls the temperature
type Thermostat struct {
 ID string

 value float64
}

// Value return the current temperature in Celsius
func (t *Thermostat) Value() float64 {
 return t.value
}

// Set tells the thermostat to set the temperature
func (t *Thermostat) Set(value float64) {
 t.value = value
}

// Kind returns the device kind
func (*Thermostat) Kind() string {
 return "thermostat"
}

func main() {
 t := Thermostat{"Living Room", 16.2}
 fmt.Printf("%s before: %.2f\n", t.ID, t.Value())
 // Living Room before: 16.20
 t.Set(18)
 fmt.Printf("%s after:  %.2f\n", t.ID, t.Value())
 // Living Room after:  18.00
}
```

### Go interfaces

```go
package main

import (
 "fmt"
)

// A Thermostat measures and controls the temperature
type Thermostat struct {
 id    string
 value float64
}

// ID return the thermostat ID
func (t *Thermostat) ID() string {
 return t.id
}

// Value return the current temperature in Celsius
func (t *Thermostat) Value() float64 {
 return t.value
}

// Kind returns the device kind
func (*Thermostat) Kind() string {
 return "thermostat"
}

// Camera is a security camera
type Camera struct {
 id string
}

// ID return the camera ID
func (c *Camera) ID() string {
 return c.id
}

func (*Camera) Kind() string {
 return "camera"
}

type Sensor interface {
 ID() string
 Kind() string
}

func printAll(sensors []Sensor) {
 for _, s := range sensors {
  fmt.Printf("%s <%s>\n", s.ID(), s.Kind())
 }
}

func main() {
 t := Thermostat{"Living Room", 16.2}
 c := Camera{"Baby room"}

 sensors := []Sensor{&t, &c}
 printAll(sensors)
 /*
  Living Room <thermostat>
  Baby room <camera>
 */
}
```

### Empty interface

```go
package main

import (
 "fmt"
 "log"
)

type ClickEvent struct {
 // ...
}

type HoverEvent struct {
 // ...
}

var eventCounts = make(map[string]int) // type -> count

func recordEvent(evt interface{}) {
 switch evt.(type) {
 case *ClickEvent:
  eventCounts["click"]++
 case *HoverEvent:
  eventCounts["hover"]++
 default:
  log.Printf("warning: unknown event: %#v of type %T\n", evt, evt)
 }
}

func main() {
 recordEvent(&ClickEvent{})
 recordEvent(&HoverEvent{})
 recordEvent(&ClickEvent{})
 recordEvent(3)
 // 2021/04/30 15:07:17 warning: unknown event: 3 of type int

 fmt.Println("event counts:", eventCounts)
 // event counts: map[click:2 hover:1]
}
```

### Iota

```go
package main

import (
 "fmt"
)

// LogLevel is a logging level
type LogLevel uint8

// Possible log levels
const (
 DebugLevel LogLevel = iota + 1
 WarningLevel
 ErrorLevel
)

// String implements the fmt.Stringer interface
func (l LogLevel) String() string {
 switch l {
 case DebugLevel:
  return "debug"
 case WarningLevel:
  return "warning"
 case ErrorLevel:
  return "error"
 }

 return fmt.Sprintf("unknown log level: %d", l)
}

func main() {
 fmt.Println(WarningLevel) // warning

 lvl := LogLevel(19)
 fmt.Println(lvl) // unknown log level: 19
}
```

### Generics

```go
package main

import (
 "fmt"
)

// LogLevel is a logging level
type LogLevel uint8

// Possible log levels
const (
 DebugLevel LogLevel = iota + 1
 WarningLevel
 ErrorLevel
)

// String implements the fmt.Stringer interface
func (l LogLevel) String() string {
 switch l {
 case DebugLevel:
  return "debug"
 case WarningLevel:
  return "warning"
 case ErrorLevel:
  return "error"
 }

 return fmt.Sprintf("unknown log level: %d", l)
}

func main() {
 fmt.Println(WarningLevel) // warning

 lvl := LogLevel(19)
 fmt.Println(lvl) // unknown log level: 19
}
```

### Challenge

```go
/*
Implement a paining program. It should support

- Circle with location (x, y), color and radius
- Rectangle with location (x, y), width, height and color

Each type should implement a `Draw(d Device)` method.

Implement an `ImageCanvas` struct which hold a slice of drawable items and has
`Draw(w io.Writer)` that writes a PNG to w (using `image/png`).
*/
package main

import (
 "fmt"
 "image"
 "image/color"
 "image/png"
 "io"
 "log"
 "math"
 "os"
)

var (
 Red   = color.RGBA{0xFF, 0, 0, 0xFF}
 Green = color.RGBA{0, 0xFF, 0, 0xFF}
 Blue  = color.RGBA{0, 0, 0xFF, 0xFF}
)

type Shape struct {
 X     int
 Y     int
 Color color.Color
}

type Circle struct {
 Shape
 Radius int
}

func NewCircle(x, y, r int, c color.Color) *Circle {
 cr := Circle{
  Shape:  Shape{x, y, c},
  Radius: r,
 }
 return &cr
}

func (c *Circle) Draw(d Device) {
 minX, minY := c.X-c.Radius, c.Y-c.Radius
 maxX, maxY := c.X+c.Radius, c.Y+c.Radius
 for x := minX; x <= maxX; x++ {
  for y := minY; y <= maxY; y++ {
   dx, dy := x-c.X, y-c.Y
   if int(math.Sqrt(float64(dx*dx+dy*dy))) <= c.Radius {
    d.Set(x, y, c.Color)
   }
  }
 }
}

type Rectangle struct {
 Shape
 Height int
 Width  int
}

func NewRectangle(x, y, h, w int, c color.Color) *Rectangle {
 r := Rectangle{
  Shape:  Shape{x, y, c},
  Height: h,
  Width:  w,
 }
 return &r
}

func (r *Rectangle) Draw(d Device) {
 minX, minY := r.X-r.Width/2, r.Y-r.Height/2
 maxX, maxY := r.X+r.Width/2, r.Y+r.Height/2
 for x := minX; x <= maxX; x++ {
  for y := minY; y <= maxY; y++ {
   d.Set(x, y, r.Color)
  }
 }
}

type Device interface {
 Set(int, int, color.Color)
}

type ImageCanvas struct {
 width  int
 height int
 shapes []Drawer
}

func NewImageCanvas(width, height int) (*ImageCanvas, error) {
 if width <= 0 || height <= 0 {
  return nil, fmt.Errorf("negative size: width=%d, height=%d", width, height)
 }

 c := ImageCanvas{
  width:  width,
  height: height,
 }
 return &c, nil
}

type Drawer interface {
 Draw(d Device)
}

func (ic *ImageCanvas) Add(d Drawer) {
 ic.shapes = append(ic.shapes, d)
}

func (ic *ImageCanvas) Draw(w io.Writer) error {
 img := image.NewRGBA(image.Rect(0, 0, ic.width, ic.height))
 for _, s := range ic.shapes {
  s.Draw(img)
 }
 return png.Encode(w, img)
}

func main() {
 ic, err := NewImageCanvas(200, 200)
 if err != nil {
  log.Fatal(err)
 }

 ic.Add(NewCircle(100, 100, 80, Green))
 ic.Add(NewCircle(60, 60, 10, Blue))
 ic.Add(NewCircle(140, 60, 10, Blue))
 ic.Add(NewRectangle(100, 130, 10, 80, Red))
 f, err := os.Create("face.png")
 if err != nil {
  log.Fatal(err)
 }
 defer f.Close()
 if err := ic.Draw(f); err != nil {
  log.Fatal(err)
 }
}
```

## JSON

### Unmarshalling JSON

```go
package main

import (
 "encoding/json"
 "fmt"
 "io"
 "log"
 "os"
 "time"
)

// Record is a weather record
type Record struct {
 Time    time.Time
 Station string
 Temp    float64 `json:"temperature"` // celsius
 Rain    float64 // millimeter
}

func readRecord(r io.Reader) (Record, error) {
 var rec Record
 dec := json.NewDecoder(r)
 if err := dec.Decode(&rec); err != nil {
  return Record{}, err
 }

 return rec, nil
}

func main() {
 file, err := os.Open("record.json")
 if err != nil {
  log.Fatal(err)
 }
 defer file.Close()

 rec, err := readRecord(file)
 if err != nil {
  log.Fatal(err)
 }

 fmt.Printf("%+v\n", rec)
 // {Time:2020-03-06 00:00:00 +0000 UTC Station:DS9 Temp:21.6 Rain:0}
}
```

### Parsing complex JSON

```go
package main

import (
 "encoding/json"
 "fmt"
 "io"
 "log"
 "os"
 "time"
)

// laggingStations return stations that are lagging in their check time
func laggingStations(r io.Reader, timeout time.Duration) ([]string, error) {
 var reply struct {
  LastCheckTime string
  Stations      []struct {
   Name      string
   Status    string
   LastCheck struct {
    Time string
   }
  }
 }

 dec := json.NewDecoder(r)
 if err := dec.Decode(&reply); err != nil {
  return nil, err
 }

 checkTime, err := parseTime(reply.LastCheckTime)
 if err != nil {
  return nil, err
 }

 var lagging []string
 for _, station := range reply.Stations {
  if station.Status != "Active" {
   continue
  }
  lastCheck, err := parseTime(station.LastCheck.Time)
  if err != nil {
   return nil, err
  }
  if checkTime.Sub(lastCheck) > timeout {
   lagging = append(lagging, station.Name)
  }
 }

 return lagging, nil
}

func parseTime(ts string) (time.Time, error) {
 return time.Parse("2006-01-02 15:04:05 PM", ts)
}

func main() {
 file, err := os.Open("stations.json")
 if err != nil {
  log.Fatal(err)
 }
 defer file.Close()

 lagging, err := laggingStations(file, time.Minute)
 if err != nil {
  log.Fatal(err)
 }

 for _, name := range lagging {
  fmt.Println(name)
 }
 // station 3
}
```

### Marshalling JSON

```go
package main

import (
 "encoding/json"
 "fmt"
 "os"
)

// Quantity is combination of value and unit (e.g. 2.7cm)
type Quantity struct {
 Value float64
 Unit  string
}

// MarshalJSON implements the json.Marshaler interface
// Example encoding: "42.195km"
func (q *Quantity) MarshalJSON() ([]byte, error) {
 if q.Unit == "" {
  return nil, fmt.Errorf("empty  unit")
 }
 text := fmt.Sprintf("%f%s", q.Value, q.Unit)
 return json.Marshal(text)
}

func main() {
 q := Quantity{1.78, "meter"}
 json.NewEncoder(os.Stdout).Encode(&q) // "1.780000meter"
}
```

### Missing values

```go
package main

import (
 "encoding/json"
 "fmt"
)

// LineItem is a line in receipt
type LineItem struct {
 SKU      string
 Price    float64
 Discount float64
 Quantity int
}

// NewLineItem returns a new line item with default values
func NewLineItem() LineItem {
 return LineItem{
  Quantity: 1,
 }
}

func unmarshalLineItem(data []byte) (LineItem, error) {
 li := NewLineItem()
 if err := json.Unmarshal(data, &li); err != nil {
  return LineItem{}, nil
 }

 if li.Quantity < 1 {
  return LineItem{}, fmt.Errorf("bad quantity")
 }

 return li, nil
}

func main() {
 data := []byte(`{"sku": "x3xs", "price": 1.2}`)
 li, err := unmarshalLineItem(data)
 if err != nil {
  fmt.Println("ERROR:", err)
 } else {
  fmt.Printf("%#v\n", li)
 }
 // main.LineItem{SKU:"x3xs", Price:1.2, Discount:0, Quantity:1}

 data = []byte(`{"sku": "x3xs", "price": 1.2, "quantity": 0}`)
 li, err = unmarshalLineItem(data)
 if err != nil {
  fmt.Println("ERROR:", err)
 } else {
  fmt.Printf("%#v\n", li)
 }
 // ERROR: bad quantity
}
```

### Map structure

```go
package main

import (
 "encoding/json"
 "fmt"

 "github.com/mitchellh/mapstructure"
)

// StartJob is a request to start a job
type StartJob struct {
 Type  string
 User  string
 Count int
}

// JobStatus is a request for job status
type JobStatus struct {
 Type string
 ID   string
}

func handleStart(req StartJob) error {
 fmt.Printf("start: %#v\n", req)
 return nil // FIXME
}

func handleStatus(req JobStatus) error {
 fmt.Printf("status: %#v\n", req)
 return nil // FIXME
}

func handleRequest(data []byte) error {
 var m map[string]interface{}
 if err := json.Unmarshal(data, &m); err != nil {
  return err
 }

 val, ok := m["type"]
 if !ok {
  return fmt.Errorf("'type' missing from JSON")
 }

 typ, ok := val.(string)
 if !ok {
  return fmt.Errorf("'type' is not a string")
 }

 switch typ {
 case "start":
  var sj StartJob
  if err := mapstructure.Decode(m, &sj); err != nil {
   return fmt.Errorf("bad 'start' request: %w", err)
  }
  return handleStart(sj)
 case "status":
  var js JobStatus
  if err := mapstructure.Decode(m, &js); err != nil {
   return fmt.Errorf("bad 'status' request: %w", err)
  }
  return handleStatus(js)
 }

 return fmt.Errorf("unknown request type: %q", typ)
}

func main() {
 data := []byte(`{"type": "start", "user": "joe", "count": 7}`)
 if err := handleRequest(data); err != nil {
  fmt.Println("ERROR:", err)
 }
 // start: main.StartJob{Type:"start", User:"joe", Count:7}

 data = []byte(`{"type": "status", "id": "seven"}`)
 if err := handleRequest(data); err != nil {
  fmt.Println("ERROR:", err)
 }
 // status: main.JobStatus{Type:"status", ID:"seven"}
}
```

### Challenge JSON

```go
// What is the maximal ride speed in rides.json?
package main

import (
 "encoding/json"
 "fmt"
 "io"
 "log"
 "os"
 "time"
)

func maxRideSpeed(r io.Reader) (float64, error) {
 dec := json.NewDecoder(r)
 maxSpeed := -1.0
 for {
  var ride struct {
   StartTime string `json:"start"`
   EndTime   string `json:"end"`
   Distance  float64
  }
  err := dec.Decode(&ride)
  if err == io.EOF {
   break
  }
  if err != nil {
   return 0, err
  }

  const timeFmt = "2006-01-02T15:04"
  startTime, err := time.Parse(timeFmt, ride.StartTime)
  if err != nil {
   return 0, err
  }
  endTime, err := time.Parse(timeFmt, ride.EndTime)
  if err != nil {
   return 0, err
  }
  dt := endTime.Sub(startTime)
  dtHour := float64(dt) / float64(time.Hour)
  speed := ride.Distance / dtHour
  if speed > maxSpeed {
   maxSpeed = speed
  }
 }

 return maxSpeed, nil
}

func main() {
 file, err := os.Open("rides.json")
 if err != nil {
  log.Fatal(err)
 }
 defer file.Close()

 speed, err := maxRideSpeed(file)
 if err != nil {
  log.Fatal(err)
 }
 fmt.Println(speed) // 40.5
}
```
