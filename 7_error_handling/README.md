# Error Handling in Go

Oct 11, 2026 · @HR

## 1. Errors are values

Go has no exceptions: a function that can fail returns an `error` as its last result, and the caller checks it straight away.

```go
// The whole definition of the built-in error type
type error interface {
    Error() string
}

func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil // nil means success
}

result, err := divide(10, 0)
if err != nil {
    fmt.Println("error:", err) // error: division by zero
    return
}
fmt.Println(result)
```

### The idiomatic shape

```go
func loadUser(id int) (*User, error) {
    row, err := queryUser(id)
    if err != nil {
        return nil, fmt.Errorf("load user %d: %w", id, err) // fail fast, add context
    }

    user, err := parseUser(row)
    if err != nil {
        return nil, fmt.Errorf("load user %d: %w", id, err)
    }

    return user, nil // the happy path stays at the left edge
}
```

- `error` is an ordinary interface, so any type with an `Error() string` method is an error.
- The error is the last return value by convention, and nil means no error.
- Check the error before touching the other results. On failure they are usually zero values.
- Handle the error case first and return early, so the success path is not nested inside `else` blocks.
- Explicit checks make every failure point visible in the code, which is the design goal.

## 2. Creating errors

Use `errors.New` for a fixed message and `fmt.Errorf` when the message needs values or must wrap another error.

```go
err1 := errors.New("connection refused")          // fixed text
err2 := fmt.Errorf("user %d not found", id)       // formatted text
err3 := fmt.Errorf("connect to db: %w", err1)     // wraps err1 and adds context

fmt.Println(err3) // connect to db: connection refused
```

### Each call creates a distinct error

```go
a := errors.New("not found")
b := errors.New("not found")
fmt.Println(a == b) // false: same text, different errors
```

Errors are compared by identity, not by message. That is why a shared error is declared once as a package-level variable.

| Convention | Example |
| --- | --- |
| Start lowercase, no trailing punctuation | `errors.New("invalid token")` |
| Say what was being attempted | `fmt.Errorf("open config: %w", err)` |
| Include the values that identify the failure | `fmt.Errorf("user %d not found", id)` |
| Do not repeat words such as "error" or "failed" at each level | `parse header: unexpected EOF` |

- Messages are lowercase because they are usually joined into a longer chain: `start server: open config: file missing`.
- Never put secrets or full request bodies into an error message. Errors end up in logs.

## 3. Sentinel errors and errors.Is

A sentinel error is a package-level error value that callers compare against with `errors.Is`.

```go
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
)

func findUser(id int) (*User, error) {
    if id != 1 {
        return nil, ErrNotFound
    }
    return &User{ID: 1, Name: "Asha"}, nil
}

_, err := findUser(2)
if errors.Is(err, ErrNotFound) {
    fmt.Println("respond with 404")
}
```

### Why errors.Is and not ==

```go
wrapped := fmt.Errorf("find user 2: %w", ErrNotFound)

fmt.Println(wrapped == ErrNotFound)          // false: it is a different, outer error
fmt.Println(errors.Is(wrapped, ErrNotFound)) // true: Is looks through the wrapping
```

### Sentinels in the standard library

| Sentinel | Meaning |
| --- | --- |
| `io.EOF` | The input has ended; expected, not a failure |
| `sql.ErrNoRows` | A query for one row found none |
| `os.ErrNotExist` | The file or directory does not exist |
| `context.Canceled` | The context was cancelled |
| `context.DeadlineExceeded` | The context's deadline passed |
| `http.ErrServerClosed` | The server stopped because `Shutdown` or `Close` was called |

```go
data, err := os.ReadFile("config.yaml")
if errors.Is(err, os.ErrNotExist) {
    data = defaultConfig // a missing file is fine: fall back
} else if err != nil {
    return err
}
```

- Name sentinels with an `Err` prefix and export them, so callers can test for them.
- A sentinel carries no detail. When the caller needs fields such as an ID, use a custom error type.
- Sentinels become part of your package's public API, so add them sparingly.

## 4. Custom error types and errors.As

When the caller needs details about a failure, define a struct that implements `error` and extract it with `errors.As`.

```go
type ValidationError struct {
    Field  string
    Reason string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("invalid %s: %s", e.Field, e.Reason)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{Field: "age", Reason: "must not be negative"}
    }
    return nil
}

err := fmt.Errorf("create user: %w", validateAge(-1)) // wrapped by a caller

var ve *ValidationError
if errors.As(err, &ve) { // finds the first *ValidationError in the chain
    fmt.Println(ve.Field, "-", ve.Reason) // age - must not be negative
}
```

### Error types in the standard library

```go
_, err := os.Open("/no/such/file")

var pe *fs.PathError
if errors.As(err, &pe) {
    fmt.Println(pe.Op, pe.Path) // open /no/such/file
}

_, err = strconv.Atoi("12x")

var ne *strconv.NumError
if errors.As(err, &ne) {
    fmt.Println(ne.Func, ne.Num) // Atoi 12x
}
```

| Need | Use |
| --- | --- |
| Only "did this specific thing happen?" | A sentinel with `errors.Is` |
| Structured details: a field name, a status code, a path | A custom type with `errors.As` |
| Just a message for a human or a log | `errors.New` or `fmt.Errorf` |

- Use a pointer receiver for `Error` and return a pointer. Then the target for `errors.As` is a variable of that pointer type.
- The second argument to `errors.As` must be a non-nil pointer; anything else panics.
- The type must match exactly. If the error is a `*ValidationError`, a target of type `ValidationError` finds nothing.
- Before Go 1.13 this was written as a type assertion, `err.(*ValidationError)`, which fails on wrapped errors.

## 5. Wrapping errors

Wrapping with `%w` adds context to an error while keeping the original reachable, so each layer can say what it was doing.

```go
func readConfig(path string) ([]byte, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, fmt.Errorf("read config %q: %w", path, err)
    }
    return data, nil
}

func start() error {
    if _, err := readConfig("app.yaml"); err != nil {
        return fmt.Errorf("start server: %w", err)
    }
    return nil
}

err := start()
fmt.Println(err)
// start server: read config "app.yaml": open app.yaml: no such file or directory

fmt.Println(errors.Is(err, os.ErrNotExist)) // true: the original is still in the chain
```

| Verb | Result | Effect on `errors.Is` and `errors.As` |
| --- | --- | --- |
| `%w` | Wraps the error | They can still find the inner error |
| `%v` | Copies only its text | The chain is cut; the inner error is hidden |

### A custom type that wraps

```go
type QueryError struct {
    Query string
    Err   error
}

func (e *QueryError) Error() string { return "query " + e.Query + ": " + e.Err.Error() }
func (e *QueryError) Unwrap() error { return e.Err } // joins the chain

err := &QueryError{Query: "SELECT 1", Err: context.DeadlineExceeded}
fmt.Println(errors.Is(err, context.DeadlineExceeded)) // true, thanks to Unwrap
```

### Combining several errors

```go
func validate(u User) error {
    var errs []error
    if u.Name == "" {
        errs = append(errs, errors.New("name is required"))
    }
    if u.Age < 0 {
        errs = append(errs, errors.New("age must not be negative"))
    }
    return errors.Join(errs...) // nil when there are no errors
}

fmt.Println(validate(User{Age: -1}))
// name is required
// age must not be negative
```

- `errors.Unwrap(err)` returns the next error down, or nil. You rarely call it yourself; `Is` and `As` do the walking.
- `errors.Join` (Go 1.20) keeps every error, and `errors.Is` matches any of them. `fmt.Errorf` also accepts several `%w` verbs since 1.20.
- Wrapping makes the inner error part of your API, because callers may start testing for it.
- Use `%v` when the inner error is an implementation detail you do not want callers to depend on.
- Add context that the lower level does not know: which user, which file, which step.

## 6. errors.Is vs errors.As

`errors.Is` asks whether a specific error value is in the chain, and `errors.As` asks whether an error of a given type is in the chain.

|  | `errors.Is` | `errors.As` |
| --- | --- | --- |
| Question | Is this exact error somewhere in the chain? | Is there an error of this type in the chain? |
| Matches by | Value, using `==` | Type |
| Second argument | An error value, usually a sentinel | A pointer to a variable of the target type |
| Gives you | true or false | true or false, plus the variable filled in |
| Typical target | `io.EOF`, `sql.ErrNoRows`, `ErrNotFound` | `*fs.PathError`, `*ValidationError` |

```go
err := fmt.Errorf("save: %w", &ValidationError{Field: "email", Reason: "empty"})

// Is compares values, and two separate pointers are never equal
fmt.Println(errors.Is(err, &ValidationError{Field: "email", Reason: "empty"})) // false

// As compares types, which is what you want for a struct error
var ve *ValidationError
fmt.Println(errors.As(err, &ve)) // true
```

### Customizing a match with an Is method

```go
var ErrTemporary = errors.New("temporary failure")

type TimeoutError struct{ Op string }

func (e *TimeoutError) Error() string { return e.Op + ": timed out" }

// Declare that every TimeoutError counts as temporary
func (e *TimeoutError) Is(target error) bool { return target == ErrTemporary }

err := fmt.Errorf("fetch: %w", &TimeoutError{Op: "read"})
fmt.Println(errors.Is(err, ErrTemporary)) // true
```

- Both functions walk the chain by calling `Unwrap` repeatedly until they find a match or reach the end.
- Both return false for a nil error, except that `errors.Is(nil, nil)` is true.
- Rule of thumb: `Is` for sentinels, `As` for types.
- Never use `==` or a type assertion on an error that may have been wrapped.

## 7. Handling strategies

Handle each error exactly once: either deal with it where you are, or add context and return it, but never both.

| Option | When to choose it |
| --- | --- |
| Return it wrapped with context | The default: you cannot fix it, but you know what you were doing |
| Return it unchanged | You have nothing useful to add |
| Recover with a fallback | A sensible default exists, such as a missing optional config file |
| Retry | The failure is temporary, such as a timeout |
| Log and continue | The work is best-effort, such as sending a metric |
| Panic | A bug or a startup failure the program cannot run without |

### Log once, at the top

```go
// Problem: every layer logs, so one failure appears three times
func loadUser(id int) (*User, error) {
    u, err := repo.Find(id)
    if err != nil {
        log.Println("find failed:", err)
        return nil, err
    }
    return u, nil
}

// Better: add context and return; only the HTTP handler or main logs it
func loadUser(id int) (*User, error) {
    u, err := repo.Find(id)
    if err != nil {
        return nil, fmt.Errorf("load user %d: %w", id, err)
    }
    return u, nil
}
```

### Keep err in a narrow scope

```go
if err := validate(input); err != nil { // err exists only inside this if
    return err
}
```

### Do not lose the error from Close

```go
func writeFile(path string, data []byte) (err error) {
    f, err := os.Create(path)
    if err != nil {
        return err
    }
    defer func() {
        err = errors.Join(err, f.Close()) // a failed Close can mean lost data
    }()

    _, err = f.Write(data)
    return err
}
```

### Ignoring an error on purpose

```go
_ = logFile.Close() // the blank identifier shows the choice was deliberate
```

- Never drop an error silently. If you ignore one, write `_ =` and ideally a comment saying why.
- When a function returns an error, treat its other results as invalid unless its documentation says otherwise.
- For a file opened only for reading, a plain `defer f.Close()` is fine. For writes, check the result.
- Linters such as `errcheck` find unchecked errors automatically.

## 8. Errors in goroutines

A goroutine cannot return a value to whoever started it, so its error must travel through an errgroup, a channel, or a shared variable.

### Stop at the first error: errgroup

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) error {
    g, ctx := errgroup.WithContext(ctx)
    for _, url := range urls {
        g.Go(func() error {
            return fetch(ctx, url) // the first error cancels ctx for the others
        })
    }
    return g.Wait() // waits for all, returns the first non-nil error
}
```

### Collect every error

```go
func processAll(jobs []Job) error {
    var (
        mu   sync.Mutex
        errs []error
        wg   sync.WaitGroup
    )
    for _, job := range jobs {
        wg.Add(1)
        go func() {
            defer wg.Done()
            if err := process(job); err != nil {
                mu.Lock()
                errs = append(errs, fmt.Errorf("job %d: %w", job.ID, err))
                mu.Unlock()
            }
        }()
    }
    wg.Wait()
    return errors.Join(errs...)
}
```

### Send errors over a channel

```go
errCh := make(chan error, len(jobs)) // buffered, so no goroutine blocks on send
for _, job := range jobs {
    go func() { errCh <- process(job) }()
}

for range jobs {
    if err := <-errCh; err != nil {
        log.Println("job failed:", err)
    }
}
```

- An error that a goroutine ignores is lost for good; nothing reports it.
- `errgroup` is the standard answer when asked how to run tasks in parallel and fail fast.
- A panic inside a goroutine is not an error value. It crashes the program unless that goroutine recovers it.
- When a context is cancelled, return `ctx.Err()` so callers can test for `context.Canceled`.

## 9. Errors across layers and in HTTP APIs

Each layer should translate low-level errors into its own vocabulary, and only the outermost layer should turn them into status codes.

```go
// Domain errors: the vocabulary the whole service shares
var (
    ErrNotFound = errors.New("not found")
    ErrConflict = errors.New("already exists")
)

// Repository: translate the driver's error into a domain error
func (r *UserRepo) Get(ctx context.Context, id int) (*User, error) {
    var u User
    err := r.db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id = $1", id).
        Scan(&u.ID, &u.Name)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrNotFound // callers never need to import database/sql
    }
    if err != nil {
        return nil, fmt.Errorf("get user %d: %w", id, err)
    }
    return &u, nil
}

// Handler: the one place that maps errors to HTTP status codes
func statusFor(err error) int {
    var ve *ValidationError
    switch {
    case errors.Is(err, ErrNotFound):
        return http.StatusNotFound
    case errors.Is(err, ErrConflict):
        return http.StatusConflict
    case errors.As(err, &ve):
        return http.StatusBadRequest
    case errors.Is(err, context.DeadlineExceeded):
        return http.StatusGatewayTimeout
    default:
        return http.StatusInternalServerError
    }
}

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
    user, err := h.service.User(r.Context(), 42)
    if err != nil {
        status := statusFor(err)
        if status == http.StatusInternalServerError {
            log.Printf("get user: %v", err) // full detail goes to the log
            http.Error(w, "internal server error", status) // a generic message goes to the client
            return
        }
        http.Error(w, err.Error(), status)
        return
    }
    json.NewEncoder(w).Encode(user)
}
```

| Layer | Responsibility |
| --- | --- |
| Repository | Convert driver errors to domain errors; wrap the rest with context |
| Service | Add business context; create validation errors |
| Handler | Map to a status code, log unexpected errors, shape the response |

- Never send raw internal errors to clients. They can leak SQL, file paths and hostnames.
- Because every layer wraps with `%w`, the handler can still recognize `ErrNotFound` three levels up.
- Keeping the status mapping in one function makes the API's behaviour consistent and easy to test.

## 10. Testing errors

Test error paths as carefully as success paths, and assert with `errors.Is` or `errors.As` rather than by comparing message text.

```go
func TestFindUser(t *testing.T) {
    tests := []struct {
        name    string
        id      int
        wantErr error
    }{
        {"found", 1, nil},
        {"missing", 2, ErrNotFound},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := findUser(tt.id)
            if !errors.Is(err, tt.wantErr) { // also true when both are nil
                t.Fatalf("got error %v, want %v", err, tt.wantErr)
            }
        })
    }
}

func TestValidateAge(t *testing.T) {
    err := validateAge(-1)

    var ve *ValidationError
    if !errors.As(err, &ve) {
        t.Fatalf("got %T, want *ValidationError", err)
    }
    if ve.Field != "age" {
        t.Errorf("field = %q, want age", ve.Field)
    }
}
```

- A message is for humans and may be reworded at any time. A test that matches on text breaks for no real reason.
- Use a fake dependency that returns a chosen error to exercise the failure branches of your code.
- Check both directions: the right error on bad input, and no error on good input.
- When a test only needs to know that something failed, a `wantErr bool` column is enough.

## 11. Gotchas

These five mistakes come up again and again in code reviews and in "what is wrong here?" questions.

### A nil pointer returned as an error is not nil

```go
type MyErr struct{}

func (*MyErr) Error() string { return "boom" }

func do() error {
    var e *MyErr // nil pointer
    return e     // becomes a non-nil interface holding a nil pointer
}

fmt.Println(do() == nil) // false
```

Return a literal `nil` on success, and declare error variables as type `error`, not as a concrete pointer type.

### Shadowing err with :=

```go
func load(useCache bool) error {
    var err error
    if !useCache {
        data, err := fetch() // := declares a NEW err inside this block
        _ = data
    }
    return err // the outer err: always nil, even when fetch failed
}
```

Declare `data` beforehand and assign with `=`, or return from inside the block.

### Using a result before checking the error

```go
resp, err := http.Get(url)
defer resp.Body.Close() // panic when err != nil, because resp is nil
if err != nil {
    return err
}

// Correct order: check first, then defer
resp, err := http.Get(url)
if err != nil {
    return err
}
defer resp.Body.Close()
```

### Comparing with == after wrapping

```go
err := fmt.Errorf("query: %w", sql.ErrNoRows)

if err == sql.ErrNoRows { } // false: never matches a wrapped error
if errors.Is(err, sql.ErrNoRows) { } // true
```

### Matching on the message text

```go
if strings.Contains(err.Error(), "not found") { } // brittle: breaks when wording changes
if errors.Is(err, ErrNotFound) { }                // robust
```

- `go vet` flags some unused results but not shadowing by default. The shadow check in golangci-lint catches shadowed variables.
- Treating `io.EOF` as a failure is another common slip. In a read loop it is the normal way to finish.
- Wrapping with `%v` instead of `%w` by accident silently breaks every `errors.Is` check above it.

## 12. Interview coding problems

These four exercises cover what interviewers ask you to build: an application error type, a central HTTP error handler, a selective retry, and a panic-to-error wrapper.

### An application error with a code and a cause

```go
type AppError struct {
    Code    int    // HTTP status to send
    Message string // safe to show to the client
    Err     error  // underlying cause, for logs
}

func (e *AppError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Err)
    }
    return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

func NotFound(msg string, cause error) error {
    return &AppError{Code: http.StatusNotFound, Message: msg, Err: cause}
}
```

### A central error handler for HTTP

```go
type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func handle(h HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        err := h(w, r)
        if err == nil {
            return
        }

        var app *AppError
        if errors.As(err, &app) {
            http.Error(w, app.Message, app.Code)
            return
        }

        log.Printf("unhandled error: %v", err)
        http.Error(w, "internal server error", http.StatusInternalServerError)
    }
}

// Handlers now simply return errors
func getUser(w http.ResponseWriter, r *http.Request) error {
    user, err := service.User(r.Context(), 42)
    if errors.Is(err, ErrNotFound) {
        return NotFound("user not found", err)
    }
    if err != nil {
        return err
    }
    return json.NewEncoder(w).Encode(user)
}

mux.HandleFunc("GET /users/{id}", handle(getUser))
```

### Retry only temporary failures

```go
var ErrTemporary = errors.New("temporary failure")

func retry(attempts int, fn func() error) error {
    var err error
    for i := 0; i < attempts; i++ {
        if err = fn(); err == nil {
            return nil
        }
        if !errors.Is(err, ErrTemporary) {
            return err // permanent: retrying will not help
        }
        time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
    }
    return fmt.Errorf("gave up after %d attempts: %w", attempts, err)
}
```

### Turn a panic into an error

```go
func safeCall(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()
    fn()
    return nil
}

err := safeCall(func() {
    var m map[string]int
    m["x"] = 1 // panics: write to a nil map
})
fmt.Println(err) // recovered from panic: assignment to entry in nil map
```

- `NotFound` returns the `error` interface type, not `*AppError`, which avoids the typed-nil trap.
- For the retry, mention the usual additions: a `context.Context` for cancellation, exponential backoff and jitter.
- `safeCall` needs a named result, because that is the only way a deferred function can set the return value.

## 13. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| What is the `error` type? | A built-in interface with one method, `Error() string`. |
| Does Go have exceptions? | No; functions return errors as values. |
| Why no try-catch? | Explicit returns keep every failure path visible at the call site. |
| `errors.New` vs `fmt.Errorf`? | A fixed message; a formatted message that can also wrap. |
| What is a sentinel error? | A package-level error value compared with `errors.Is`, such as `io.EOF`. |
| What does `%w` do? | Wraps an error so it stays reachable through the chain. |
| `%w` vs `%v`? | `%w` keeps the chain; `%v` keeps only the text. |
| `errors.Is` vs `errors.As`? | `Is` matches a specific value; `As` matches a type and extracts it. |
| Why not `==` on errors? | It fails as soon as the error has been wrapped. |
| How does a custom type join the chain? | By implementing `Unwrap() error`. |
| What does `errors.Join` do? | Combines several errors into one; returns nil if there are none. |
| When should you use a custom error type? | When callers need structured details, not just identity. |
| Is `io.EOF` a failure? | No; it signals the normal end of input. |
| Should you log an error and also return it? | No; handle it once to avoid duplicate logs. |
| Why can a returned error be non-nil unexpectedly? | A nil pointer of a concrete type was stored in the interface. |
| How do goroutines report errors? | Through an errgroup, a channel, or a mutex-guarded slice. |
| Error or panic? | Errors for expected failures; panic for bugs and impossible states. |
| How should error strings be written? | Lowercase, no trailing punctuation, with context. |
| How do you map errors to HTTP status codes? | In one place, using `errors.Is` and `errors.As`. |
| How do you assert errors in tests? | With `errors.Is` or `errors.As`, never by message text. |
