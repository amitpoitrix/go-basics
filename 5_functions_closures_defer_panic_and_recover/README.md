# Functions, Closures, Defer, Panic and Recover in Go

Oct 11, 2026 · @HR

## 1. Function basics

A Go function can return several values, name its results, and accept a variable number of arguments.

```go
func add(a, b int) int { return a + b }

// Multiple results: the error comes last by convention
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

// Variadic parameter plus named results
func bounds(nums ...int) (lo, hi int) {
    if len(nums) == 0 {
        return // naked return: sends back lo and hi as they are (0, 0)
    }
    lo, hi = nums[0], nums[0]
    for _, n := range nums[1:] {
        if n < lo {
            lo = n
        }
        if n > hi {
            hi = n
        }
    }
    return
}

lo, hi := bounds(3, 1, 4) // 1 4

vals := []int{5, 2, 9}
lo, hi = bounds(vals...) // spread a slice: 2 9

_, err := divide(1, 0) // _ discards a result you do not need
```

### A variadic call shares the slice

```go
func zeroFirst(nums ...int) { nums[0] = 0 }

vals := []int{5, 2, 9}
zeroFirst(vals...)
fmt.Println(vals) // [0 2 9]: no copy is made when you spread a slice
```

- Arguments are always passed by value. To change the caller's data, pass a pointer.
- Go has no default parameter values and no function overloading.
- A variadic parameter must be last, and inside the function it is an ordinary slice. With no arguments it is nil.
- Named results start at their zero values and act as local variables.
- Naked returns are fine in short functions and hurt readability in long ones.
- An unused local variable or import is a compile error. An unused parameter is allowed.

## 2. Functions as values

Functions are first-class values: you can store them in variables, pass them as arguments, return them, and put them in maps.

```go
type Op func(int, int) int // a named function type

func apply(a, b int, op Op) int { return op(a, b) }

add := func(a, b int) int { return a + b } // anonymous function in a variable

fmt.Println(apply(2, 3, add))                                 // 5
fmt.Println(apply(2, 3, func(a, b int) int { return a * b })) // 6

ops := map[string]Op{
    "add": add,
    "sub": func(a, b int) int { return a - b },
}
fmt.Println(ops["sub"](5, 2)) // 3

func() {
    fmt.Println("runs immediately")
}() // define and call in one step

var f func()          // the zero value of a function type is nil
fmt.Println(f == nil) // true
// f()                // panic: calling a nil function
```

### A function type can have methods

```go
type Handler interface {
    Handle(string) string
}

type HandlerFunc func(string) string

func (f HandlerFunc) Handle(s string) string { return f(s) } // the adapter

var h Handler = HandlerFunc(strings.ToUpper) // a plain function now satisfies the interface
fmt.Println(h.Handle("go"))                   // GO
```

This is exactly how `http.HandlerFunc` lets an ordinary function serve as an `http.Handler`.

- A function's type is its signature: parameter types and result types.
- Functions can be compared only against nil, so they cannot be map keys.
- A function that takes or returns another function is called a higher-order function.
- Passing a function is the idiomatic way to inject behaviour: comparators, callbacks, middleware and options.

## 3. Closures

A closure is a function that uses variables from the scope where it was defined, and it captures the variables themselves, not copies of their values.

```go
func counter() func() int {
    n := 0
    return func() int {
        n++ // n lives on after counter returns
        return n
    }
}

a, b := counter(), counter()
fmt.Println(a(), a(), a()) // 1 2 3
fmt.Println(b())           // 1: each call to counter creates a separate n
```

### Capture is by reference

```go
x := 1
show := func() { fmt.Println(x) }
x = 2
show() // 2: the closure reads the variable when it runs
```

### Loop variables

```go
var fns []func()
for i := 0; i < 3; i++ {
    fns = append(fns, func() { fmt.Print(i, " ") })
}
for _, f := range fns {
    f()
}
// Go 1.22+: 0 1 2 (each iteration has its own i)
// Before 1.22: 3 3 3 (all closures share one i)

// The fix needed before 1.22, still seen in older code:
for i := 0; i < 3; i++ {
    i := i // a new variable for each iteration
    fns = append(fns, func() { fmt.Print(i, " ") })
}
```

### Practical uses

```go
// Middleware: wrap a handler with extra behaviour
func logging(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next(w, r)
        log.Println(r.URL.Path, time.Since(start))
    }
}

// Memoization: remember results between calls
func memoize(f func(int) int) func(int) int {
    cache := map[int]int{}
    return func(n int) int {
        if v, ok := cache[n]; ok {
            return v
        }
        v := f(n)
        cache[n] = v
        return v
    }
}
```

- A closure is the function plus the variables it refers to. Those variables move to the heap so they outlive the enclosing call.
- Closures give you private state without a struct, as in the counter above.
- A captured variable shared between goroutines needs the same synchronization as any shared data.
- The Go 1.22 change applies only to variables declared by the `for` statement. A variable declared outside the loop is still shared.

## 4. Defer

`defer` schedules a call to run when the surrounding function returns, whether it returns normally or through a panic.

```go
func readFile(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err // check the error first: f is nil here
    }
    defer f.Close() // runs on every return path below

    data, err := io.ReadAll(f)
    return string(data), err
}
```

### Rule 1: last in, first out

```go
func main() {
    defer fmt.Println("1")
    defer fmt.Println("2")
    defer fmt.Println("3")
    fmt.Println("body")
}
// body, 3, 2, 1
```

### Rule 2: arguments are evaluated at the defer statement

```go
x := 1
defer fmt.Println("argument:", x)             // x is read now: prints 1
defer func() { fmt.Println("closure:", x) }() // x is read when it runs: prints 2
x = 2
// closure: 2, then argument: 1
```

### Rule 3: a deferred closure can change named results

```go
func named() (result int) {
    defer func() { result *= 2 }()
    return 5 // sets result to 5, then the defer doubles it: returns 10
}

func unnamed() int {
    x := 5
    defer func() { x *= 2 }()
    return x // the return value is copied first: returns 5
}
```

### Useful patterns

```go
// Unlock on every path
mu.Lock()
defer mu.Unlock()

// Add context to any error the function returns
func process() (err error) {
    defer func() {
        if err != nil {
            err = fmt.Errorf("process: %w", err)
        }
    }()
    return step()
}

// Time a function: timer("work") runs now, the function it returns runs at exit
func timer(name string) func() {
    start := time.Now()
    return func() { fmt.Println(name, "took", time.Since(start)) }
}

func work() {
    defer timer("work")()
    // ...
}
```

### Defer inside a loop

```go
// Problem: every file stays open until the function returns
for _, p := range paths {
    f, _ := os.Open(p)
    defer f.Close()
}

// Fix: give each iteration its own function scope
for _, p := range paths {
    func() {
        f, err := os.Open(p)
        if err != nil {
            return
        }
        defer f.Close() // runs at the end of this iteration
        // use f
    }()
}
```

- A `return` statement does three things in order: sets the result, runs the deferred calls, then returns to the caller.
- Deferred calls run on a normal return and during a panic. They do not run after `os.Exit`.
- Values returned by a deferred function are discarded.
- Deferring a nil function panics when the surrounding function exits, not at the `defer` line.
- Since Go 1.14 a defer costs almost nothing in the common case, so use it freely for cleanup.

## 5. Panic

A panic stops the normal flow of a goroutine, runs its deferred calls while unwinding the stack, and crashes the program unless something recovers it.

```go
func main() {
    defer fmt.Println("deferred call runs")
    panic("something went wrong")
    fmt.Println("never reached")
}
// deferred call runs
// panic: something went wrong
// goroutine 1 [running]: ...stack trace...
// exit status 2
```

### What causes a panic

| Cause | Example |
| --- | --- |
| Explicit call | `panic("bad state")` or `panic(err)` |
| Nil pointer dereference | Reading a field through a nil pointer |
| Index out of range | `s[5]` on a slice of length 3 |
| Write to a nil map | `m["k"] = 1` when m is nil |
| Failed type assertion | `v.(int)` when v holds a string, in the single-result form |
| Integer division by zero | `a / b` when b is 0 |
| Channel misuse | Sending on a closed channel, or closing one twice |
| Calling a nil function | `f()` when f is nil |

### Fatal errors are different

| Fatal error | Recoverable |
| --- | --- |
| Concurrent map writes | No |
| All goroutines asleep (deadlock) | No |
| Stack overflow from unbounded recursion | No |
| Out of memory | No |

- `panic` accepts any value, usually a string or an error.
- The panic travels up the call stack of one goroutine, running each function's deferred calls on the way.
- If it reaches the top of that goroutine unrecovered, the whole program exits, not just the goroutine.
- Use panic for programmer mistakes and impossible states. Return an error for anything a caller could reasonably expect.
- A library should not let a panic escape its public API.

## 6. Recover

`recover` stops a panic and returns its value, but only when called directly inside a deferred function in the panicking goroutine.

```go
func safeDivide(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered: %v", r) // turn the panic into an error
        }
    }()
    return a / b, nil
}

fmt.Println(safeDivide(10, 2)) // 5 <nil>
fmt.Println(safeDivide(1, 0))  // 0 recovered: runtime error: integer divide by zero
```

### Where recover works

```go
defer func() { recover() }() // works: called directly by the deferred function

defer recover() // does not work: recover is the deferred call itself

defer func() {
    helper() // does not work if recover is called inside helper
}()
```

### Recover is per goroutine

```go
func main() {
    defer func() { recover() }() // cannot catch a panic in another goroutine

    go func() {
        panic("boom") // crashes the whole program
    }()
    time.Sleep(time.Second)
}
```

### Standard patterns

```go
// Start a goroutine that cannot take the program down
func safeGo(fn func()) {
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("goroutine panic: %v\n%s", r, debug.Stack())
            }
        }()
        fn()
    }()
}

// HTTP middleware: one bad request must not kill the server
func recoverer(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if rec := recover(); rec != nil {
                log.Printf("panic: %v", rec)
                http.Error(w, "internal server error", http.StatusInternalServerError)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

- `recover` returns nil when there is no panic in progress.
- After a recover, the panicking function does not resume. It returns to its caller with whatever its named results hold.
- To recover and still return a useful value, the function needs named results, as `safeDivide` has.
- You can recover, inspect the value, and call `panic(r)` again if you cannot handle it.
- `net/http` already recovers a panic in a handler and logs it, but the client just sees a dropped connection. The middleware returns a proper 500.
- Fatal errors such as concurrent map writes cannot be recovered.
- Do not use panic and recover as general exception handling. Errors are the normal path.

## 7. Errors vs panic, and exit paths

Return an error for anything that can go wrong in normal operation, and reserve panic for bugs and for startup failures the program cannot survive.

| Situation | Use |
| --- | --- |
| File not found, bad user input, network timeout | Return an error |
| Database row missing | Return an error |
| Index out of range caused by your own logic bug | Panic (the runtime does it for you) |
| A switch reaches a case that should be impossible | Panic |
| Invalid hard-coded regular expression or template at startup | Panic, through a `Must` function |

### The Must pattern

```go
// Panics at startup if the pattern is wrong, which is a programmer error
var emailRe = regexp.MustCompile(`^[a-z]+@[a-z]+$`)

// A generic Must for your own setup code
func Must[T any](v T, err error) T {
    if err != nil {
        panic(err)
    }
    return v
}

cfg := Must(loadConfig("app.yaml"))
```

### Which exits run deferred calls

| How the function ends | Deferred calls run |
| --- | --- |
| `return` | Yes |
| `panic` | Yes, while the stack unwinds |
| `runtime.Goexit()` | Yes |
| `os.Exit(code)` | No |
| `log.Fatal(...)` | No, because it calls `os.Exit(1)` |

```go
// Problem: cleanup never runs
func main() {
    defer fmt.Println("cleanup")
    log.Fatal("bad config")
}

// Fix: keep main tiny and put the real work in a function that returns an error
func main() {
    if err := run(); err != nil {
        log.Fatal(err) // run has already returned, so its defers have run
    }
}

func run() error {
    defer fmt.Println("cleanup")
    return errors.New("bad config")
}
```

- Go has no exceptions and no try-catch. Explicit error returns keep the failure path visible.
- Reserve `Must` functions for package-level initialization, never for request handling.
- `log.Panic` logs and then panics, so deferred calls do run.

## 8. init and program startup

`init` is a special function that Go runs automatically before `main`, once per declaration.

```go
package main

var config = loadDefaults() // 1. package-level variables are initialized first

func init() { // 2. then every init function runs
    fmt.Println("init")
}

func main() { // 3. finally main starts
    fmt.Println("main")
}
```

The startup order is:

1. Imported packages are initialized first, each exactly once, dependencies before dependents.
2. Within a package, package-level variables are initialized in dependency order.
3. Then the package's `init` functions run, in the order the files are given to the compiler.
4. After every package is ready, `main.main` runs.

- `init` takes no arguments, returns nothing, and cannot be called from your code.
- A package, and even a single file, may contain several `init` functions.
- A blank import such as `import _ "net/http/pprof"` exists only to run that package's `init`.
- Keep `init` small. Hidden work there makes programs hard to test and slow to start.
- When `main` returns, the program exits immediately without waiting for other goroutines.

## 9. Predict the output

Interviewers love these five puzzles; cover the answer comments and work each one out first.

### Puzzle 1: defer and return values

```go
func a() int {
    x := 1
    defer func() { x++ }()
    return x
}

func b() (x int) {
    x = 1
    defer func() { x++ }()
    return x
}

func c() (x int) {
    defer func() { x++ }()
    return 5
}

fmt.Println(a(), b(), c()) // 1 2 6
```

Only a named result can be changed by a defer. In `a`, the return value was copied before `x++` ran.

### Puzzle 2: defer with a method receiver

```go
type T struct{ name string }

func (t T) Print()     { fmt.Println(t.name) }
func (t *T) PrintPtr() { fmt.Println(t.name) }

func main() {
    t := T{"first"}
    defer t.Print()    // the receiver value is copied now
    defer t.PrintPtr() // only the address is saved now
    t.name = "second"
}
// second
// first
```

### Puzzle 3: defer order around a recover

```go
func main() {
    defer fmt.Println("A")
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered:", r)
        }
    }()
    defer fmt.Println("B")
    panic("boom")
}
// B
// recovered: boom
// A
```

Deferred calls run last in, first out. After the recover, the remaining defers still run and the program exits normally.

### Puzzle 4: a closure over an outer variable

```go
func main() {
    var fns []func()
    x := 0
    for i := 0; i < 3; i++ {
        x += i
        fns = append(fns, func() { fmt.Print(x, " ") })
    }
    for _, f := range fns {
        f()
    }
}
// 3 3 3
```

`x` is declared outside the loop, so all three closures share it, even in Go 1.22 and later.

### Puzzle 5: evaluation inside a defer statement

```go
func enter(s string) string {
    fmt.Println("enter", s)
    return s
}

func leave(s string) { fmt.Println("leave", s) }

func work() {
    defer leave(enter("work")) // enter runs now; leave runs at exit
    fmt.Println("working")
}
// enter work
// working
// leave work
```

## 10. Interview coding problems

These four exercises show the main reasons to pass functions around: retrying, caching, wrapping and transforming.

### Retry with exponential backoff

```go
func retry(attempts int, delay time.Duration, fn func() error) error {
    var err error
    for i := 0; i < attempts; i++ {
        if err = fn(); err == nil {
            return nil
        }
        if i < attempts-1 {
            time.Sleep(delay)
            delay *= 2 // wait twice as long before the next try
        }
    }
    return fmt.Errorf("after %d attempts: %w", attempts, err)
}

err := retry(3, 100*time.Millisecond, func() error {
    return callAPI()
})
```

### Memoized Fibonacci with a recursive closure

```go
func fibMemo() func(int) int {
    cache := map[int]int{}

    var fib func(int) int // declare first, so the closure can call itself
    fib = func(n int) int {
        if n < 2 {
            return n
        }
        if v, ok := cache[n]; ok {
            return v
        }
        cache[n] = fib(n-1) + fib(n-2)
        return cache[n]
    }
    return fib
}

fib := fibMemo()
fmt.Println(fib(50)) // 12586269025, computed in O(n)
```

### Middleware chain

```go
type Middleware func(http.Handler) http.Handler

func chain(h http.Handler, mws ...Middleware) http.Handler {
    for i := len(mws) - 1; i >= 0; i-- { // wrap in reverse, so the first listed runs first
        h = mws[i](h)
    }
    return h
}

handler := chain(mux, recoverer, requestLog, auth) // each is a func(http.Handler) http.Handler
// a request passes through recoverer, then requestLog, then auth, then mux
```

### Generic Filter, Map and Reduce

```go
func Filter[T any](xs []T, keep func(T) bool) []T {
    var out []T
    for _, x := range xs {
        if keep(x) {
            out = append(out, x)
        }
    }
    return out
}

func Map[T, U any](xs []T, f func(T) U) []U {
    out := make([]U, 0, len(xs))
    for _, x := range xs {
        out = append(out, f(x))
    }
    return out
}

func Reduce[T, A any](xs []T, initial A, f func(A, T) A) A {
    acc := initial
    for _, x := range xs {
        acc = f(acc, x)
    }
    return acc
}

evens := Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 }) // [2 4]
squares := Map(evens, func(n int) int { return n * n })                  // [4 16]
sum := Reduce(squares, 0, func(acc, n int) int { return acc + n })       // 20
```

- For the retry, expect follow-ups on adding a `context.Context` for cancellation and random jitter to the delay.
- The memoized closure is not safe for concurrent use; a mutex around the cache fixes that.
- The `safeGo` and `recoverer` functions in section 6 are also common whiteboard requests.

## 11. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| Does Go support function overloading or default arguments? | No to both. |
| What is a variadic function? | One whose last parameter takes any number of arguments, received as a slice. |
| What are named results? | Result variables declared in the signature, which a deferred closure can modify. |
| Are functions first-class in Go? | Yes; they can be stored, passed and returned. |
| Can two functions be compared? | Only against nil. |
| What is a closure? | A function together with the outer variables it uses. |
| Does a closure capture values or variables? | Variables, so it sees later changes. |
| What changed for loop variables in Go 1.22? | Each iteration gets its own copy, which fixed the classic closure bug. |
| In what order do deferred calls run? | Last in, first out. |
| When are a deferred call's arguments evaluated? | At the `defer` statement, not when the call runs. |
| Can a defer change the return value? | Yes, if the result is named. |
| Do defers run after a panic? | Yes, as the stack unwinds. |
| Do defers run after `os.Exit` or `log.Fatal`? | No. |
| What is wrong with defer in a loop? | Resources stay open until the whole function returns. |
| What does panic do? | Stops normal flow, runs defers, and crashes the program if not recovered. |
| Where must recover be called? | Directly inside a deferred function, in the panicking goroutine. |
| Can one goroutine recover another's panic? | No. |
| What does recover return with no panic in progress? | nil. |
| Can every crash be recovered? | No; fatal errors such as concurrent map writes and deadlocks cannot. |
| Error or panic? | Errors for expected failures; panic for bugs and impossible states. |
| When does `init` run? | Before `main`, after package-level variables are initialized. |
