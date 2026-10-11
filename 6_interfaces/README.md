# Interfaces in Go

Oct 11, 2026 · @HR

## 1. What an interface is

An interface is a set of method signatures, and any type that has those methods satisfies it automatically, with no `implements` keyword.

```go
type Notifier interface {
    Notify(msg string) error
}

type Email struct{ To string }

func (e Email) Notify(msg string) error {
    fmt.Println("email to", e.To, ":", msg)
    return nil
}

type SMS struct{ Phone string }

func (s SMS) Notify(msg string) error {
    fmt.Println("sms to", s.Phone, ":", msg)
    return nil
}

// alert works with any Notifier, including types written later
func alert(n Notifier) {
    _ = n.Notify("server down")
}

alert(Email{To: "ops@x.com"})
alert(SMS{Phone: "98765"})

channels := []Notifier{Email{"a@x.com"}, SMS{"12345"}} // mixed types in one slice
for _, c := range channels {
    _ = c.Notify("deploy finished")
}
```

- Satisfaction is implicit and checked at compile time. `Email` never mentions `Notifier`.
- This lets you define an interface for a type you do not own, including standard library types.
- An interface describes behaviour only. It has no fields and no implementation.
- Interfaces are how Go does polymorphism; there is no inheritance.
- A type can satisfy many interfaces at once, and an interface can be satisfied by many types.

## 2. Interface values

An interface value is a pair: the concrete type currently stored in it, called the dynamic type, and the value of that type.

```go
var w io.Writer             // nil interface: no type, no value
fmt.Printf("%T %v\n", w, w) // <nil> <nil>

w = os.Stdout          // now holds (type *os.File, value the pointer)
fmt.Printf("%T\n", w)  // *os.File

w = &bytes.Buffer{}    // the same variable can hold a different type later
fmt.Printf("%T\n", w)  // *bytes.Buffer
```

| State | Dynamic type | Dynamic value | `== nil` |
| --- | --- | --- | --- |
| Declared, never assigned | none | none | true |
| Holds `os.Stdout` | `*os.File` | a valid pointer | false |
| Holds a nil `*bytes.Buffer` | `*bytes.Buffer` | nil | false |

### An interface stores a copy

```go
type Getter interface{ Get() int }

type Point struct{ X int }

func (p Point) Get() int { return p.X }

p := Point{1}
var g Getter = p // the interface holds its own copy of p
p.X = 2
fmt.Println(g.Get()) // 1

g = &p // store a pointer and the interface follows the original
p.X = 3
fmt.Println(g.Get()) // 3
```

- The static type of `w` is always `io.Writer`. Its dynamic type changes with each assignment.
- Through an interface you can call only the methods the interface lists, even if the concrete type has more.
- Calling a method on a nil interface panics, because there is no type to dispatch to.
- `%T` prints the dynamic type, which is the quickest way to debug an interface value.

## 3. The empty interface and any

The empty interface has no methods, so every type satisfies it; `any` is simply an alias for `interface{}` (Go 1.18).

```go
var x any = 42
x = "hello"
x = []int{1, 2, 3} // any value fits

func describe(v any) {
    fmt.Printf("%T: %v\n", v, v)
}

describe(3.14)  // float64: 3.14
describe("go")  // string: go

// var n int = x // compile error: you must assert the type first
// x + 1         // compile error: any supports no operations
```

### A slice of any is not a slice of anything

```go
names := []string{"a", "b"}
// var xs []any = names // compile error: different memory layouts

xs := make([]any, len(names)) // convert element by element
for i, n := range names {
    xs[i] = n
}
```

### Dynamic JSON

```go
var data map[string]any
json.Unmarshal([]byte(`{"name":"Asha","age":30,"tags":["go"]}`), &data)

fmt.Printf("%T\n", data["name"]) // string
fmt.Printf("%T\n", data["age"])  // float64: every JSON number decodes to float64
fmt.Printf("%T\n", data["tags"]) // []interface {}
```

- `any` gives up compile-time type safety, so every use needs a type assertion or a type switch.
- Legitimate uses: `fmt.Println`, JSON of unknown shape, and generic containers written before Go 1.18.
- In new code, prefer a specific interface or generics over `any`.
- Assigning a concrete value to `any` usually allocates, because the value is boxed.

## 4. Type assertions and type switches

A type assertion extracts the concrete value from an interface, and a type switch does the same across several possible types.

```go
var v any = "hello"

s := v.(string)   // hello: single-result form
n, ok := v.(int)  // 0 false: comma-ok form never panics
// _ = v.(int)    // panic: interface conversion: interface {} is string, not int

if s, ok := v.(string); ok {
    fmt.Println("string of length", len(s))
}
```

### Asserting to another interface

```go
func closeIfPossible(r io.Reader) {
    if c, ok := r.(io.Closer); ok { // does the dynamic type also have Close?
        c.Close()
    }
}
```

### Type switch

```go
func describe(v any) string {
    switch x := v.(type) {
    case nil:
        return "nil"
    case int, int64: // several types: x keeps the type any here
        return fmt.Sprintf("integer %v", x)
    case string: // one type: x is a string here
        return "string of length " + strconv.Itoa(len(x))
    case error:
        return "error: " + x.Error()
    case fmt.Stringer:
        return "stringer: " + x.String()
    default:
        return fmt.Sprintf("unknown type %T", x)
    }
}
```

- Use the comma-ok form whenever the type is not guaranteed. The single-result form panics on a mismatch.
- An assertion to a concrete type checks for an exact match with the dynamic type.
- An assertion to an interface type checks whether the dynamic type has the required methods.
- Cases in a type switch are tried top to bottom, and the first match wins. Put specific types before broad interfaces.
- An assertion on a nil interface always fails: `ok` is false, or the single-result form panics.

## 5. The nil interface trap

An interface is nil only when it holds no type and no value, so an interface that holds a nil pointer is not nil.

```go
type MyErr struct{}

func (*MyErr) Error() string { return "boom" }

func do() error {
    var e *MyErr // a nil pointer
    return e     // the interface now holds (type *MyErr, value nil)
}

func main() {
    err := do()
    fmt.Println(err == nil) // false
    fmt.Printf("%T\n", err) // *main.MyErr: the type is set, so the interface is not nil
}
```

### The same thing in two lines

```go
var p *bytes.Buffer   // nil pointer
var w io.Writer = p   // interface holding (type *bytes.Buffer, value nil)

fmt.Println(p == nil, w == nil) // true false
```

### The fix

```go
func do() error {
    var e *MyErr
    if failed {
        e = &MyErr{}
        return e
    }
    return nil // return a literal nil, never a typed nil pointer
}

// Simpler still: declare the variable as the interface type
func do2() error {
    var err error // nil interface
    // ...
    return err
}
```

- This is the most common Go interview trap. The answer to "why is it not nil?" is "the interface has a type".
- It bites hardest with `error`, because callers always write `if err != nil`.
- A function should declare its return type as `error`, not as a concrete pointer type such as `*MyErr`.
- Calling a method on such a value does not panic by itself. The method runs with a nil receiver.

## 6. Method sets: pointer vs value receivers

Whether a type satisfies an interface depends on its method set: a value has only the value-receiver methods, while a pointer has all of them.

| Methods declared with | Does `T` satisfy the interface? | Does `*T` satisfy it? |
| --- | --- | --- |
| Value receivers only | Yes | Yes |
| At least one pointer receiver | No | Yes |

```go
type Shape interface{ Area() float64 }

type Rect struct{ W, H float64 }

func (r Rect) Area() float64 { return r.W * r.H } // value receiver

type Square struct{ Side float64 }

func (s *Square) Area() float64 { return s.Side * s.Side } // pointer receiver

var a Shape = Rect{2, 3}  // ok
var b Shape = &Rect{2, 3} // ok: a pointer has the value methods too
var c Shape = &Square{2}  // ok
// var d Shape = Square{2} // compile error: Square does not implement Shape
//                         // (method Area has pointer receiver)
```

### Compile-time check

```go
var _ Shape = Rect{}         // fails to compile if Rect stops satisfying Shape
var _ Shape = (*Square)(nil) // the same check for a pointer type, with no allocation
```

- The reason: a value stored in an interface is not addressable, so Go cannot take its address to call a pointer method.
- On an ordinary variable, `sq.Area()` works with a pointer receiver because Go rewrites it to `(&sq).Area()`.
- If any method needs a pointer receiver, use pointer receivers for all and pass pointers around.

## 7. Composing interfaces

Large interfaces are built by embedding small ones, and a struct can embed an interface to wrap or partially implement it.

### Interface embedding

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface { // the union of both method sets
    Reader
    Writer
}
```

This is how the standard library defines `io.ReadWriter`, `io.ReadCloser` and `io.ReadWriteCloser`.

### A struct that embeds an interface

```go
type Store interface {
    Get(id int) (string, error)
    Save(id int, v string) error
}

// Decorator: override one method, and the embedded Store handles the rest
type LoggingStore struct {
    Store
}

func (l LoggingStore) Save(id int, v string) error {
    log.Println("saving", id)
    return l.Store.Save(id, v) // delegate to the wrapped implementation
}

var s Store = LoggingStore{Store: realStore} // LoggingStore is itself a Store
```

- A type that satisfies `ReadWriter` automatically satisfies `Reader` and `Writer` too.
- Ask for the smallest interface a function needs: take an `io.Reader`, not an `*os.File`.
- An embedded interface that was never set is nil, and calling one of its methods panics.
- The decorator pattern above is how logging, caching, retries and metrics are layered onto a dependency.

## 8. Standard library interfaces to know

A handful of small interfaces hold the standard library together, and most have only one method.

| Interface | Methods | Purpose |
| --- | --- | --- |
| `error` | `Error() string` | Every error value |
| `fmt.Stringer` | `String() string` | Custom text when printed |
| `io.Reader` | `Read(p []byte) (n int, err error)` | Any source of bytes |
| `io.Writer` | `Write(p []byte) (n int, err error)` | Any destination for bytes |
| `io.Closer` | `Close() error` | Releasing a resource |
| `sort.Interface` | `Len()`, `Less(i, j int) bool`, `Swap(i, j int)` | Custom sorting |
| `http.Handler` | `ServeHTTP(w http.ResponseWriter, r *http.Request)` | Handling an HTTP request |
| `json.Marshaler` | `MarshalJSON() ([]byte, error)` | Custom JSON encoding |
| `context.Context` | `Deadline`, `Done`, `Err`, `Value` | Cancellation and deadlines |

### One function, many destinations

```go
func greet(w io.Writer, name string) {
    fmt.Fprintf(w, "hello %s\n", name)
}

greet(os.Stdout, "Asha") // the terminal

var buf bytes.Buffer
greet(&buf, "Ravi")      // memory, which is handy in tests

// The same function also writes to files, network connections,
// HTTP responses and gzip streams, because all of them are io.Writers.
```

### Writing your own io.Writer

```go
type CountingWriter struct{ N int }

func (c *CountingWriter) Write(p []byte) (int, error) {
    c.N += len(p)
    return len(p), nil
}

cw := &CountingWriter{}
fmt.Fprint(cw, "hello")
fmt.Println(cw.N) // 5
```

### Stringer and error

```go
type Money struct{ Paise int }

func (m Money) String() string { return fmt.Sprintf("Rs %d.%02d", m.Paise/100, m.Paise%100) }

fmt.Println(Money{12550}) // Rs 125.50: fmt finds and calls String
```

- `io.Reader` signals the end of data by returning `io.EOF`.
- If a type has both `Error` and `String`, the `fmt` package uses `Error`.
- Since Go 1.21, `slices.SortFunc` usually replaces writing a `sort.Interface` by hand.

## 9. Designing with interfaces

Define a small interface in the package that uses it, accept that interface as a parameter, and return concrete types.

```go
// Package orders is the consumer, so it declares exactly what it needs.
type UserFinder interface {
    FindUser(ctx context.Context, id int) (*User, error)
}

type Service struct {
    users UserFinder
}

func NewService(users UserFinder) *Service { // accept an interface
    return &Service{users: users}            // return a concrete type
}

func (s *Service) Greeting(ctx context.Context, id int) (string, error) {
    u, err := s.users.FindUser(ctx, id)
    if err != nil {
        return "", fmt.Errorf("greeting: %w", err)
    }
    return "hello " + u.Name, nil
}

// Package postgres is the producer. It never mentions UserFinder.
type UserRepo struct{ db *sql.DB }

func (r *UserRepo) FindUser(ctx context.Context, id int) (*User, error) {
    // query the database
}
```

### Testing with a fake

```go
type fakeUsers struct {
    user *User
    err  error
}

func (f fakeUsers) FindUser(ctx context.Context, id int) (*User, error) {
    return f.user, f.err
}

func TestGreeting(t *testing.T) {
    svc := NewService(fakeUsers{user: &User{Name: "Asha"}})

    got, err := svc.Greeting(context.Background(), 1)
    if err != nil || got != "hello Asha" {
        t.Fatalf("got %q, %v", got, err)
    }
}

func TestGreetingError(t *testing.T) {
    svc := NewService(fakeUsers{err: errors.New("db down")})

    if _, err := svc.Greeting(context.Background(), 1); err == nil {
        t.Fatal("expected an error")
    }
}
```

| Guideline | Why |
| --- | --- |
| Keep interfaces small, one to three methods | Small interfaces are easy to satisfy and to fake |
| Define the interface where it is used | The consumer knows what it needs; the producer stays decoupled |
| Accept interfaces, return structs | Callers get flexibility in, and the full concrete type out |
| Do not create an interface until you need one | One implementation and no test double means no interface yet |
| Name one-method interfaces with an -er suffix | `Reader`, `Writer`, `Stringer`, `Notifier` |
| Avoid `any` in public APIs | It moves type errors from compile time to runtime |

- A Go proverb sums it up: the bigger the interface, the weaker the abstraction.
- This is dependency injection without a framework: pass dependencies into constructors as interfaces.
- Hand-written fakes like the one above are often enough. Tools such as `mockgen` generate mocks for larger interfaces.

## 10. Interfaces and generics

Since Go 1.18 an interface can also describe a set of types, which is how generic functions state their constraints.

```go
// A constraint: any type whose underlying type is one of these
type Number interface {
    ~int | ~int64 | ~float64
}

func Sum[T Number](xs []T) T {
    var total T
    for _, x := range xs {
        total += x // + is allowed because every type in the set supports it
    }
    return total
}

fmt.Println(Sum([]int{1, 2, 3}))      // 6
fmt.Println(Sum([]float64{1.5, 2.5})) // 4

// var n Number // compile error: a constraint interface cannot be a variable's type
```

### Interface or generics?

```go
// Interface: the behaviour differs for each type
func PrintAll(items []fmt.Stringer) {
    for _, it := range items {
        fmt.Println(it.String())
    }
}

// Generics: the same algorithm for every type, and the type is preserved
func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

m := Max(3, 7)     // m is an int, no assertion needed
s := Max("a", "b") // s is a string
```

|  | Interface | Generics |
| --- | --- | --- |
| Decided at | Runtime | Compile time |
| Best for | Different behaviour behind one method set | One algorithm over many types |
| Mixed types in one collection | Yes, as in `[]Shape` | No, one type per instantiation |
| Result type | The interface, so you may need an assertion | The exact type you passed in |
| Typical use | Plugins, dependencies, I/O | Containers, `Map`, `Filter`, `Min`, `Max` |

- `any` and `comparable` are the two built-in constraints. `cmp.Ordered` covers types that support the ordering operators.
- The `~` prefix means "any type whose underlying type is this", so `type Age int` matches `~int`.
- A constraint can combine a type set with methods, but such interfaces work only as constraints.
- If you find yourself writing `any` plus type assertions, generics are probably the better tool.

## 11. Internals and performance

An interface value occupies two machine words: one identifies the type and its methods, and the other points to the data.

| Kind | First word | Second word |
| --- | --- | --- |
| Empty interface (`any`) | Pointer to the type descriptor | Pointer to the data |
| Interface with methods | Pointer to an itab: the type plus a table of its methods for this interface | Pointer to the data |

- A method call through an interface looks up the function in the itab and calls it indirectly. This is dynamic dispatch.
- The runtime builds an itab once for each pair of interface and concrete type, then caches it.
- An indirect call costs a little more than a direct one and usually cannot be inlined.
- Storing a non-pointer value in an interface usually copies it to the heap, which is an allocation.
- In hot loops, a concrete type or a generic function avoids both costs. Everywhere else, the cost does not matter.

### Comparing interface values

```go
var a, b any = 1, 1
fmt.Println(a == b) // true: same dynamic type and equal values

var c any = int64(1)
fmt.Println(a == c) // false: int and int64 are different dynamic types

var x, y any = []int{1}, []int{1}
// fmt.Println(x == y) // panic: runtime error: comparing uncomparable type []int
```

- Two interface values are equal when their dynamic types are identical and their dynamic values are equal.
- The comparison compiles for any interface, and panics at runtime if the dynamic type is not comparable.
- For the same reason an interface can be a map key, and inserting a slice or map as that key panics.

## 12. Predict the output

Cover the answer comments and work out each of these four before reading on.

### Puzzle 1: does it compile?

```go
type Animal interface{ Sound() string }

type Dog struct{}

func (d *Dog) Sound() string { return "woof" }

func main() {
    var a Animal = Dog{}
    fmt.Println(a.Sound())
}
// Compile error: Dog does not implement Animal (method Sound has pointer receiver).
// Fix: var a Animal = &Dog{}
```

### Puzzle 2: copy or reference?

```go
type Valuer interface{ Value() int }

type Counter struct{ n int }

func (c Counter) Value() int { return c.n }

func main() {
    c := Counter{1}

    var v Valuer = c // stores a copy
    c.n = 100
    fmt.Println(v.Value()) // 1

    var p Valuer = &c // stores the address
    c.n = 200
    fmt.Println(p.Value()) // 200
}
```

### Puzzle 3: which case matches?

```go
type Both struct{}

func (Both) Error() string  { return "as error" }
func (Both) String() string { return "as stringer" }

func kind(v any) string {
    switch v.(type) {
    case fmt.Stringer:
        return "stringer"
    case error:
        return "error"
    }
    return "other"
}

func main() {
    fmt.Println(kind(Both{})) // stringer: the first matching case wins
    fmt.Println(Both{})       // as error: fmt checks for error before Stringer
}
```

### Puzzle 4: nil or not?

```go
type Logger interface{ Log(string) }

type FileLogger struct{}

func (f *FileLogger) Log(s string) { fmt.Println("log:", s) }

func newLogger(enabled bool) *FileLogger {
    if enabled {
        return &FileLogger{}
    }
    return nil
}

func main() {
    var l Logger = newLogger(false)
    if l != nil {
        l.Log("hello") // this line runs
    }
}
// log: hello
```

The interface holds a nil `*FileLogger`, so it is not nil. The call succeeds only because `Log` never touches a field of its nil receiver. The fix is for `newLogger` to return the `Logger` interface type.

## 13. Interview coding problems

These four exercises show the main jobs interfaces do: swapping behaviour, classifying errors, plugging into the standard library, and wrapping a dependency.

### Strategy: interchangeable payment methods

```go
type PaymentMethod interface {
    Pay(amount float64) string
}

type Card struct{ Last4 string }

func (c Card) Pay(amount float64) string {
    return fmt.Sprintf("paid %.2f with card ending %s", amount, c.Last4)
}

type UPI struct{ ID string }

func (u UPI) Pay(amount float64) string {
    return fmt.Sprintf("paid %.2f via UPI %s", amount, u.ID)
}

func checkout(m PaymentMethod, amount float64) {
    fmt.Println(m.Pay(amount))
}

checkout(Card{"4242"}, 499)    // paid 499.00 with card ending 4242
checkout(UPI{"asha@upi"}, 250) // paid 250.00 via UPI asha@upi
```

Adding a wallet later means adding one type. `checkout` does not change.

### A custom error type, found with errors.As

```go
type NotFoundError struct {
    Resource string
    ID       int
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s %d not found", e.Resource, e.ID)
}

func getUser(id int) error {
    return fmt.Errorf("getUser: %w", &NotFoundError{"user", id}) // wrapped
}

err := getUser(7)
fmt.Println(err) // getUser: user 7 not found

var nf *NotFoundError
if errors.As(err, &nf) { // walks the chain looking for this type
    fmt.Println(nf.Resource, nf.ID) // user 7
}
```

### Implementing sort.Interface

```go
type Person struct {
    Name string
    Age  int
}

type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

people := []Person{{"Ravi", 35}, {"Asha", 28}, {"Meera", 31}}
sort.Sort(ByAge(people))
fmt.Println(people) // [{Asha 28} {Meera 31} {Ravi 35}]
```

### Decorator: add retries to any Notifier

```go
type Notifier interface {
    Notify(msg string) error
}

type RetryNotifier struct {
    next     Notifier // the wrapped implementation
    attempts int
}

func (r RetryNotifier) Notify(msg string) error {
    var err error
    for i := 0; i < r.attempts; i++ {
        if err = r.next.Notify(msg); err == nil {
            return nil
        }
    }
    return fmt.Errorf("after %d attempts: %w", r.attempts, err)
}

var n Notifier = RetryNotifier{next: Email{To: "ops@x.com"}, attempts: 3}
_ = n.Notify("disk full") // callers cannot tell it is wrapped
```

- Be ready to explain why each receiver is a value or a pointer.
- For any design question, start by naming the interface and its one or two methods, then write the types.
- A common follow-up: "how would you test this?" The answer is a fake that satisfies the same interface.

## 14. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| What is an interface? | A set of method signatures that a type satisfies by having those methods. |
| How does a type implement an interface? | Implicitly; there is no `implements` keyword. |
| What does an interface value contain? | A dynamic type and a dynamic value, stored in two words. |
| When is an interface nil? | Only when both its type and its value are unset. |
| Why can an error be non-nil while holding a nil pointer? | The interface still carries the pointer's type. |
| What is `any`? | An alias for `interface{}`, satisfied by every type. |
| Type assertion: one result or two? | One result panics on a mismatch; the comma-ok form does not. |
| What is a type switch? | A switch on the dynamic type of an interface value. |
| Can a value satisfy an interface with pointer-receiver methods? | No; only the pointer type can. |
| How do you check satisfaction at compile time? | `var _ Iface = (*T)(nil)`. |
| How do you compose interfaces? | Embed small interfaces in a larger one. |
| Where should an interface be defined? | In the package that uses it. |
| What does "accept interfaces, return structs" mean? | Take the smallest interface you need and return a concrete type. |
| How large should an interface be? | As small as possible, often one method. |
| Does Go have inheritance? | No; polymorphism comes from interfaces and reuse from embedding. |
| Can you compare two interface values? | Yes, but it panics if the dynamic type is not comparable. |
| Is `[]string` assignable to `[]any`? | No; convert it element by element. |
| Interface or generics? | An interface for varying behaviour; generics for one algorithm over many types. |
| What is an itab? | The runtime table linking a concrete type to an interface's methods. |
| Does an interface call cost more than a direct call? | Slightly: it is an indirect call and usually not inlined. |
