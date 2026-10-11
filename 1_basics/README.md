# Go Basics: Variables, Types, Control Flow and Loops

Oct 11, 2026 · @HR

## 1. Program structure

Every Go file belongs to a package, and an executable program starts at the `main` function of package `main`.

```go
package main // an executable must be in package main

import (
    "fmt"
    "os"
)

func main() { // the entry point: no arguments, no return value
    fmt.Println("Hello, Go")

    if len(os.Args) > 1 {
        fmt.Println("first argument:", os.Args[1])
    }
}
```

```
go run main.go      # compile and run in one step
go build -o app .   # produce a single static binary named app
```

- Go is compiled and statically typed, and it produces one self-contained binary with no runtime to install.
- An unused import or an unused local variable is a compile error, not a warning.
- The opening brace must sit on the same line as its statement, because Go inserts semicolons at line ends automatically.
- A name that starts with an uppercase letter is exported from its package. Lowercase names are private to it.
- `gofmt` formats all Go code the same way, so style is never debated.
- Go is garbage collected, so you never free memory by hand.

## 2. Variables, zero values and scope

A variable is declared with `var` or, inside a function, with the short form `:=`, and it always starts with a usable zero value.

```go
var a int           // 0: declared without a value
var b string = "go" // explicit type and value
var c = 3.14        // type inferred: float64
d := 10             // short declaration: only inside functions

var e, f = 1, "x" // several at once
x, y := 1, 2
x, y = y, x // swap without a temporary variable

var ( // grouped declaration, common at package level
    host = "localhost"
    port = 8080
)

_, err := strconv.Atoi("42") // _ discards a value you do not need
```

| Type | Zero value |
| --- | --- |
| Integers and floats | 0 |
| bool | false |
| string | "" (empty string) |
| Pointer, slice, map, channel, function, interface | nil |
| Struct | Every field at its own zero value |
| Array | Every element at its zero value |

### Scope and shadowing

```go
x := 1
if true {
    x := 2 // a NEW x that exists only inside this block
    x++
}
fmt.Println(x) // 1: the outer x was never touched

var err error
if true {
    n, err := strconv.Atoi("abc") // := also declares a new err here
    _ = n
    _ = err
}
fmt.Println(err) // <nil>: the failure was lost
```

- `:=` declares and assigns. `=` only assigns to a variable that already exists.
- `:=` needs at least one new variable on its left. Existing ones in the same scope are simply assigned.
- `:=` cannot be used at package level. Use `var` there.
- There are no uninitialized variables in Go, so reading a fresh variable is always safe.
- A variable declared in an `if`, `for` or `switch` header lives only inside that statement.
- Go has no `undefined` or `null` for basic types. An `int` is never nil.

## 3. Basic types and conversions

Go has a small set of built-in types and never converts between them implicitly, so every conversion is written out.

| Type | Details |
| --- | --- |
| `int`, `uint` | 64 bits on 64-bit systems; `int` is the default for whole numbers |
| `int8`, `int16`, `int32`, `int64` | Signed integers of a fixed size |
| `uint8`, `uint16`, `uint32`, `uint64` | Unsigned integers of a fixed size |
| `float32`, `float64` | Floating point; `float64` is the default for decimals |
| `bool` | `true` or `false`; never converts to or from a number |
| `string` | An immutable sequence of bytes |
| `byte` | Alias for `uint8` |
| `rune` | Alias for `int32`; one Unicode code point |

```go
var i int = 42
var f float64 = float64(i) // conversion must be explicit
var u uint = uint(f)

var a int32 = 1
var b int64 = 2
// fmt.Println(a + b)      // compile error: mismatched types int32 and int64
fmt.Println(int64(a) + b) // 3

s := strconv.Itoa(42)        // number to text: "42"
n, err := strconv.Atoi("42") // text to number: 42, nil

fmt.Println(7/2, 7%2, 7.0/2) // 3 1 3.5: integer division drops the fraction
```

### Overflow and floating point

```go
var m int8 = 127
m++
fmt.Println(m) // -128: integer overflow wraps around silently

var big uint8 = 255
big++
fmt.Println(big) // 0

x, y := 0.1, 0.2
fmt.Println(x+y == 0.3) // false: float64 cannot represent these exactly
fmt.Println(x + y)      // 0.30000000000000004
```

### Number literals

```go
million := 1_000_000 // underscores for readability
bin := 0b1010        // binary: 10
oct := 0o17          // octal: 15
hex := 0xFF          // hexadecimal: 255
ch := 'A'            // a rune: 65
```

- Integer overflow does not panic, so check ranges yourself when sizes matter.
- Never compare floats with `==`. Compare the difference against a small tolerance, and use integers or a decimal library for money.
- `string(65)` gives "A", not "65". Use `strconv.Itoa` to turn a number into digits.
- The size of `int` depends on the platform. Use `int64` or `int32` when a file format or protocol fixes the size.

## 4. Constants and iota

A constant is a value fixed at compile time, and `iota` generates a numbered sequence of constants, which is how Go writes enums.

```go
const Pi = 3.14159

const (
    StatusOK       = 200
    StatusNotFound = 404
)

// An enum: a defined type plus iota
type Weekday int

const (
    Sunday  Weekday = iota // 0
    Monday                 // 1
    Tuesday                // 2
)

// Skipping a value and using an expression
const (
    _  = iota             // 0 is discarded
    KB = 1 << (10 * iota) // 1024
    MB                    // 1048576
    GB                    // 1073741824
)

// Bit flags
const (
    Read    = 1 << iota // 1
    Write               // 2
    Execute             // 4
)

perm := Read | Write          // 3
fmt.Println(perm&Write != 0)   // true: write is set
fmt.Println(perm&Execute != 0) // false
```

### Typed and untyped constants

```go
const n = 5 // untyped: takes whatever numeric type the context needs

var i int = n
var f float64 = n // fine, no conversion

const typed int = 5
// var g float64 = typed // compile error: typed is an int
```

- `iota` starts at 0 in each `const` block and goes up by one per line.
- A line with no expression repeats the previous expression with the new `iota` value.
- Only numbers, strings, booleans and runes can be constants. A slice, map or struct cannot.
- A constant has no address, so `&Pi` does not compile.
- Go has no `enum` keyword. Give the type a `String()` method to print names instead of numbers.
- Start enums at 1, or reserve 0 for "unknown", when the zero value should not mean a real option.

## 5. Operators

Go's operators are the familiar C-style set, with three differences: `++` is a statement, there is no ternary operator, and both operands must have the same type.

```go
// Arithmetic
fmt.Println(7+2, 7-2, 7*2, 7/2, 7%2) // 9 5 14 3 1
fmt.Println(-7/2, -7%2)              // -3 -1: division truncates toward zero
fmt.Println("go" + "lang")           // golang: + also joins strings

// Comparison: the result is always a bool
fmt.Println(3 == 3, 3 != 4, 3 < 4, 3 >= 4) // true true true false

// Logical: short-circuit evaluation
if user != nil && user.Active { // the right side is skipped when user is nil
}

// Bitwise: 6 is 110 and 3 is 011 in binary
fmt.Println(6&3, 6|3, 6^3) // 2 7 5: AND, OR, XOR
fmt.Println(1<<3, 16>>2)   // 8 4: shift left, shift right
fmt.Println(5 &^ 3)        // 4: AND NOT clears the bits set in 3

// Assignment forms
x := 10
x += 5 // 15
x++    // 16: a statement, not an expression
```

### Things that do not exist

```go
// y := x++          // compile error: ++ produces no value
// ++x               // compile error: there is no prefix form
// max := a > b ? a : b // compile error: no ternary operator

max := b
if a > b {
    max = a
}
```

- An integer divided by zero panics at runtime. A float divided by zero gives infinity.
- Go has no operator overloading and no implicit numeric promotion.
- Precedence from highest to lowest: multiplicative and shifts, additive, comparison, `&&`, then `||`. Add parentheses when in doubt.
- `n%2 == 0` tests for even numbers. `n&1 == 0` does the same with a bit test.
- Since Go 1.21 the built-in functions `min` and `max` cover the most common use of a ternary.

## 6. if and switch

Conditions need no parentheses but always need braces, and a `switch` case ends by itself without a `break`.

### if

```go
if score >= 90 {
    fmt.Println("A")
} else if score >= 75 {
    fmt.Println("B")
} else {
    fmt.Println("C")
}

// An init statement: the variable lives only inside the if-else chain
if n := len(name); n == 0 {
    fmt.Println("empty")
} else if n > 50 {
    fmt.Println("too long")
}

if v, ok := cache[key]; ok {
    fmt.Println("hit", v)
}

if err := save(user); err != nil {
    return err
}
```

### switch

```go
// Several values in one case, and no break needed
switch day {
case "sat", "sun":
    fmt.Println("weekend")
case "mon":
    fmt.Println("week starts")
default:
    fmt.Println("weekday")
}

// With no tag, a switch is a cleaner if-else chain
switch {
case score >= 90:
    grade = "A"
case score >= 75:
    grade = "B"
default:
    grade = "C"
}

// fallthrough runs the next case without testing its condition
switch n := 1; n {
case 1:
    fmt.Println("one")
    fallthrough
case 2:
    fmt.Println("two") // also printed
case 3:
    fmt.Println("three") // not printed
}

// A type switch branches on the dynamic type of an interface value
switch x := v.(type) {
case int:
    fmt.Println("int", x)
case string:
    fmt.Println("string", x)
}
```

- The condition of an `if` must be a `bool`. A number or a pointer is never treated as true or false.
- Cases are tried top to bottom, and only the first match runs.
- Case values do not have to be constants. Any expression of the right type works.
- `fallthrough` is rare. Listing several values in one case is the usual way to share a body.
- Prefer the early-return style: handle the error or edge case and return, so the main path is not nested in `else`.

## 7. Loops

Go has a single loop keyword, `for`, and it covers the counting loop, the while loop, the infinite loop and iteration with `range`.

### The four basic forms

```go
// 1. Classic three-part loop: init; condition; post
for i := 0; i < 3; i++ {
    fmt.Print(i, " ") // 0 1 2
}

// 2. Condition only: Go's while loop
n := 5
for n > 0 {
    n--
}

// 3. No condition: an infinite loop, left with break or return
for {
    line, err := reader.ReadString('\n')
    if err != nil {
        break
    }
    process(line)
}

// 4. range: iterate over a collection
for i, v := range []string{"a", "b"} {
    fmt.Println(i, v) // 0 a, then 1 b
}
```

### What range yields

| Ranging over | First value | Second value |
| --- | --- | --- |
| Slice or array | Index | A copy of the element |
| String | Byte index | The rune starting there |
| Map | Key | A copy of the value; order is random |
| Channel | The received value | None; the loop ends when the channel is closed |
| Integer n (Go 1.22) | 0 up to n-1 | None |
| Iterator function (Go 1.23) | Whatever it yields | Whatever it yields |

```go
nums := []int{10, 20, 30}

for i, v := range nums { } // index and value
for i := range nums { }    // index only
for _, v := range nums { } // value only
for range nums { }         // just repeat len(nums) times

for k, v := range map[string]int{"a": 1, "b": 2} {
    fmt.Println(k, v) // order changes from run to run
}

for i, r := range "héy" {
    fmt.Print(i, ":", string(r), " ") // 0:h 1:é 3:y, because é takes two bytes
}

for v := range ch { // receives until ch is closed
    fmt.Println(v)
}

for i := range 3 { // Go 1.22: range over an integer
    fmt.Print(i, " ") // 0 1 2
}
```

### More shapes

```go
// Two variables: the usual two-pointer loop
for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    s[i], s[j] = s[j], s[i] // reverses s in place
}

// Counting down
for i := len(nums) - 1; i >= 0; i-- {
    fmt.Println(nums[i])
}

// Do-while: run the body at least once
for {
    attempt++
    if tryConnect() || attempt >= 3 {
        break
    }
}

// Nested loops
for row := 0; row < 3; row++ {
    for col := 0; col < 3; col++ {
        fmt.Print(row*3+col, " ")
    }
    fmt.Println()
}

// Range over a function (Go 1.23): a custom iterator
func Countdown(n int) iter.Seq[int] {
    return func(yield func(int) bool) {
        for i := n; i > 0; i-- {
            if !yield(i) {
                return // the caller used break
            }
        }
    }
}

for v := range Countdown(3) {
    fmt.Print(v, " ") // 3 2 1
}
```

- There is no `while`, `do-while` or `foreach` keyword. All of them are spelled `for`.
- Any of the three parts of a classic loop may be left out; `for ;; {}` is the same as `for {}`.
- Variables declared in the loop header are visible only inside the loop.
- Since Go 1.22, each iteration gets its own copy of the loop variable.
- Ranging over a nil slice or nil map is safe. The body simply runs zero times.

## 8. Loop control: break, continue, labels and goto

`break` leaves the innermost loop, `continue` skips to its next iteration, and a label lets either one target an outer loop.

```go
for i := 0; i < 10; i++ {
    if i%2 == 0 {
        continue // skip even numbers
    }
    if i > 7 {
        break // stop the loop entirely
    }
    fmt.Print(i, " ") // 1 3 5 7
}
```

### Labels for nested loops

```go
outer:
    for i := 0; i < 3; i++ {
        for j := 0; j < 3; j++ {
            if j == 2 {
                continue outer // go to the next i
            }
            if i == 2 {
                break outer // leave both loops
            }
            fmt.Println(i, j)
        }
    }
// 0 0
// 0 1
// 1 0
// 1 1
```

### break inside a switch or select

```go
// Surprise: this break leaves only the switch
for _, v := range []int{1, 2, 3} {
    switch v {
    case 2:
        break
    }
    fmt.Print(v, " ")
}
// 1 2 3

// With a label, break leaves the loop
loop:
    for _, v := range []int{1, 2, 3} {
        switch v {
        case 2:
            break loop
        }
        fmt.Print(v, " ")
    }
// 1
```

The same applies to `select` inside a `for`, which is the usual shape of a goroutine's main loop.

### Searching with a flag or a function

```go
// Often cleaner than a label: put the loops in a function and return
func find(grid [][]int, target int) (int, int, bool) {
    for r, row := range grid {
        for c, v := range row {
            if v == target {
                return r, c, true
            }
        }
    }
    return 0, 0, false
}
```

- A label names the statement that follows it, and `break` or `continue` with that label acts on that loop.
- `continue` in a classic loop still runs the post statement, so `i++` happens.
- `goto` exists but is almost never used. It cannot jump over a variable declaration or into a block.
- An unused label is a compile error, like an unused variable.

## 9. Loop gotchas

These six behaviours cause most loop bugs and most "what does this print?" questions.

### 1. The range value is a copy

```go
type User struct{ Age int }

users := []User{{30}, {25}}

for _, u := range users {
    u.Age++ // changes a copy; users is untouched
}
fmt.Println(users) // [{30} {25}]

for i := range users {
    users[i].Age++ // index into the slice to change the real element
}
fmt.Println(users) // [{31} {26}]
```

### 2. range evaluates its operand once

```go
s := []int{1, 2, 3}
for range s {
    s = append(s, 0) // the loop still runs exactly 3 times
}
fmt.Println(len(s)) // 6

arr := [3]int{1, 2, 3}
for i, v := range arr { // ranging over an ARRAY iterates over a copy of it
    arr[2] = 100
    if i == 2 {
        fmt.Println(v) // 3, not 100
    }
}
```

### 3. Closures and the loop variable

```go
var fns []func()
for i := 0; i < 3; i++ {
    fns = append(fns, func() { fmt.Print(i, " ") })
}
for _, f := range fns {
    f()
}
// Go 1.22+: 0 1 2
// Before 1.22: 3 3 3, because every closure shared one i
```

### 4. An unsigned counter never goes below zero

```go
var i uint
for i = 3; i >= 0; i-- { // always true: after 0 comes the largest uint
    fmt.Println(i)       // 3 2 1 0 18446744073709551615 ... forever
}
```

### 5. Removing elements while iterating

```go
nums := []int{1, 2, 2, 3}

// Wrong: after a removal the next element slides into index i and is skipped
for i := 0; i < len(nums); i++ {
    if nums[i] == 2 {
        nums = append(nums[:i], nums[i+1:]...)
    }
}
fmt.Println(nums) // [1 2 3]: one 2 survived

// Right: build a filtered slice, or use slices.DeleteFunc
nums = slices.DeleteFunc(nums, func(n int) bool { return n == 2 })
```

### 6. defer inside a loop

```go
for _, path := range paths {
    f, _ := os.Open(path)
    defer f.Close() // nothing closes until the whole function returns
}
```

- Map iteration order is random. Sort the keys when the order matters.
- Deleting from a map during `range` is safe. Adding to it may or may not show up in the same loop.
- Use integer loop counters. A float counter accumulates rounding error and may run one time too many or too few.
- Ranging over a string yields runes; an index loop over the same string yields bytes.

## 10. Defined types, new vs make, and pass by value

Three small rules explain much of Go's behaviour: a defined type is distinct from its underlying type, `new` and `make` do different jobs, and every argument is passed by value.

### Defined type vs alias

```go
type Celsius float64 // a new, distinct type that can have its own methods
type Temp = float64  // an alias: just a second name for float64

var c Celsius = 36.6
var f float64 = float64(c) // conversion required between distinct types

var t Temp = 36.6
f = t // fine: Temp IS float64
```

Defined types catch mistakes at compile time: a `UserID` cannot be passed where an `OrderID` is expected, even if both are integers underneath.

### new vs make

|  | `new(T)` | `make(T, ...)` |
| --- | --- | --- |
| Works with | Any type | Only slices, maps and channels |
| Returns | A pointer, `*T` | The value itself, `T` |
| State | Zeroed memory | Initialized and ready to use |

```go
p := new(int)            // *int pointing at 0
s := make([]int, 0, 10)  // empty slice with room for 10
m := make(map[string]int) // a map you can write to
ch := make(chan int, 5)  // buffered channel

var nilMap map[string]int
// nilMap["a"] = 1 // panic: a map must be created with make or a literal first
```

### Everything is passed by value

```go
func changeInt(n int)        { n = 99 }       // changes a copy
func changePtr(n *int)       { *n = 99 }      // changes the original through the pointer
func changeSlice(s []int)    { s[0] = 99 }    // the copied header points at the same array
func changeMap(m map[string]int) { m["k"] = 99 } // the copied pointer refers to the same table

n := 1
changeInt(n)
fmt.Println(n) // 1
changePtr(&n)
fmt.Println(n) // 99

s := []int{1, 2}
changeSlice(s)
fmt.Println(s) // [99 2]
```

- Go has no pass by reference. What looks like it is a copied pointer or a copied slice header.
- Arrays and structs are copied in full when passed. Pass a pointer to avoid the copy or to allow changes.
- Slices, maps and channels are cheap to pass, because only a small header or pointer is copied.
- `&T{}` is the usual way to get a pointer to a struct; `new` is mostly seen with basic types.

## 11. Printing and reading input

The `fmt` package prints with format verbs, and `bufio.Scanner` is the reliable way to read input in a coding round.

| Function | Writes to | Notes |
| --- | --- | --- |
| `fmt.Println` | Standard output | Adds spaces between values and a newline |
| `fmt.Printf` | Standard output | Uses a format string; no newline unless you add one |
| `fmt.Sprintf` | A returned string | Same formatting, nothing printed |
| `fmt.Fprintf` | Any `io.Writer` | Files, HTTP responses, buffers |
| `fmt.Errorf` | A returned error | Formatting plus error wrapping |

```go
type Point struct{ X, Y int }
p := Point{1, 2}

fmt.Printf("%v %+v %#v %T\n", p, p, p, p)
// {1 2} {X:1 Y:2} main.Point{X:1, Y:2} main.Point

fmt.Printf("%d %05d %x %b\n", 42, 42, 255, 5) // 42 00042 ff 101
fmt.Printf("%.2f|%8.3f|\n", 3.14159, 3.14159) // 3.14|   3.142|
fmt.Printf("%s %q %t\n", "go", "go", true)    // go "go" true
fmt.Printf("|%5d|%-5d|\n", 42, 42)            // |   42|42   |
fmt.Printf("%c %U\n", 'A', 'A')               // A U+0041
```

| Verb | Meaning |
| --- | --- |
| `%v` | Default format for any value |
| `%+v` | Struct with field names |
| `%#v` | Go syntax, useful when debugging |
| `%T` | The type of the value |
| `%d`, `%b`, `%x` | Integer in decimal, binary, hexadecimal |
| `%f`, `%.2f` | Float, with a fixed number of decimals |
| `%s`, `%q` | String, plain or quoted |
| `%t` | Boolean |
| `%c`, `%U` | Rune as a character, or as a Unicode code point |
| `%p` | Pointer address |

### Reading input

```go
// Line by line: the safest choice
scanner := bufio.NewScanner(os.Stdin)
for scanner.Scan() {
    line := scanner.Text()
    fmt.Println("read:", line)
}

// Whitespace-separated values
var n int
var name string
fmt.Scan(&n, &name) // needs pointers, so it can store into the variables
```

- A width goes before the verb: `%5d` right-aligns in five columns and `%-5d` left-aligns.
- `%v` works for everything, so reach for it when unsure.
- `fmt.Println(a, b)` puts a space between its arguments. `fmt.Print` does so only between non-string values.
- `go vet` checks that format verbs match their arguments.

## 12. Packages, modules and the go tool

A package is a directory of Go files compiled together, and a module is a versioned collection of packages described by a `go.mod` file.

```
myapp/
  go.mod            module github.com/asha/myapp, plus the Go version and dependencies
  go.sum            checksums of every dependency
  main.go           package main
  internal/
    store/
      store.go      package store: importable only from inside myapp
  pkg/
    mathx/
      mathx.go      package mathx
```

```go
// pkg/mathx/mathx.go
package mathx

func Add(a, b int) int { return a + b } // exported: starts with an uppercase letter
func helper() {}                         // unexported: visible only inside mathx

// main.go
package main

import (
    "fmt"

    "github.com/asha/myapp/pkg/mathx"
)

func main() {
    fmt.Println(mathx.Add(2, 3)) // call it through the package name
}
```

| Command | Purpose |
| --- | --- |
| `go mod init <module path>` | Start a new module and create `go.mod` |
| `go run .` | Compile and run the package in the current directory |
| `go build ./...` | Compile everything |
| `go test ./...` | Run all tests |
| `go get <package>` | Add or upgrade a dependency |
| `go mod tidy` | Add missing dependencies and remove unused ones |
| `go fmt ./...` | Format the code |
| `go vet ./...` | Report suspicious constructs |
| `go doc fmt.Println` | Show documentation in the terminal |

- All files in one directory must declare the same package name.
- Visibility has only two levels: exported (uppercase) and package-private (lowercase). There is no `protected`.
- Import cycles are forbidden. If package A imports B, then B cannot import A.
- The `internal` directory is enforced by the compiler: code outside the parent tree cannot import it.
- Each package can have `init()` functions, which run once before `main`, after package-level variables are set.
- `go.mod` records the minimum version of each dependency, and `go.sum` guarantees the downloaded code has not changed.

## 13. Predict the output

Cover the answer comments and work out each of these five before reading on.

### Puzzle 1: shadowing

```go
func main() {
    x := 1
    {
        x := 2
        x += 10
    }
    if x := 5; x > 3 {
        x++
    }
    fmt.Println(x)
}
// 1: both inner x variables are separate from the outer one
```

### Puzzle 2: constants vs variables

```go
func main() {
    const a, b = 0.1, 0.2
    fmt.Println(a+b == 0.3)

    x, y := 0.1, 0.2
    fmt.Println(x+y == 0.3)
}
// true
// false
```

Constant expressions are evaluated exactly at compile time. Variables are `float64`, which cannot hold 0.1 or 0.2 exactly.

### Puzzle 3: iota with gaps

```go
const (
    A = iota * 10
    B
    C = "x"
    D
    E = iota
)

func main() {
    fmt.Println(A, B, C, D, E)
}
// 0 10 x x 4
```

`iota` counts lines in the block whether or not it is used, so `E` on the fifth line is 4. `D` repeats the expression above it.

### Puzzle 4: overflow and integer division

```go
func main() {
    var a uint8 = 250
    a += 10
    fmt.Println(a)

    fmt.Println(5/2, 5.0/2, -5/2, 5%-2)
}
// 4
// 2 2.5 -2 1
```

260 does not fit in a `uint8`, so it wraps to 260 minus 256. Integer division truncates toward zero, and the remainder takes the sign of the left operand.

### Puzzle 5: ranging over an array vs a slice

```go
func main() {
    arr := [3]int{1, 2, 3}
    for i, v := range arr {
        if i == 0 {
            arr[2] = 100
        }
        if i == 2 {
            fmt.Println("array:", v)
        }
    }

    sl := []int{1, 2, 3}
    for i, v := range sl {
        if i == 0 {
            sl[2] = 100
        }
        if i == 2 {
            fmt.Println("slice:", v)
        }
    }
}
// array: 3
// slice: 100
```

`range` copies its operand once. Copying an array copies all its elements; copying a slice copies only the header, which still points at the same data.

## 14. Interview coding problems

These five warm-up problems test loops, conditions and arithmetic, and they are often the first question in a coding round.

### FizzBuzz

```go
func fizzBuzz(n int) {
    for i := 1; i <= n; i++ {
        switch {
        case i%15 == 0:
            fmt.Println("FizzBuzz")
        case i%3 == 0:
            fmt.Println("Fizz")
        case i%5 == 0:
            fmt.Println("Buzz")
        default:
            fmt.Println(i)
        }
    }
}
```

### Check whether a number is prime

```go
func isPrime(n int) bool {
    if n < 2 {
        return false
    }
    for i := 2; i*i <= n; i++ { // testing up to the square root is enough
        if n%i == 0 {
            return false
        }
    }
    return true
}
// isPrime(17) is true; isPrime(18) is false
```

### Fibonacci, iteratively

```go
func fib(n int) int {
    a, b := 0, 1
    for i := 0; i < n; i++ {
        a, b = b, a+b // both right-hand values are computed before assigning
    }
    return a
}
// fib(0) through fib(7): 0 1 1 2 3 5 8 13
```

### Reverse an integer and test for a palindrome

```go
func reverseInt(n int) int {
    rev := 0
    for n != 0 {
        rev = rev*10 + n%10 // take the last digit
        n /= 10             // drop the last digit
    }
    return rev
}

func isPalindrome(n int) bool {
    return n >= 0 && n == reverseInt(n)
}
// reverseInt(1234) is 4321; reverseInt(-120) is -21; isPalindrome(1221) is true
```

### Largest and second largest in one pass

```go
func twoLargest(nums []int) (first, second int) {
    first, second = math.MinInt, math.MinInt
    for _, n := range nums {
        switch {
        case n > first:
            first, second = n, first // the old largest becomes second
        case n > second && n != first:
            second = n
        }
    }
    return first, second
}
// twoLargest([]int{4, 9, 2, 9, 7}) returns 9, 7
```

- State the complexity: every solution here is one pass, O(n), except the prime test, which is O(√n).
- Say the edge cases out loud: zero, negative numbers, an empty slice, and overflow for large inputs.
- For Fibonacci, mention that `int` overflows past the 92nd number, and that `math/big` handles larger values.
- Parallel assignment, as in `a, b = b, a+b`, is the idiomatic way to swap or rotate values in Go.

## 15. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| Is Go compiled or interpreted? | Compiled, to a single static binary. |
| `var` vs `:=`? | `var` works anywhere; `:=` declares with inference, only inside functions. |
| What is a zero value? | The default every variable gets: 0, false, an empty string, or nil. |
| Does Go convert types implicitly? | No; every conversion is explicit. |
| What is the default type of `x := 5` and `y := 5.0`? | `int` and `float64`. |
| byte vs rune? | `uint8`, one byte; `int32`, one Unicode code point. |
| What happens on integer overflow? | The value wraps around silently. |
| What is `iota`? | A counter in a `const` block that starts at 0 and rises by one per line. |
| Does Go have enums? | Not as a keyword; use a defined type with `iota` constants. |
| Typed vs untyped constant? | An untyped constant adapts to the context; a typed one does not. |
| Is there a ternary operator? | No; use `if`, or `min` and `max`. |
| Is `i++` an expression? | No, it is a statement, and there is no `++i`. |
| Does `switch` fall through? | No; each case ends by itself unless you write `fallthrough`. |
| How many loop keywords does Go have? | One: `for`. |
| How do you write a while loop? | `for condition { }`. |
| What does `range` return for a slice? | The index and a copy of the element. |
| Is map iteration ordered? | No; the order is deliberately random. |
| How do you break out of nested loops? | A labelled `break`, or return from a function. |
| What does `break` do inside a `switch` in a loop? | It leaves only the `switch`. |
| What changed for loop variables in Go 1.22? | Each iteration gets its own variable. |
| `new` vs `make`? | `new` returns a zeroed pointer; `make` initializes slices, maps and channels. |
| Pass by value or by reference? | Always by value. |
| How is visibility controlled? | Uppercase names are exported; lowercase names are package-private. |
| What does an unused variable or import cause? | A compile error. |
| What is `go.mod`? | The file that names a module and lists its dependencies. |
