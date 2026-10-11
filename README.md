# Go Interview Prep: Basics to Advanced

Oct 10, 2026 · @HR

## How to use this guide

Each section gives the concept, a short example, and the interview points that follow from it. Examples assume Go 1.22 or newer, and notes mark where older versions behave differently. For a last-minute revision, skip to the rapid-fire table at the end.

## 1. Basics

Go is statically typed, compiled and garbage collected, and every variable starts with a usable zero value.

### Declarations and zero values

```go
var a int            // 0
var s string         // ""
var p *int           // nil
var sl []int         // nil slice, len 0, safe to append
var m map[string]int // nil map: reads work, writes panic
b := 10              // short declaration, only inside functions
```

| Type | Zero value |
| --- | --- |
| Numeric types | 0 |
| bool | false |
| string | "" (empty string) |
| Pointer, slice, map, channel, function, interface | nil |
| Struct | Every field set to its own zero value |

### Constants and iota

```go
type Status int

const (
    Pending Status = iota // 0
    Active                // 1
    Closed                // 2
)

const (
    _  = iota             // skip 0
    KB = 1 << (10 * iota) // 1024
    MB                    // 1048576
)
```

### Variable shadowing

```go
var err error
if true {
    _, err := strconv.Atoi("x") // := creates a NEW err in this scope
    _ = err
}
fmt.Println(err) // <nil>: the outer err was never set
```

### Interview points

- Everything is passed by value. Slices, maps and channels copy a small header or pointer, so they behave like references.
- `new(T)` returns a zeroed `*T`. `make` works only for slices, maps and channels and returns an initialized value.
- `:=` needs at least one new variable on the left. In a nested scope it shadows instead of assigning.
- There is no implicit conversion: adding an `int32` to an `int64` does not compile.
- Constants are untyped until used, so `const x = 5` works as an int or a float64.
- Names starting with an uppercase letter are exported from the package; lowercase names are private to it.
- `init()` runs automatically before `main`, after package-level variables are initialized.

## 2. Arrays, slices and maps

A slice is a three-word header (pointer, length, capacity) over a backing array, and most slice questions test whether you know that.

### Arrays vs slices

```go
arr := [3]int{1, 2, 3} // fixed size, the size is part of the type
sl := []int{1, 2, 3}   // header pointing at a backing array

arr2 := arr // full copy
sl2 := sl   // same backing array
sl2[0] = 99
fmt.Println(arr2[0], sl[0]) // 1 99
```

### Shared backing array

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]        // [2 3], len 2, cap 4
b = append(b, 100) // fits in cap, so it overwrites a[3]
fmt.Println(a)     // [1 2 3 100 5]

c := a[1:3:3]      // full slice expression: cap limited to 2
c = append(c, 200) // exceeds cap, so a new array is allocated
fmt.Println(a)     // [1 2 3 100 5], unchanged this time
```

### Append inside a function

```go
func add(s []int) { s = append(s, 4) }

func main() {
    s := []int{1, 2, 3}
    add(s)
    fmt.Println(s) // [1 2 3]: the caller's header still has len 3
}
```

The fix is to return the new slice and assign it: `s = add(s)`.

### nil slice vs empty slice

```go
var a []int  // nil: a == nil is true
b := []int{} // empty but not nil
// Both have len 0 and work with append and range.
// json.Marshal gives null for a and [] for b.
```

### Maps

```go
m := map[string]int{"a": 1}
v, ok := m["zzz"] // 0 false: comma-ok separates missing from zero
delete(m, "a")    // no-op if the key is absent

var nm map[string]int
_ = nm["x"]  // fine, returns 0
nm["x"] = 1  // panic: assignment to entry in nil map

set := map[string]struct{}{} // set idiom: struct{} takes zero bytes
set["go"] = struct{}{}
```

### Interview points

- Slice capacity roughly doubles while small, then grows by about 1.25x. Preallocate with `make([]T, 0, n)` when you know the size.
- A small slice of a huge array keeps the whole array alive. Copy the part you need with `slices.Clone`.
- Map iteration order is random by design. Sort the keys when you need a stable order.
- Maps are not safe for concurrent use. Concurrent writes cause a fatal error that `recover` cannot catch.
- Map keys must be comparable: no slices, maps or functions. Structs and arrays of comparable types are fine.
- Map values are not addressable, so `m["k"].field = 1` does not compile. Store pointers or reassign the whole value.
- Deleting entries while ranging over a map is safe.

## 3. Strings, bytes and runes

A string is an immutable sequence of bytes, usually UTF-8: `len` counts bytes and `range` yields runes.

```go
s := "héllo"
fmt.Println(len(s))                    // 6 bytes: é takes 2
fmt.Println(utf8.RuneCountInString(s)) // 5 characters
fmt.Println(s[1])                      // 195: a byte, not 'é'

for i, r := range s { // i is the byte index, r is a rune
    fmt.Printf("%d:%c ", i, r) // 0:h 1:é 3:l 4:l 5:o
}
```

### Reverse a string correctly

```go
func reverse(s string) string {
    r := []rune(s) // runes, so multi-byte characters survive
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}
```

### Efficient concatenation

```go
var sb strings.Builder
for i := 0; i < 1000; i++ {
    sb.WriteString("x")
}
out := sb.String()
```

### Interview points

- `byte` is an alias for `uint8`; `rune` is an alias for `int32` and holds one Unicode code point.
- Strings are immutable, so `s[0] = 'x'` does not compile. Converting to `[]byte` or `[]rune` makes a copy.
- `+=` in a loop allocates a new string each time, which is O(n²). Use `strings.Builder`.
- A substring shares memory with the original string; no bytes are copied.
- `string(65)` gives "A", not "65". Use `strconv.Itoa` to convert a number to text.

## 4. Structs, pointers, methods and embedding

Go has no classes or inheritance: you attach methods to types and reuse code through embedding.

### Value vs pointer receivers

```go
type Counter struct{ n int }

func (c Counter) IncValue()    { c.n++ } // works on a copy
func (c *Counter) IncPointer() { c.n++ } // modifies the original

func main() {
    c := Counter{}
    c.IncValue()
    c.IncPointer()   // Go rewrites this as (&c).IncPointer()
    fmt.Println(c.n) // 1
}
```

Use a pointer receiver when the method mutates the value, the struct is large, or it contains a mutex. If one method needs a pointer receiver, use pointer receivers for all of them.

### Method sets

```go
type Shape interface{ Area() float64 }

type Square struct{ side float64 }

func (s *Square) Area() float64 { return s.side * s.side }

var _ Shape = &Square{2} // ok
var _ Shape = Square{2}  // compile error: Area has a pointer receiver
```

The rule: the method set of `T` holds only value-receiver methods, while the method set of `*T` holds both kinds.

### Embedding

```go
type Logger struct{ prefix string }

func (l Logger) Log(msg string) { fmt.Println(l.prefix, msg) }

type Service struct {
    Logger // embedded: Log is promoted to Service
    name   string
}

s := Service{Logger: Logger{"[svc]"}, name: "orders"}
s.Log("started") // same as s.Logger.Log("started")
```

Embedding is composition, not inheritance. A promoted method still receives the inner `Logger`, never the outer `Service`, so there is no method overriding in the class sense.

### Interview points

- Structs are comparable with `==` only if every field is comparable.
- Returning a pointer to a local variable is safe; the compiler moves that variable to the heap.
- Struct tags such as `json:"name,omitempty"` are metadata read through reflection.
- `encoding/json` skips unexported (lowercase) fields.
- Go has pointers but no pointer arithmetic outside the `unsafe` package.
- Constructors are plain functions by convention: `func NewService(...) *Service`.

## 5. Functions, closures, defer, panic and recover

Functions are first-class values, and `defer` runs calls in last-in-first-out order when the surrounding function returns.

### Defer order and argument evaluation

```go
func order() {
    for i := 0; i < 3; i++ {
        defer fmt.Print(i, " ") // prints 2 1 0
    }
}

func args() {
    x := 1
    defer fmt.Println("deferred:", x) // x is evaluated now: prints 1
    x = 2
}
```

### Defer with named return values

```go
func f() (result int) {
    defer func() { result *= 2 }()
    return 5 // sets result to 5, then the defer runs: returns 10
}
```

### Recover from a panic

```go
func safeDivide(a, b int) (res int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r)
        }
    }()
    return a / b, nil // panics when b is 0
}
```

### Closures

```go
func counter() func() int {
    n := 0
    return func() int { n++; return n }
}

next := counter()
next() // 1
next() // 2: n lives on because the closure captured it
```

### Loop variable capture

```go
var fns []func()
for i := 0; i < 3; i++ {
    fns = append(fns, func() { fmt.Print(i) })
}
for _, f := range fns {
    f()
}
// Go 1.22+: 012 (a new i for each iteration)
// Before 1.22: 333 (one i shared by all closures)
```

### Interview points

- `recover` works only when called directly inside a deferred function, in the same goroutine as the panic.
- An unrecovered panic in any goroutine crashes the whole program. One goroutine cannot recover another's panic.
- A `defer` inside a loop holds resources until the function returns. Wrap the loop body in a function instead.
- `os.Exit` and `log.Fatal` skip deferred calls.
- Variadic functions take `nums ...int`; pass an existing slice with `sum(nums...)`.
- Use `panic` for programmer bugs and impossible states. Return an `error` for anything a caller can expect.

## 6. Interfaces

Interfaces are satisfied implicitly: a type implements one simply by having its methods, with no `implements` keyword.

```go
type Notifier interface{ Notify(msg string) error }

type Email struct{ to string }

func (e Email) Notify(msg string) error { fmt.Println("email", e.to, msg); return nil }

type SMS struct{ phone string }

func (s SMS) Notify(msg string) error { fmt.Println("sms", s.phone, msg); return nil }

func alert(n Notifier) { _ = n.Notify("server down") } // accepts either type

var _ Notifier = Email{} // compile-time check that Email implements Notifier
```

### Type assertion and type switch

```go
var v any = "hello"

s, ok := v.(string) // safe form: ok is false on a mismatch
_ = v.(int)         // unsafe form: panics, because v holds a string

switch x := v.(type) {
case string:
    fmt.Println("string of length", len(x))
case int:
    fmt.Println("int", x)
default:
    fmt.Println("unknown type")
}
```

### The nil interface trap

```go
type MyErr struct{}

func (*MyErr) Error() string { return "boom" }

func do() error {
    var e *MyErr = nil
    return e // the interface now holds (type=*MyErr, value=nil)
}

func main() {
    err := do()
    fmt.Println(err == nil) // false
}
```

An interface value is nil only when both its type and its value are nil. The fix is to return a literal `nil` on the success path.

### Interview points

- An interface value is two words: a pointer to type information and a pointer to the data.
- Keep interfaces small. `io.Reader` and `io.Writer` have one method each.
- Accept interfaces, return structs. Define the interface in the package that uses it, not the one that implements it.
- `any` is an alias for `interface{}` since Go 1.18.
- Interfaces can embed other interfaces: `io.ReadWriter` combines `Reader` and `Writer`.
- Comparing two interface values that hold a non-comparable type, such as a slice, panics at runtime.
- Interfaces are how you mock dependencies in tests.

## 7. Error handling

Errors are ordinary values returned last: wrap them with `%w` to add context, and inspect them with `errors.Is` and `errors.As`.

```go
var ErrNotFound = errors.New("not found") // sentinel error

type ValidationError struct{ Field string } // custom error type

func (e *ValidationError) Error() string { return "invalid " + e.Field }

func findUser(id int) error {
    if id <= 0 {
        return &ValidationError{Field: "id"}
    }
    return fmt.Errorf("findUser %d: %w", id, ErrNotFound) // wrap with context
}

func main() {
    err := findUser(42)
    fmt.Println(err) // findUser 42: not found

    if errors.Is(err, ErrNotFound) { // walks the wrap chain
        fmt.Println("respond 404")
    }

    var ve *ValidationError
    if errors.As(findUser(-1), &ve) { // extracts a typed error from the chain
        fmt.Println("bad field:", ve.Field)
    }
}
```

### Interview points

- `err == ErrNotFound` stops working once the error is wrapped. Use `errors.Is`.
- `%w` wraps and keeps the chain; `%v` only formats the text and cuts the chain.
- `errors.Is` matches a specific error value; `errors.As` matches an error type and gives you its fields.
- `errors.Join` (Go 1.20) combines several errors into one.
- A custom error type joins the chain by implementing `Unwrap() error`.
- Handle an error once: either log it or return it, not both.
- By convention, error strings start lowercase and have no trailing punctuation.
- Go has no exceptions. `panic` is not a substitute for returning errors.

## 8. Generics

Generics (Go 1.18+) let functions and types take type parameters that are restricted by constraints.

```go
type Number interface {
    ~int | ~int64 | ~float64 // ~ also allows types whose underlying type matches
}

func Sum[T Number](xs []T) T {
    var total T
    for _, x := range xs {
        total += x
    }
    return total
}

func Map[T, U any](xs []T, f func(T) U) []U {
    out := make([]U, 0, len(xs))
    for _, x := range xs {
        out = append(out, f(x))
    }
    return out
}

type Stack[T any] struct{ items []T }

func (s *Stack[T]) Push(v T) { s.items = append(s.items, v) }

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    v := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return v, true
}

func main() {
    fmt.Println(Sum([]float64{1.5, 2.5})) // 4: the type argument is inferred
    fmt.Println(Map([]int{1, 2}, strconv.Itoa)) // [1 2] as strings

    st := &Stack[string]{}
    st.Push("a")
}
```

### Interview points

- `any` allows every type. `comparable` allows types that support `==` and is required for map keys.
- `cmp.Ordered` covers types that support `<` and `>`.
- The standard library has generic helpers: the `slices` and `maps` packages, and the built-ins `min`, `max` and `clear` (Go 1.21).
- Use generics for containers and algorithms that work the same across types. Use interfaces when behaviour differs per type.
- Before generics, the options were `interface{}` with type assertions, or code generation.

## 9. Concurrency: goroutines, channels and select

A goroutine is a lightweight thread managed by the Go runtime, and channels pass data between goroutines while synchronizing them.

### Goroutines

```go
func main() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1) // Add before starting the goroutine, never inside it
        go func() {
            defer wg.Done()
            fmt.Println("worker", i)
        }()
    }
    wg.Wait() // without this, main can exit before the goroutines run
}
```

|  | Goroutine | OS thread |
| --- | --- | --- |
| Initial stack | About 2 KB, grows as needed | Usually 1 to 8 MB, fixed |
| Scheduled by | Go runtime, in user space | Operating system kernel |
| Create and switch cost | Very cheap | Expensive |
| Practical count | Hundreds of thousands | A few thousand |

### Channels

```go
unbuf := make(chan int)  // a send blocks until a receiver is ready
buf := make(chan int, 2) // a send blocks only when the buffer is full

go func() { unbuf <- 1 }()
fmt.Println(<-unbuf)

buf <- 1
buf <- 2
close(buf)
for v := range buf { // range ends once the channel is closed and drained
    fmt.Println(v)
}

v, ok := <-buf // 0 false: closed and empty
```

| Operation | nil channel | Closed channel | Open channel |
| --- | --- | --- | --- |
| Send | Blocks forever | Panics | Blocks until there is space or a receiver |
| Receive | Blocks forever | Returns buffered values, then the zero value with ok false | Blocks until a value arrives |
| Close | Panics | Panics | Succeeds |

### Select

```go
select {
case v := <-dataCh:
    fmt.Println("received", v)
case resultCh <- 42:
    fmt.Println("sent")
case <-time.After(2 * time.Second):
    fmt.Println("timeout")
case <-ctx.Done():
    return ctx.Err()
}
```

### Interview points

- When several `select` cases are ready, Go picks one at random.
- A `default` case makes `select` non-blocking. An empty `select {}` blocks forever.
- A nil channel in a `select` case is never chosen, which is a way to switch a case off.
- Only the sender should close a channel. Closing is needed only to tell receivers no more values are coming.
- Directional types restrict use: `chan<- int` is send-only and `<-chan int` is receive-only.
- Concurrency is structuring a program as independent tasks. Parallelism is running tasks at the same instant on several cores.
- Use channels to hand off data or coordinate goroutines. Use a mutex to protect shared state.
- The Go proverb: do not communicate by sharing memory; share memory by communicating.

## 10. Concurrency: sync, atomic and context

Use `sync` primitives to protect shared state, and `context` to carry cancellation and deadlines across goroutines.

### Mutex

```go
type SafeCounter struct {
    mu sync.Mutex
    m  map[string]int
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.m[key]++
}
```

### Once, atomic and Pool

```go
var (
    once sync.Once
    db   *sql.DB
)

func GetDB() *sql.DB {
    once.Do(func() { db = connect() }) // runs exactly once, even under concurrency
    return db
}

var hits atomic.Int64 // lock-free counter
hits.Add(1)
fmt.Println(hits.Load())

var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

buf := bufPool.Get().(*bytes.Buffer)
buf.Reset()
defer bufPool.Put(buf)
```

| Primitive | Use it for |
| --- | --- |
| `sync.Mutex` | Exclusive access to shared state |
| `sync.RWMutex` | Read-heavy state: many readers or one writer |
| `sync.WaitGroup` | Waiting for a set of goroutines to finish |
| `sync.Once` | One-time initialization, singletons |
| `sync.Pool` | Reusing temporary objects to reduce garbage collection work |
| `sync.Cond` | Waking goroutines that wait on a condition; rarely needed |
| `sync.Map` | Concurrent map for keys written once and read many times |
| `sync/atomic` | Lock-free counters and flags |

### Context

```go
func fetch(ctx context.Context, url string) error {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel() // always call cancel: it releases the timer

    req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
    if err != nil {
        return err
    }
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return err // wraps context.DeadlineExceeded after 2 seconds
    }
    defer resp.Body.Close()
    return nil
}

func worker(ctx context.Context, jobs <-chan int) {
    for {
        select {
        case <-ctx.Done():
            return // cancelled or timed out
        case j, ok := <-jobs:
            if !ok {
                return // channel closed
            }
            process(j)
        }
    }
}
```

### Interview points

- Pass `ctx` as the first parameter. Do not store it in a struct, and use `context.TODO()` rather than nil.
- Cancelling a parent context cancels all its children, never the other way round.
- `ctx.Err()` returns `context.Canceled` or `context.DeadlineExceeded`.
- `context.WithValue` is for request-scoped data such as a request ID, not for optional function parameters.
- Never copy a `Mutex` or `WaitGroup` after first use; pass a pointer. `go vet` reports copies.
- Go mutexes are not reentrant: locking the same mutex twice in one goroutine deadlocks.
- Go 1.25 added `wg.Go(f)`, which does the `Add` and `Done` for you.

## 11. Concurrency patterns and classic coding questions

These six programs cover most live-coding concurrency rounds, so practise writing each one from memory.

### Worker pool

```go
func worker(jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs {
        results <- j * 2
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    var wg sync.WaitGroup

    for w := 0; w < 3; w++ {
        wg.Add(1)
        go worker(jobs, results, &wg)
    }

    for i := 1; i <= 9; i++ {
        jobs <- i
    }
    close(jobs) // lets each worker's range loop end

    go func() { wg.Wait(); close(results) }()

    for r := range results {
        fmt.Println(r)
    }
}
```

### Print odd and even numbers with two goroutines

```go
func main() {
    oddTurn, evenTurn := make(chan struct{}), make(chan struct{})
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        for i := 1; i <= 9; i += 2 {
            <-oddTurn
            fmt.Println(i)
            evenTurn <- struct{}{}
        }
    }()

    go func() {
        defer wg.Done()
        for i := 2; i <= 10; i += 2 {
            <-evenTurn
            fmt.Println(i)
            if i < 10 {
                oddTurn <- struct{}{}
            }
        }
    }()

    oddTurn <- struct{}{} // odd goes first
    wg.Wait()
}
```

### Pipeline with fan-out and fan-in

```go
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            out <- n
        }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            out <- n * n
        }
    }()
    return out
}

func merge(chans ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    for _, c := range chans {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }(c)
    }
    go func() { wg.Wait(); close(out) }()
    return out
}

func main() {
    in := gen(1, 2, 3, 4)
    // fan-out: two square workers share one input; fan-in: merge joins them
    for v := range merge(square(in), square(in)) {
        fmt.Println(v)
    }
}
```

### Limit concurrency

```go
// Semaphore: a buffered channel caps work at 5 at a time.
sem := make(chan struct{}, 5)
var wg sync.WaitGroup
for _, u := range urls {
    wg.Add(1)
    sem <- struct{}{} // acquire
    go func(u string) {
        defer wg.Done()
        defer func() { <-sem }() // release
        _ = fetch(ctx, u)
    }(u)
}
wg.Wait()

// errgroup (golang.org/x/sync/errgroup): same limit, plus error handling.
g, ctx := errgroup.WithContext(ctx)
g.SetLimit(5)
for _, u := range urls {
    g.Go(func() error { return fetch(ctx, u) }) // the first error cancels ctx
}
if err := g.Wait(); err != nil {
    return err
}
```

### Rate limiter

```go
ticker := time.NewTicker(200 * time.Millisecond) // 5 requests per second
defer ticker.Stop()

for _, req := range requests {
    <-ticker.C // wait for the next tick
    go handle(req)
}
```

In production, use the token bucket in `golang.org/x/time/rate`: create it with `rate.NewLimiter(5, 10)` and call `limiter.Wait(ctx)` before each request.

### Concurrent-safe cache

```go
type Cache[K comparable, V any] struct {
    mu sync.RWMutex
    m  map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{m: make(map[K]V)}
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
    c.mu.RLock() // many readers at once
    defer c.mu.RUnlock()
    v, ok := c.m[k]
    return v, ok
}

func (c *Cache[K, V]) Set(k K, v V) {
    c.mu.Lock() // one writer, no readers
    defer c.mu.Unlock()
    c.m[k] = v
}
```

Common follow-ups are adding a TTL per entry and LRU eviction, which is a map plus `container/list`.

## 12. Concurrency bugs: races, deadlocks and leaks

Run tests with `go test -race`: it catches data races, the most common Go concurrency bug.

### Data race

```go
count := 0
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        count++ // read, add, write: not atomic
    }()
}
wg.Wait()
fmt.Println(count) // often less than 1000
```

Three fixes: guard `count` with a mutex, use `atomic.Int64`, or let one goroutine own the value and receive updates over a channel.

### Deadlock

```go
func main() {
    ch := make(chan int)
    ch <- 1 // no receiver exists yet
    fmt.Println(<-ch)
}
// fatal error: all goroutines are asleep - deadlock!
```

Other causes: a missing `wg.Done()`, locking the same mutex twice, two goroutines taking two locks in opposite order, and ranging over a channel nobody closes.

### Goroutine leak

```go
func search(query string) string {
    ch := make(chan string) // unbuffered
    go func() { ch <- slowBackend(query) }()

    select {
    case r := <-ch:
        return r
    case <-time.After(time.Second):
        return "timeout" // the goroutine stays blocked on its send forever
    }
}
```

The fix is `make(chan string, 1)` so the send always completes, or a context that the goroutine watches.

### Interview points

- A data race is two goroutines accessing the same variable at once, with at least one writing and no synchronization.
- The race detector reports only races on code paths that actually run, and it slows the program several times over. Use it in tests and CI.
- The runtime reports a deadlock only when every goroutine is blocked. A partial deadlock is silent.
- Find leaks with the pprof goroutine profile or `runtime.NumGoroutine()`.
- Every goroutine you start needs a known way to stop.

## 13. Internals

Senior interviews probe the scheduler, memory allocation and the garbage collector: you need the mental model, not the runtime source.

### Scheduler: the G-M-P model

| Piece | What it is |
| --- | --- |
| G | A goroutine: its stack, instruction pointer and state |
| M | An OS thread that executes code |
| P | A logical processor with a local run queue; there are GOMAXPROCS of them |

- An M must hold a P to run Go code. Each P has a local run queue, and there is one global queue.
- Work stealing: a P with an empty queue takes half the goroutines from another P.
- On a blocking system call, the M blocks and its P is handed to another M, so other goroutines keep running.
- Network I/O goes through the netpoller (epoll or kqueue): the goroutine parks and no thread is blocked.
- Preemption has been asynchronous since Go 1.14: a goroutine in a tight loop is interrupted after about 10 ms.
- GOMAXPROCS defaults to the number of CPU cores. Since Go 1.25 it also respects container CPU limits on Linux.

### Stack, heap and escape analysis

```go
func onStack() int {
    x := 42
    return x // x stays on the stack
}

func onHeap() *int {
    x := 42
    return &x // x escapes: it outlives the function, so it goes to the heap
}
```

Run `go build -gcflags=-m` to see what the compiler decided.

- The compiler decides where a value lives, not your use of `new` or `&`.
- Common causes of escape: returning a pointer, storing a value in an interface, capture by a closure that outlives the function, and very large or dynamically sized values.
- Goroutine stacks start at about 2 KB and grow by copying to a stack twice the size.
- Stack allocation is nearly free and needs no garbage collection, so fewer heap allocations mean less GC work.

### Garbage collector

- It is a concurrent, tri-color mark-and-sweep collector. It is neither generational nor compacting.
- White objects are not yet reached, grey objects are reached but not scanned, and black objects are fully scanned. Whatever is still white at the end is freed.
- A write barrier keeps the marking correct while your program keeps changing pointers.
- Each cycle has two short stop-the-world pauses, usually well under a millisecond.
- `GOGC` (default 100) starts a new cycle when the heap has grown 100% over the live heap. A higher value trades memory for less CPU.
- `GOMEMLIMIT` (Go 1.19) sets a soft memory limit, which is useful in containers.
- To reduce GC pressure: preallocate slices, reuse objects with `sync.Pool`, and avoid needless pointers.

### How built-in types are laid out

| Type | Layout on 64-bit | Consequence |
| --- | --- | --- |
| Slice | Pointer, length, capacity: 24 bytes | Copying the header shares the array |
| String | Pointer, length: 16 bytes | Substrings share memory |
| Interface | Type pointer, data pointer: 16 bytes | Nil only when both words are nil |
| Map | Pointer to a runtime hash table | Passing a map shares it |
| Channel | Pointer to a struct with a ring buffer, a lock and queues of waiting senders and receivers | Passing a channel shares it |

Since Go 1.24 the map is built on a Swiss-table design. Earlier versions used buckets of eight entries with overflow buckets.

### Struct padding

```go
type Bad struct {
    a bool  // 1 byte + 7 padding
    b int64 // 8 bytes
    c bool  // 1 byte + 7 padding
} // 24 bytes

type Good struct {
    b int64 // 8 bytes
    a bool  // 1 byte
    c bool  // 1 byte + 6 padding
} // 16 bytes
```

Order fields from largest to smallest, and check the result with `unsafe.Sizeof`.

## 14. Testing, benchmarking and profiling

Go's standard `testing` package covers unit tests, benchmarks and fuzzing, and table-driven tests are the expected style.

### Table-driven test and benchmark

```go
func Add(a, b int) int { return a + b }

func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -1, -2, -3},
        {"zero", 0, 0, 0},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            if got := Add(tt.a, tt.b); got != tt.want {
                t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}

func BenchmarkAdd(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Add(2, 3)
    }
}
```

### Mocking through an interface

```go
type UserStore interface{ Get(id int) (string, error) }

func Greet(s UserStore, id int) string {
    name, _ := s.Get(id)
    return "hello " + name
}

type fakeStore struct{}

func (fakeStore) Get(id int) (string, error) { return "alice", nil }

func TestGreet(t *testing.T) {
    if got := Greet(fakeStore{}, 1); got != "hello alice" {
        t.Fatalf("got %q", got)
    }
}
```

### Testing an HTTP handler

```go
func TestHealth(t *testing.T) {
    req := httptest.NewRequest(http.MethodGet, "/health", nil)
    rec := httptest.NewRecorder()

    healthHandler(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status = %d, want 200", rec.Code)
    }
}
```

### Commands to know

| Command | Purpose |
| --- | --- |
| `go test ./...` | Run every test in the module |
| `go test -race ./...` | Run tests with the race detector |
| `go test -cover ./...` | Report test coverage |
| `go test -bench=. -benchmem` | Run benchmarks with allocation counts |
| `go test -fuzz=FuzzParse` | Fuzz a function with generated inputs |
| `go vet ./...` | Report suspicious code such as copied locks |
| `go tool pprof cpu.prof` | Analyse a CPU or heap profile |

### Interview points

- Test files end in `_test.go`, and test functions look like `TestXxx(t *testing.T)`.
- `t.Error` records a failure and continues; `t.Fatal` stops that test immediately.
- Know `t.Parallel()`, `t.Helper()`, `t.Cleanup()` and `t.TempDir()`.
- Go 1.24 added `for b.Loop() { ... }` as a simpler benchmark loop than `b.N`.
- Importing `net/http/pprof` exposes profiles at `/debug/pprof`: CPU, heap, goroutine, block and mutex.
- For "memory keeps growing": compare two heap profiles taken minutes apart, then check the goroutine profile for leaks.

## 15. Practical backend Go

Production questions centre on three things: timeouts everywhere, clean shutdown, and correct use of connection pools.

### HTTP server with middleware and graceful shutdown

```go
func logging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        slog.Info("request", "method", r.Method, "path", r.URL.Path, "took", time.Since(start))
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /users/{id}", func(w http.ResponseWriter, r *http.Request) {
        id := r.PathValue("id") // method and wildcard routing: Go 1.22+
        json.NewEncoder(w).Encode(map[string]string{"id": id})
    })

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      logging(mux),
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    go func() {
        if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
            log.Fatal(err)
        }
    }()

    <-ctx.Done() // wait for SIGINT or SIGTERM

    shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    // Shutdown stops accepting connections and waits for in-flight requests.
    if err := srv.Shutdown(shutdownCtx); err != nil {
        log.Printf("shutdown: %v", err)
    }
}
```

### Database access

```go
db, err := sql.Open("postgres", dsn) // creates a pool; does not connect yet
if err != nil {
    return err
}
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(25)
db.SetConnMaxLifetime(5 * time.Minute)

rows, err := db.QueryContext(ctx, "SELECT id, name FROM users WHERE active = $1", true)
if err != nil {
    return err
}
defer rows.Close() // otherwise the connection never returns to the pool

for rows.Next() {
    var id int
    var name string
    if err := rows.Scan(&id, &name); err != nil {
        return err
    }
}
return rows.Err() // errors met during iteration surface here
```

### Transactions

```go
tx, err := db.BeginTx(ctx, nil)
if err != nil {
    return err
}
defer tx.Rollback() // does nothing once Commit has succeeded

if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance - $1 WHERE id = $2", amt, from); err != nil {
    return err
}
if _, err := tx.ExecContext(ctx, "UPDATE accounts SET balance = balance + $1 WHERE id = $2", amt, to); err != nil {
    return err
}
return tx.Commit()
```

### Interview points

- `*sql.DB` is a connection pool that is safe for concurrent use. Create it once and share it.
- Always use placeholders for query values; this prevents SQL injection.
- Each incoming HTTP request runs in its own goroutine, so shared state in handlers needs synchronization.
- `r.Context()` is cancelled when the client disconnects. Pass it to database and downstream calls.
- `http.DefaultClient` has no timeout. Set one, reuse a single client, and always close `resp.Body`.
- JSON encoding sees only exported fields. Control names with tags such as `json:"name,omitempty"` and `json:"-"`.
- Layout convention: `cmd/` for binaries and `internal/` for packages the compiler keeps private to the module.
- Dependency injection needs no framework: pass dependencies into constructors as interfaces.
- `go.mod` lists dependencies and `go.sum` holds their checksums. `go mod tidy` syncs both with the code.

## 16. Rapid-fire revision

One-line answers to the 25 questions that come up most often.

| Question | Short answer |
| --- | --- |
| Is Go pass by value or by reference? | Always by value; slices, maps and channels copy a small header or pointer. |
| `new` vs `make`? | `new(T)` returns a zeroed `*T`; `make` initializes slices, maps and channels. |
| Array vs slice? | An array is a fixed-size value; a slice is a resizable view over an array. |
| nil slice vs empty slice? | Both have length 0; only the nil one equals nil and marshals to JSON null. |
| Why is map iteration order random? | By design, so code never depends on the order. |
| Is a map safe for concurrent use? | No; concurrent writes crash the program. Use a mutex or `sync.Map`. |
| What does `len` return for a string? | The number of bytes, not characters. |
| Value or pointer receiver? | Pointer if the method mutates, the struct is large, or it holds a lock. |
| Does Go have inheritance? | No; it has embedding (composition) and interfaces. |
| When is an interface nil? | Only when both its type and its value are nil. |
| `errors.Is` vs `errors.As`? | `Is` matches a specific error value; `As` extracts an error type. |
| In what order do defers run? | Last in, first out; arguments are evaluated at the `defer` statement. |
| Can you recover another goroutine's panic? | No; `recover` works only in a deferred call in the same goroutine. |
| Buffered vs unbuffered channel? | An unbuffered send waits for a receiver; a buffered send waits only when full. |
| What happens on a closed channel? | Sending panics; receiving returns the zero value with ok false. |
| Mutex or channel? | A mutex protects shared state; a channel passes data and coordinates goroutines. |
| Concurrency vs parallelism? | Concurrency structures independent tasks; parallelism runs them at the same instant. |
| What is `context` for? | Cancellation, deadlines and request-scoped values across API boundaries. |
| What is a goroutine leak? | A goroutine blocked forever, usually on a channel nobody reads or closes. |
| How do you detect data races? | `go test -race` or `go run -race`. |
| What is escape analysis? | The compiler deciding whether a variable lives on the stack or the heap. |
| What kind of GC does Go use? | Concurrent tri-color mark-and-sweep; not generational, not compacting. |
| What does GOMAXPROCS control? | The number of Ps, which caps how many goroutines run Go code in parallel. |
| What does `init()` do? | It runs automatically before `main`, after package variables are initialized. |
| Why use the empty struct? | It takes zero bytes; used for sets and signal-only channels. |
