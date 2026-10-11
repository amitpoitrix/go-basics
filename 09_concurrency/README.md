# Concurrency in Go

## 1. Goroutines

A goroutine is a function running independently of its caller, started with the `go` keyword and scheduled by the Go runtime rather than the operating system.

```go
func main() {
    go fmt.Println("from a goroutine") // may never print
    fmt.Println("from main")
} // when main returns, the program exits and every other goroutine is killed
```

### Waiting for goroutines

```go
func main() {
    var wg sync.WaitGroup
    for i := 1; i <= 3; i++ {
        wg.Add(1) // count the goroutine before starting it
        go func() {
            defer wg.Done()
            fmt.Println("worker", i)
        }()
    }
    wg.Wait() // blocks until the counter is back to zero
}
// worker 1, worker 2, worker 3, in any order

// Go 1.25 and later: wg.Go does the Add and Done for you
var wg sync.WaitGroup
for i := 1; i <= 3; i++ {
    wg.Go(func() { fmt.Println("worker", i) })
}
wg.Wait()
```

### Concurrency is not parallelism

|  | Concurrency | Parallelism |
| --- | --- | --- |
| Meaning | Structuring a program as independently running tasks | Executing several tasks at the same instant |
| Needs multiple cores | No | Yes |
| Analogy | One cook switching between three dishes | Three cooks, one dish each |
| In Go | Goroutines and channels | What the runtime does with them when cores are available |

### Goroutines vs OS threads

|  | Goroutine | OS thread |
| --- | --- | --- |
| Initial stack | About 2 KB, grows as needed | Usually 1 to 8 MB, fixed |
| Scheduled by | Go runtime, in user space | Operating system kernel |
| Create and switch cost | Very cheap | Expensive |
| Practical count | Hundreds of thousands | A few thousand |

- The order in which goroutines run is not defined. Never rely on it.
- A goroutine has no ID and no handle. You cannot kill one from outside; it must return on its own.
- The value a goroutine's function returns is discarded. Results travel back through channels or shared memory.
- An unrecovered panic in any goroutine crashes the whole program.
- `time.Sleep` is not a way to wait for a goroutine. Use a `WaitGroup` or a channel.

## 2. Channels

A channel is a typed pipe that goroutines use to pass values to each other, and it synchronizes them as a side effect.

```go
unbuf := make(chan int)  // unbuffered: a send blocks until a receiver takes the value
buf := make(chan int, 2) // buffered: a send blocks only when the buffer is full

go func() { unbuf <- 42 }() // send
v := <-unbuf                // receive: both sides meet at this moment

buf <- 1
buf <- 2 // a third send would block
fmt.Println(len(buf), cap(buf)) // 2 2
```

### Closing and ranging

```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
close(ch) // no more values will be sent

for v := range ch { // receives until the channel is closed and empty
    fmt.Println(v) // 1, then 2
}

v, ok := <-ch // 0 false: closed and drained
```

### What each operation does

| Operation | nil channel | Closed channel | Open channel |
| --- | --- | --- | --- |
| Send | Blocks forever | Panics | Blocks until there is space or a receiver |
| Receive | Blocks forever | Returns buffered values, then the zero value with ok false | Blocks until a value arrives |
| Close | Panics | Panics | Succeeds |

### Direction and ownership

```go
func produce(n int) <-chan int { // returns a receive-only channel
    out := make(chan int)
    go func() {
        defer close(out) // the sender owns the channel, so the sender closes it
        for i := 0; i < n; i++ {
            out <- i
        }
    }()
    return out
}

func consume(in <-chan int) { // can only receive
    for v := range in {
        fmt.Println(v)
    }
}

func sink(out chan<- int) { out <- 1 } // can only send
```

### A channel as a signal

```go
done := make(chan struct{})
go func() {
    work()
    close(done) // closing wakes every receiver at once
}()
<-done
```

- An unbuffered channel is a hand-off: sender and receiver synchronize on every value.
- A buffered channel decouples them up to its capacity. It is a queue, not a way to fix a deadlock.
- Only the sender should close, and only when receivers need to know that no more values are coming.
- You do not have to close every channel. An unreferenced channel is garbage collected like any value.
- With several senders, coordinate the close with a `WaitGroup`, so it happens exactly once after all of them finish.

## 3. select

`select` waits on several channel operations at once and runs whichever one is ready first.

```go
select {
case v := <-dataCh:
    fmt.Println("received", v)
case resultCh <- 42:
    fmt.Println("sent")
case <-time.After(2 * time.Second):
    fmt.Println("timed out")
case <-ctx.Done():
    return ctx.Err()
}
```

### The for-select loop

```go
func worker(ctx context.Context, jobs <-chan Job) {
    for {
        select {
        case <-ctx.Done():
            return // told to stop
        case job, ok := <-jobs:
            if !ok {
                return // channel closed: no more work
            }
            process(job)
        }
    }
}
```

### Non-blocking operations with default

```go
select {
case v := <-ch:
    fmt.Println("got", v)
default:
    fmt.Println("nothing ready") // runs at once when no case is ready
}

select {
case ch <- event:
default:
    dropped++ // the buffer is full: drop instead of blocking
}
```

### Turning a case off with a nil channel

```go
func merge(a, b <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for a != nil || b != nil {
            select {
            case v, ok := <-a:
                if !ok {
                    a = nil // a nil channel is never ready, so this case is disabled
                    continue
                }
                out <- v
            case v, ok := <-b:
                if !ok {
                    b = nil
                    continue
                }
                out <- v
            }
        }
    }()
    return out
}
```

- When several cases are ready, `select` picks one at random, which prevents starvation.
- With no ready case and no `default`, `select` blocks. An empty `select {}` blocks forever.
- A receive from a closed channel is always ready, so a closed channel in a loop spins unless you set it to nil.
- A `default` inside a `for` loop becomes a busy loop that burns a CPU core. Block on a channel instead.
- Since Go 1.23, timers that are no longer referenced are garbage collected, so `time.After` in a loop no longer piles up memory.

## 4. The sync and atomic packages

When goroutines share memory, the `sync` package provides locks and coordination tools, and `sync/atomic` provides lock-free operations on single values.

### Mutex and RWMutex

```go
type Counter struct {
    mu sync.Mutex
    n  map[string]int
}

func (c *Counter) Inc(key string) {
    c.mu.Lock()
    defer c.mu.Unlock() // released on every return path, even after a panic
    c.n[key]++
}

type Config struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Config) Get(k string) string {
    c.mu.RLock() // many readers may hold the read lock together
    defer c.mu.RUnlock()
    return c.data[k]
}

func (c *Config) Set(k, v string) {
    c.mu.Lock() // a writer excludes everyone else
    defer c.mu.Unlock()
    c.data[k] = v
}
```

### Once and Pool

```go
var (
    once sync.Once
    db   *sql.DB
)

func DB() *sql.DB {
    once.Do(func() { db = connect() }) // runs exactly once, even with many callers
    return db
}

// Go 1.21: the same thing in one line
var getDB = sync.OnceValue(connect)

// Pool: reuse temporary objects to reduce garbage
var bufPool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

func render() string {
    buf := bufPool.Get().(*bytes.Buffer)
    buf.Reset() // a pooled object still holds its old contents
    defer bufPool.Put(buf)
    buf.WriteString("hello")
    return buf.String()
}
```

### Atomic operations

```go
var requests atomic.Int64
requests.Add(1)
fmt.Println(requests.Load())

var ready atomic.Bool
ready.Store(true)

var state atomic.Int32
if state.CompareAndSwap(0, 1) { // set to 1 only if it is still 0
    // exactly one goroutine gets here
}
```

| Tool | Use it for |
| --- | --- |
| `sync.Mutex` | Exclusive access to shared state |
| `sync.RWMutex` | State that is read far more often than written |
| `sync.WaitGroup` | Waiting for a group of goroutines to finish |
| `sync.Once`, `sync.OnceValue` | One-time initialization |
| `sync.Pool` | Reusing short-lived objects in hot paths |
| `sync.Cond` | Waking goroutines that wait for a condition; rarely needed, a channel usually does it |
| `sync.Map` | A concurrent map for write-once, read-many keys |
| `sync/atomic` | Counters, flags and pointer swaps on a single value |

- Never copy a `Mutex`, `WaitGroup` or any struct containing one after first use. Pass a pointer. `go vet` reports copies.
- Go mutexes are not reentrant. Locking one you already hold deadlocks.
- Keep critical sections short, and never do I/O or send on a channel while holding a lock if you can avoid it.
- When you need two locks, always take them in the same order everywhere.
- An `RWMutex` is slower than a `Mutex` when writes are frequent. Use it only for read-heavy data.
- Atomics protect one value. An invariant that spans two fields needs a mutex.
- Objects in a `sync.Pool` can be removed by the garbage collector at any time, so it is a cache, not storage.

## 5. context

A `context.Context` carries a cancellation signal, a deadline and request-scoped values down a call chain and across goroutines.

### Cancellation

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())

    go worker(ctx)

    time.Sleep(time.Second)
    cancel() // closes ctx.Done(), which every goroutine holding ctx can see
    time.Sleep(100 * time.Millisecond)
}

func worker(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            fmt.Println("stopping:", ctx.Err()) // stopping: context canceled
            return
        default:
            doSomeWork()
        }
    }
}
```

### Timeouts and deadlines

```go
func fetchUser(ctx context.Context, id int) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
    defer cancel() // always call cancel: it releases the timer

    var u User
    err := db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", id).
        Scan(&u.ID, &u.Name)
    if errors.Is(err, context.DeadlineExceeded) {
        return nil, fmt.Errorf("fetch user %d: timed out", id)
    }
    return &u, err
}
```

### Request-scoped values

```go
type ctxKey string // a private key type prevents collisions between packages

const requestIDKey ctxKey = "requestID"

ctx = context.WithValue(ctx, requestIDKey, "abc-123")

id, ok := ctx.Value(requestIDKey).(string) // abc-123 true
```

| Function | Creates a context that |
| --- | --- |
| `context.Background()` | Is never cancelled; the root, used in `main` and tests |
| `context.TODO()` | Is a placeholder when you are not yet sure which context to use |
| `context.WithCancel(parent)` | Is cancelled when you call `cancel` |
| `context.WithTimeout(parent, d)` | Is cancelled after a duration |
| `context.WithDeadline(parent, t)` | Is cancelled at a point in time |
| `context.WithValue(parent, k, v)` | Carries one extra key-value pair |

- Pass `ctx` as the first parameter of a function. Do not store it in a struct.
- Cancelling a parent cancels all its children. Cancelling a child never affects the parent.
- `ctx.Err()` returns `context.Canceled` or `context.DeadlineExceeded` once the context is done.
- Cancellation is cooperative. A goroutine stops only if its code checks `ctx.Done()` or calls something that does.
- Use values only for data that belongs to the request, such as a request ID or the authenticated user. Never use them for optional function parameters.
- In an HTTP handler, `r.Context()` is cancelled when the client disconnects. Pass it to every database and outbound call.

## 6. Core patterns: worker pool, pipeline, fan-out and fan-in

These three patterns structure most concurrent Go programs, and each is a standard live-coding request.

### Worker pool

```go
func worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    for j := range jobs { // exits when jobs is closed and empty
        results <- j * j
    }
}

func main() {
    jobs := make(chan int, 100)
    results := make(chan int, 100)
    var wg sync.WaitGroup

    for w := 1; w <= 3; w++ { // a fixed number of workers bounds the concurrency
        wg.Add(1)
        go worker(w, jobs, results, &wg)
    }

    for i := 1; i <= 9; i++ {
        jobs <- i
    }
    close(jobs) // tell the workers there is no more work

    go func() {
        wg.Wait()
        close(results) // close results only after every worker has finished
    }()

    for r := range results {
        fmt.Println(r)
    }
}
```

### Pipeline

```go
// Each stage takes a receive-only channel and returns one
func generate(nums ...int) <-chan int {
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

for v := range square(square(generate(1, 2, 3))) {
    fmt.Println(v) // 1, 16, 81
}
```

### Fan-out and fan-in

```go
// Fan-in: merge several channels into one
func merge(chans ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    for _, c := range chans {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for v := range c {
                out <- v
            }
        }()
    }
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}

func main() {
    in := generate(1, 2, 3, 4, 5, 6)

    // Fan-out: three square stages read from the same input channel
    w1, w2, w3 := square(in), square(in), square(in)

    for v := range merge(w1, w2, w3) {
        fmt.Println(v) // the six squares, in no fixed order
    }
}
```

- The worker pool's two closes are the part people get wrong: close `jobs` after sending, and close `results` only after `wg.Wait()`.
- In a pipeline, each stage closes its own output when its input is exhausted, so shutdown flows downstream by itself.
- Fan-out gives parallelism for CPU-heavy or slow stages. Fan-in brings the results back to one consumer.
- Add a `ctx context.Context` to every stage in real code, so a consumer that stops early does not leave upstream goroutines blocked.

## 7. Limiting and coordinating

Unbounded goroutines exhaust memory, connections and downstream services, so production code limits how much runs at once.

### Semaphore: cap concurrency with a buffered channel

```go
sem := make(chan struct{}, 5) // at most 5 running at a time
var wg sync.WaitGroup

for _, url := range urls {
    wg.Add(1)
    sem <- struct{}{} // acquire: blocks while 5 are already running
    go func() {
        defer wg.Done()
        defer func() { <-sem }() // release
        fetch(ctx, url)
    }()
}
wg.Wait()
```

### errgroup: limit, wait, and stop on the first error

```go
import "golang.org/x/sync/errgroup"

g, ctx := errgroup.WithContext(ctx)
g.SetLimit(5)

for _, url := range urls {
    g.Go(func() error {
        return fetch(ctx, url) // the first error cancels ctx for the rest
    })
}

if err := g.Wait(); err != nil { // waits for all; returns the first error
    return err
}
```

### Rate limiting

```go
// Simple: one request per tick
ticker := time.NewTicker(200 * time.Millisecond) // 5 per second
defer ticker.Stop()

for _, req := range requests {
    <-ticker.C
    go handle(req)
}

// Production: a token bucket with bursts
import "golang.org/x/time/rate"

limiter := rate.NewLimiter(rate.Limit(5), 10) // 5 per second, bursts of up to 10

for _, req := range requests {
    if err := limiter.Wait(ctx); err != nil { // blocks until a token is free
        return err
    }
    go handle(req)
}
```

### First result wins

```go
func fastest(ctx context.Context, replicas []string) string {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel() // stops the slower requests once we have an answer

    ch := make(chan string, len(replicas)) // buffered, so the losers do not block forever
    for _, addr := range replicas {
        go func() { ch <- query(ctx, addr) }()
    }
    return <-ch
}
```

| Need | Tool |
| --- | --- |
| At most N things running at once | A buffered channel as a semaphore, or `errgroup.SetLimit` |
| Stop everything on the first failure | `errgroup.WithContext` |
| At most N operations per second | `time.Ticker`, or `rate.Limiter` for bursts |
| Collapse duplicate concurrent calls for one key | `golang.org/x/sync/singleflight` |
| Run something once | `sync.Once` |

- A semaphore limits how many run together. A rate limiter limits how many start per unit of time. They solve different problems.
- `errgroup` replaces a `WaitGroup`, an error channel and a cancel call with one type.
- `singleflight` is the standard answer to a cache stampede: many requests for the same missing key trigger only one load.

## 8. Channels or a mutex

Use a channel to pass data or signals between goroutines, and a mutex to protect state that several goroutines read and write in place.

| Situation | Better tool |
| --- | --- |
| Handing work or results from one goroutine to another | Channel |
| Signalling that something finished or should stop | Channel, or `context` |
| Distributing jobs to workers | Channel |
| Guarding a map, cache or struct that many goroutines touch | Mutex |
| A simple counter or flag | `sync/atomic` |
| One-time setup | `sync.Once` |

### The same counter, three ways

```go
// 1. Mutex: the usual choice for shared state
type Counter struct {
    mu sync.Mutex
    n  int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.n++
    c.mu.Unlock()
}

// 2. Atomic: the fastest for a single number
var n atomic.Int64
n.Add(1)

// 3. Channel: one goroutine owns the value, others send it requests
func counter(inc <-chan struct{}, read chan<- int) {
    n := 0
    for {
        select {
        case <-inc:
            n++
        case read <- n:
        }
    }
}
```

- The proverb is "do not communicate by sharing memory; share memory by communicating". It describes a default, not a law.
- The channel version above is correct but slower and longer. Do not force channels onto a problem a mutex solves in three lines.
- A good test: are you moving ownership of data, or guarding it in place? Moving suggests a channel; guarding suggests a mutex.
- The standard library itself uses mutexes heavily. Using one is not un-idiomatic.

## 9. Concurrency bugs

Three bugs account for nearly every concurrency failure in Go: data races, deadlocks and goroutine leaks.

### Data race

```go
count := 0
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        count++ // read, add, write: three steps that can interleave
    }()
}
wg.Wait()
fmt.Println(count) // often less than 1000
```

A data race is two goroutines accessing the same variable at the same time, at least one of them writing, with no synchronization. Fix it with a mutex, an atomic, or by giving the variable to a single goroutine.

### Deadlock

```go
// 1. Sending with no receiver
func main() {
    ch := make(chan int)
    ch <- 1 // fatal error: all goroutines are asleep - deadlock!
}

// 2. Lock ordering: each goroutine holds one lock and waits for the other
func transfer(from, to *Account, amount int) {
    from.mu.Lock()
    defer from.mu.Unlock()
    to.mu.Lock() // transfer(a, b) and transfer(b, a) together can deadlock
    defer to.mu.Unlock()
    from.balance -= amount
    to.balance += amount
}

// Fix: always lock in one global order, for example by account ID
first, second := from, to
if to.id < from.id {
    first, second = to, from
}
first.mu.Lock()
defer first.mu.Unlock()
second.mu.Lock()
defer second.mu.Unlock()
```

### Goroutine leak

```go
func search(query string) string {
    ch := make(chan string) // unbuffered
    go func() { ch <- slowBackend(query) }()

    select {
    case r := <-ch:
        return r
    case <-time.After(time.Second):
        return "timeout" // nobody will ever receive, so the goroutine blocks forever
    }
}

// Fix: make(chan string, 1), so the send always completes and the goroutine exits
```

| Bug | Typical causes | How to find it |
| --- | --- | --- |
| Data race | Unprotected shared variable, map written from several goroutines, captured variable | `go test -race`, `go run -race` |
| Deadlock | No receiver, missing `wg.Done()`, lock ordering, locking twice, range over a channel nobody closes | The runtime reports it only when every goroutine is blocked; otherwise a goroutine dump |
| Goroutine leak | Blocked send or receive, missing cancellation, a ticker never stopped | The `goroutine` profile in pprof; the `goroutineleak` profile since Go 1.27 |

- The race detector finds races only on code paths that actually run, and it slows the program several times over. Use it in tests and CI.
- A partial deadlock, where some goroutines are stuck but others still run, is silent. The program just stops making progress.
- Go 1.27 made a dedicated leak profile generally available at `/debug/pprof/goroutineleak`, which reports goroutines blocked on something that can never unblock ([Go 1.27 release notes](https://go.dev/doc/go1.27)).
- The rule that prevents leaks: never start a goroutine without knowing how it will stop.

## 10. The scheduler: G, M and P

The Go runtime multiplexes many goroutines onto a small number of OS threads using three entities, known as the G-M-P model.

| Piece | What it is |
| --- | --- |
| G | A goroutine: its stack, instruction pointer and status |
| M | A machine: an OS thread that executes code |
| P | A processor: a scheduling context with a local queue of runnable goroutines; there are GOMAXPROCS of them |

- An M must hold a P to run Go code, so at most GOMAXPROCS goroutines execute Go code at the same instant.
- Each P has a local run queue. There is also one global queue, which Ps check periodically.
- Work stealing: when a P's queue is empty, it takes half the goroutines from another P's queue.
- GOMAXPROCS defaults to the number of CPU cores. Since Go 1.25 it also respects a container's CPU limit on Linux.

### What happens when a goroutine blocks

| The goroutine blocks on | What the runtime does |
| --- | --- |
| A channel, mutex or `select` | Parks the goroutine; the thread picks up another one. No thread is blocked |
| Network I/O | Parks the goroutine on the netpoller (epoll or kqueue) until the socket is ready |
| `time.Sleep` or a timer | Parks the goroutine until the timer fires |
| A blocking system call, such as file I/O | The thread blocks; its P is handed to another thread so other goroutines keep running |

### Preemption

- Since Go 1.14 preemption is asynchronous: a background monitor signals any goroutine that has run for about 10 ms, even in a tight loop with no function calls.
- Before 1.14 a goroutine could be preempted only at function calls, so an empty `for {}` loop could freeze a program.
- `runtime.Gosched()` voluntarily yields the processor. You rarely need it.

### Why goroutines are cheap

- A goroutine starts with a stack of about 2 KB that grows and shrinks on demand by copying.
- Switching between goroutines happens in user space and saves only a few registers, with no kernel call.
- Parking a goroutine costs almost nothing, which is why blocking code in Go scales without callbacks or async/await.

## 11. The memory model: happens-before

The Go memory model says a goroutine is guaranteed to see another goroutine's write only if a synchronization event orders the write before the read.

```go
// Broken: no synchronization between the write and the read
var done bool

func main() {
    go func() { done = true }()
    for !done {
    } // may loop forever: the compiler is allowed to assume done never changes
}

// Correct: an atomic, a mutex or a channel creates the ordering
var done atomic.Bool

func main() {
    go func() { done.Store(true) }()
    for !done.Load() {
    }
}
```

### The guarantees worth remembering

| Event | Guarantee |
| --- | --- |
| `go f()` | The `go` statement happens before `f` starts running |
| Channel send | A send happens before the matching receive completes |
| Unbuffered channel | The receive happens before the send completes |
| `close(ch)` | The close happens before a receive that returns because the channel is closed |
| Mutex | An `Unlock` happens before the next `Lock` returns |
| `once.Do(f)` | The single call of `f` finishes before any `Do` returns |
| `wg.Wait()` | Returns after the `Done` calls that brought the counter to zero |

- "It works on my machine" proves nothing. A racy program can pass a million runs and fail on another CPU or compiler version.
- There is no such thing as a benign data race in Go. Even a single-word flag needs an atomic or a lock.
- You do not need to memorize the model. Follow the rule: if two goroutines touch a variable and one writes, synchronize.

## 12. Testing concurrent code

Run every test with the race detector, and test time-dependent code with a fake clock so tests are fast and never flaky.

```go
// Hammer shared state from many goroutines, then run: go test -race
func TestCounterConcurrent(t *testing.T) {
    var c Counter
    var wg sync.WaitGroup

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            c.Inc()
        }()
    }
    wg.Wait()

    if got := c.Value(); got != 100 {
        t.Fatalf("got %d, want 100", got)
    }
}
```

### Fake time with testing/synctest (Go 1.25)

```go
func TestTimeout(t *testing.T) {
    synctest.Test(t, func(t *testing.T) {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()

        <-ctx.Done() // returns at once: inside the bubble the clock is simulated

        if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
            t.Fatalf("got %v", ctx.Err())
        }
    })
}
```

| Command or tool | Purpose |
| --- | --- |
| `go test -race ./...` | Detect data races while the tests run |
| `go test -count=100 -race` | Repeat a test to shake out rare interleavings |
| `t.Parallel()` | Run tests concurrently, which also exposes shared-state bugs |
| `testing/synctest` | Run goroutines in a bubble with a simulated clock |
| `go.uber.org/goleak` | Fail a test if goroutines are still running at the end |
| `go tool pprof` with the goroutine profile | See where goroutines are blocked |

- Never use `time.Sleep` in a test to wait for a goroutine. Wait on a channel or a `WaitGroup`.
- Inside a `synctest` bubble, time moves forward only when every goroutine is blocked, so a five-second timeout takes no real time.
- A passing test without `-race` says little about concurrent code. Make `-race` part of CI.

## 13. Predict the output

Cover the answer comments and work out each of these five before reading on.

### Puzzle 1: a closed, buffered channel

```go
func main() {
    ch := make(chan int, 2)
    ch <- 1
    ch <- 2
    close(ch)

    for v := range ch {
        fmt.Print(v, " ")
    }
    v, ok := <-ch
    fmt.Println(v, ok)
}
// 1 2 0 false
```

Closing does not discard buffered values. Receivers drain them first, then get the zero value with `ok` false.

### Puzzle 2: a WaitGroup passed by value

```go
func worker(wg sync.WaitGroup) {
    defer wg.Done()
    fmt.Println("working")
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go worker(wg)
    wg.Wait()
}
// working
// fatal error: all goroutines are asleep - deadlock!
```

The worker calls `Done` on its own copy, so the original counter never reaches zero. Pass `*sync.WaitGroup`.

### Puzzle 3: two ready cases

```go
func main() {
    a, b := make(chan string, 1), make(chan string, 1)
    a <- "from a"
    b <- "from b"

    select {
    case v := <-a:
        fmt.Println(v)
    case v := <-b:
        fmt.Println(v)
    }
}
// "from a" or "from b": select chooses at random among ready cases
```

### Puzzle 4: goroutines in a loop

```go
func main() {
    var wg sync.WaitGroup
    for i := 0; i < 3; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            fmt.Print(i, " ")
        }()
    }
    wg.Wait()
}
// Go 1.22+: 0, 1 and 2 in some order, such as "2 0 1"
// Before 1.22: usually "3 3 3", because all goroutines shared one i
```

### Puzzle 5: locking twice

```go
var mu sync.Mutex

func outer() {
    mu.Lock()
    defer mu.Unlock()
    inner()
}

func inner() {
    mu.Lock() // the same goroutine already holds mu
    defer mu.Unlock()
    fmt.Println("inner")
}

func main() { outer() }
// fatal error: all goroutines are asleep - deadlock!
```

Go mutexes are not reentrant. Split `inner` into a public method that locks and a private helper that assumes the lock is held.

## 14. Interview coding problems

These four programs, together with the worker pool and pipeline in section 6, cover most live-coding concurrency rounds.

### Print odd and even numbers alternately with two goroutines

```go
func main() {
    oddTurn, evenTurn := make(chan struct{}), make(chan struct{})
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        for i := 1; i <= 9; i += 2 {
            <-oddTurn // wait for my turn
            fmt.Println(i)
            evenTurn <- struct{}{} // hand over
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
// 1 2 3 4 5 6 7 8 9 10, one per line
```

### Fetch in parallel with a limit, keeping results in order

```go
func fetchAll(ctx context.Context, urls []string) ([]string, error) {
    results := make([]string, len(urls))

    g, ctx := errgroup.WithContext(ctx)
    g.SetLimit(5) // at most 5 requests in flight

    for i, url := range urls {
        g.Go(func() error {
            body, err := fetch(ctx, url)
            if err != nil {
                return fmt.Errorf("fetch %s: %w", url, err)
            }
            results[i] = body // each goroutine writes its own index, so no lock is needed
            return nil
        })
    }

    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

### Run a function with a timeout

```go
func callWithTimeout(ctx context.Context, d time.Duration, fn func(context.Context) (string, error)) (string, error) {
    ctx, cancel := context.WithTimeout(ctx, d)
    defer cancel()

    type result struct {
        val string
        err error
    }
    ch := make(chan result, 1) // buffered: the goroutine can always send and exit

    go func() {
        v, err := fn(ctx)
        ch <- result{v, err}
    }()

    select {
    case r := <-ch:
        return r.val, r.err
    case <-ctx.Done():
        return "", ctx.Err() // context.DeadlineExceeded
    }
}
```

### Workers with graceful shutdown

```go
func main() {
    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer stop()

    jobs := make(chan int)
    var wg sync.WaitGroup

    for w := 0; w < 3; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := range jobs { // finishes the job in hand, then sees the close
                process(j)
            }
        }()
    }

produce:
    for i := 0; ; i++ {
        select {
        case <-ctx.Done(): // Ctrl+C or SIGTERM
            break produce
        case jobs <- i:
        }
    }

    close(jobs) // stop handing out work
    wg.Wait()   // let in-flight jobs finish before exiting
    fmt.Println("shut down cleanly")
}
```

- For odd and even, a follow-up is to generalize to N goroutines printing in turn: use a ring of N channels.
- The timeout wrapper cancels the context, but `fn` must honour `ctx` for the work itself to stop.
- For graceful shutdown, mention a deadline on the wait so a stuck job cannot block the exit forever.
- Explain your shutdown path out loud in every answer: who closes which channel, and how each goroutine ends.

## 15. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| What is a goroutine? | A lightweight function execution scheduled by the Go runtime. |
| Goroutine vs thread? | About 2 KB of growable stack and user-space scheduling, against megabytes and the kernel. |
| Concurrency vs parallelism? | Structuring independent tasks; running them at the same instant. |
| What happens when main returns? | The program exits and all other goroutines are killed. |
| Buffered vs unbuffered channel? | An unbuffered send waits for a receiver; a buffered send waits only when full. |
| Send on a closed channel? | Panics. |
| Receive from a closed channel? | Returns remaining values, then the zero value with ok false. |
| Send or receive on a nil channel? | Blocks forever. |
| Who should close a channel? | The sender, and only when receivers need the signal. |
| What does select do with several ready cases? | Picks one at random. |
| How do you make a channel operation non-blocking? | A `select` with a `default` case. |
| Mutex vs RWMutex? | One holder at a time; many readers or one writer. |
| Is a Go mutex reentrant? | No; locking it twice in one goroutine deadlocks. |
| What does WaitGroup do? | Waits until a counted set of goroutines has finished. |
| What is `sync.Once` for? | Running initialization exactly once. |
| When do you use atomic? | For a single counter, flag or pointer. |
| What is context for? | Cancellation, deadlines and request-scoped values. |
| What is a data race? | Unsynchronized concurrent access to a variable with at least one write. |
| How do you detect races? | `go test -race` or `go run -race`. |
| What is a goroutine leak? | A goroutine blocked forever with no way to finish. |
| What is the G-M-P model? | Goroutines run on OS threads through processors that hold run queues. |
| What does GOMAXPROCS set? | The number of Ps, so how many goroutines run Go code in parallel. |
| What is work stealing? | An idle P taking half the queue of a busy one. |
| Channel or mutex? | A channel to pass data or signals; a mutex to guard shared state. |
| What does errgroup add over WaitGroup? | Error propagation, cancellation on first failure, and a concurrency limit. |
