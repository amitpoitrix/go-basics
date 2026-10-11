# Generics in Go

Oct 11, 2026 · @HR

## 1. Why generics exist

Generics let you write one function or type that works for many types, while the compiler still checks every use.

```go
// Before Go 1.18, option 1: copy the code for each type
func MaxInt(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func MaxFloat(a, b float64) float64 {
    if a > b {
        return a
    }
    return b
}

// Option 2: use any, and lose type safety
func Contains(items []any, target any) bool {
    for _, it := range items {
        if it == target {
            return true
        }
    }
    return false
}
// Callers had to convert a []string into a []any first, and mistakes showed up only at runtime.

// With generics: one definition, fully type-checked
func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

fmt.Println(Max(3, 7))     // 7
fmt.Println(Max(2.5, 1.5)) // 2.5
fmt.Println(Max("a", "b")) // b
// Max(1, "a")             // compile error: mismatched types
```

- Generics arrived in Go 1.18. The formal name is type parameters.
- There are three pieces: type parameters on functions, type parameters on types, and constraints that limit which types are allowed.
- The benefit over `any` is that the concrete type is known at compile time: no assertions, no boxing, and errors are caught early.
- The benefit over copying code is one implementation to test and maintain.

## 2. Generic functions and type inference

A generic function declares type parameters in square brackets before its ordinary parameters, and the compiler usually works out the type arguments from the call.

```go
func Map[T, U any](xs []T, f func(T) U) []U {
    out := make([]U, 0, len(xs))
    for _, x := range xs {
        out = append(out, f(x))
    }
    return out
}

// Inferred: T is string and U is int
lengths := Map([]string{"go", "rust"}, func(s string) int { return len(s) }) // [2 4]

// Explicit type arguments: allowed, rarely needed
strs := Map[int, string]([]int{1, 2}, strconv.Itoa) // [1 2] as strings
```

### When inference cannot help

```go
// T appears only in the result, so there is nothing to infer it from
func Zero[T any]() T {
    var z T
    return z
}

n := Zero[int]()    // 0
s := Zero[string]() // ""

// Put the type the caller must supply first, so the rest can still be inferred
func Convert[To, From ~int | ~float64](x From) To {
    return To(x)
}

f := Convert[float64](42) // To is given; From is inferred as int
```

### Instantiating without calling

```go
intMax := Max[int] // a plain func(int, int) int
fmt.Println(intMax(2, 9)) // 9
```

- The syntax is `func Name[TypeParam Constraint](params) results`.
- Type parameters are conventionally one capital letter: `T` for a type, `K` and `V` for map keys and values, `E` for elements.
- Inference uses the argument types. It cannot look at how the result will be used.
- Supplying type arguments is called instantiation, and it happens at compile time.
- Parameters sharing a constraint can be grouped, as in `[T, U any]`.

## 3. Constraints

A constraint is an interface that names the set of types a type parameter may be, and it decides which operations the function body can use.

| Constraint | Allows | Body may use |
| --- | --- | --- |
| `any` | Every type | Assignment and passing only; no operators |
| `comparable` | Types that support equality | `==` and `!=`; required for map keys |
| `cmp.Ordered` | Integers, floats and strings | The ordering operators as well as equality |
| A union such as `~int \| ~float64` | Exactly the listed types | Operators that every listed type supports |
| An interface with methods | Types that have those methods | Those methods |

```go
// comparable: equality is available
func Index[T comparable](xs []T, target T) int {
    for i, x := range xs {
        if x == target {
            return i
        }
    }
    return -1
}

// A custom type set: arithmetic is available
type Number interface {
    ~int | ~int32 | ~int64 | ~float32 | ~float64
}

func Sum[T Number](xs []T) T {
    var total T
    for _, x := range xs {
        total += x
    }
    return total
}

type Celsius float64

fmt.Println(Sum([]int{1, 2, 3}))          // 6
fmt.Println(Sum([]Celsius{20.5, 21.5}))   // 42: accepted because of the ~

// A method constraint: the method is available
func JoinAll[T fmt.Stringer](xs []T, sep string) string {
    parts := make([]string, len(xs))
    for i, x := range xs {
        parts[i] = x.String()
    }
    return strings.Join(parts, sep)
}

// An inline constraint, for one-off use
func Double[T ~int | ~float64](x T) T { return x * 2 }

// Type set and method together
type PrintableNumber interface {
    ~int | ~float64
    String() string
}
```

- The `~` means "any type whose underlying type is this". Without it, a named type such as `Celsius` is rejected.
- The body may use only the operations that every type in the set supports. That is why `any` permits so little.
- An interface containing a union can be used only as a constraint, never as the type of a variable.
- `golang.org/x/exp/constraints` offers ready-made sets such as `Integer`, `Float` and `Signed`.
- Choose the loosest constraint that still compiles the body. It keeps the function usable with more types.

## 4. Generic types

A struct, slice, map or other named type can take type parameters too, which is how you build type-safe containers.

```go
type Stack[T any] struct {
    items []T
}

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

func (s *Stack[T]) Len() int { return len(s.items) }

var ints Stack[int] // the type argument is written out
ints.Push(1)
// ints.Push("a")   // compile error: this stack holds ints

names := &Stack[string]{}
names.Push("asha")
```

### More shapes

```go
// Two type parameters
type Pair[K comparable, V any] struct {
    Key   K
    Value V
}

p := Pair[string, int]{Key: "age", Value: 30}

// A constructor function restores inference
func NewPair[K comparable, V any](k K, v V) Pair[K, V] {
    return Pair[K, V]{Key: k, Value: v}
}

q := NewPair("age", 30) // Pair[string, int], no brackets needed

// A generic type built on a map
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T)      { s[v] = struct{}{} }
func (s Set[T]) Has(v T) bool { _, ok := s[v]; return ok }

seen := Set[string]{}
seen.Add("go")

// A recursive generic type
type Tree[T cmp.Ordered] struct {
    Left, Right *Tree[T]
    Val         T
}
```

### Methods with their own type parameters (Go 1.27)

```go
// Go 1.27 and later: a method may declare type parameters of its own
func (s *Stack[T]) MapTo[U any](f func(T) U) []U {
    out := make([]U, 0, len(s.items))
    for _, v := range s.items {
        out = append(out, f(v))
    }
    return out
}

labels := ints.MapTo(strconv.Itoa) // U inferred as string

// Before Go 1.27 the only option was a package-level function
func MapStack[T, U any](s *Stack[T], f func(T) U) []U {
    out := make([]U, 0, len(s.items))
    for _, v := range s.items {
        out = append(out, f(v))
    }
    return out
}
```

Generic methods arrived in August 2026 with Go 1.27. Interface methods still cannot declare type parameters, and a generic method cannot be used to satisfy an interface method ([Go 1.27 release notes](https://go.dev/doc/go1.27)).

- A method on a generic type repeats the type parameters in its receiver: `(s *Stack[T])`.
- `Stack[int]` and `Stack[string]` are two distinct types.
- Type arguments for a generic type are written explicitly. Wrap construction in a function when you want inference.
- Many codebases are still on older versions, so know the package-level function form as well.

## 5. Generics in the standard library

Most generic code you need already exists in the standard library, so reach for it before writing your own helper.

| Package or built-in | Since | What it offers |
| --- | --- | --- |
| `min`, `max`, `clear` | Go 1.21 | Built-in functions for ordered values, and for emptying maps and slices |
| `slices` | Go 1.21 | `Sort`, `SortFunc`, `Contains`, `Index`, `Reverse`, `Max`, `Min`, `BinarySearch`, `Compact`, `Clone`, `Equal`, `Insert`, `Delete` |
| `maps` | Go 1.21 | `Keys`, `Values`, `Clone`, `Equal`, `DeleteFunc` |
| `cmp` | Go 1.21 | The `Ordered` constraint and the `Compare` function |
| `sync` | Go 1.21 | `OnceValue` and `OnceValues` for lazy, one-time initialization |
| `sync/atomic` | Go 1.19 | `atomic.Pointer[T]` for lock-free pointer swaps |
| `iter` | Go 1.23 | `Seq` and `Seq2`, the iterator types used by `range` over functions |

```go
nums := []int{3, 1, 2}
slices.Sort(nums)                        // [1 2 3]
idx, found := slices.BinarySearch(nums, 2) // 1 true
fmt.Println(slices.Max(nums), min(4, 2), max(1.5, 2.5)) // 3 2 2.5

m := map[string]int{"b": 2, "a": 1}
keys := slices.Sorted(maps.Keys(m)) // [a b]; Go 1.23

type Person struct {
    Name string
    Age  int
}

people := []Person{{"Ravi", 35}, {"Asha", 28}}
slices.SortFunc(people, func(a, b Person) int {
    return cmp.Compare(a.Age, b.Age)
}) // [{Asha 28} {Ravi 35}]

// Type-safe atomic pointer: no casts from unsafe.Pointer
var current atomic.Pointer[Config]
current.Store(&Config{Port: 8080})
cfg := current.Load() // *Config

// Compute once, on first use, safely across goroutines
loadConfig := sync.OnceValue(func() *Config {
    return readConfigFile()
})
cfg = loadConfig() // later calls return the same value
```

- Before Go 1.21 the equivalents were `sort.Ints`, `sort.Slice` and hand-written loops.
- Since Go 1.23, `maps.Keys` and `maps.Values` return iterators. Wrap them in `slices.Collect` or `slices.Sorted` to get a slice.
- These functions are a good model for your own signatures. Read `slices.Index` and `slices.SortFunc` once.

## 6. Useful generic helpers

These six short functions are the ones interviewers most often ask you to write from scratch.

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

func Reduce[T, A any](xs []T, initial A, f func(A, T) A) A {
    acc := initial
    for _, x := range xs {
        acc = f(acc, x)
    }
    return acc
}

func Keys[K comparable, V any](m map[K]V) []K {
    keys := make([]K, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return keys
}

func Unique[T comparable](xs []T) []T {
    seen := make(map[T]struct{}, len(xs))
    var out []T
    for _, x := range xs {
        if _, ok := seen[x]; !ok {
            seen[x] = struct{}{}
            out = append(out, x)
        }
    }
    return out
}

func GroupBy[T any, K comparable](xs []T, key func(T) K) map[K][]T {
    groups := make(map[K][]T)
    for _, x := range xs {
        k := key(x)
        groups[k] = append(groups[k], x)
    }
    return groups
}

// Ptr returns a pointer to any value: handy for optional struct fields
func Ptr[T any](v T) *T { return &v }
```

```go
evens := Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 }) // [2 4]
sum := Reduce(evens, 0, func(acc, n int) int { return acc + n })         // 6
uniq := Unique([]string{"a", "b", "a"})                                  // [a b]

byLen := GroupBy([]string{"go", "rust", "c", "java"}, func(s string) int { return len(s) })
// map[1:[c] 2:[go] 4:[rust java]]

type Patch struct {
    Name *string
    Age  *int
}

p := Patch{Name: Ptr("Asha"), Age: Ptr(30)} // no temporary variables needed
```

- Notice the constraints: `comparable` appears only where a value is used as a map key.
- `Map` from section 2 completes the usual trio of `Map`, `Filter` and `Reduce`.
- In production, check `slices` and `maps` first; several of these already exist there under similar names.

## 7. Limitations and gotchas

Go's generics are deliberately small, and these restrictions are the ones that surprise people and show up in interviews.

### Operators need a constraint

```go
func Add[T any](a, b T) T {
    return a + b // compile error: operator + not defined on a, constrained by any
}

func Add[T ~int | ~float64 | ~string](a, b T) T {
    return a + b // fine: every type in the set supports +
}
```

### There is no nil or literal zero for T

```go
func First[T any](xs []T) (T, bool) {
    if len(xs) == 0 {
        // return nil, false // compile error: T might be an int
        var zero T
        return zero, false
    }
    return xs[0], true
}
```

### You cannot switch on T directly

```go
func Describe[T any](v T) string {
    switch any(v).(type) { // convert to any first
    case int:
        return "int"
    case string:
        return "string"
    }
    return "other"
}
```

Needing this is usually a sign that an interface would fit better than a type parameter.

### A union does not give access to shared fields

```go
type A struct{ ID int }
type B struct{ ID int }

func GetID[T A | B](v T) int {
    return v.ID // compile error, even though both types have ID
}

// Use a method in the constraint instead
type HasID interface{ GetID() int }

func IDOf[T HasID](v T) int { return v.GetID() }
```

| Limitation | Detail |
| --- | --- |
| Generic methods | Not allowed before Go 1.27. Allowed since, but they cannot satisfy interface methods, and interface methods cannot have type parameters ([release notes](https://go.dev/doc/go1.27)) |
| No specialization | You cannot supply a different implementation for one particular type argument |
| No variadic type parameters | A function has a fixed number of type parameters |
| Constraint interfaces | An interface containing a union cannot be used as a variable's type |
| Type arguments on types | Must be written out; only function calls get inference |
| `comparable` with interfaces | Interface types satisfy `comparable` since Go 1.20, but `==` can panic at runtime if the dynamic type is not comparable |

- Do not reach for generics by default. Start with concrete types and add a type parameter when you see real duplication.
- Heavy use of type parameters and long constraint lists makes signatures hard to read, which Go code reviews push back on.

## 8. Generics vs interfaces vs any

Use an interface when behaviour differs per type, generics when one algorithm serves many types, and `any` only when the type truly cannot be known.

|  | Interface | Generics | `any` |
| --- | --- | --- | --- |
| Type known at | Runtime | Compile time | Runtime |
| Type safety | Method set is checked | Fully checked | None until you assert |
| Mixed types in one collection | Yes | No | Yes |
| Result keeps the concrete type | No | Yes | No |
| Typical use | Dependencies, plugins, I/O | Containers and algorithms | JSON of unknown shape, `fmt.Println` |

### A type parameter that adds nothing

```go
// Unnecessary: the body only calls a method
func WriteAll[W io.Writer](w W, data []byte) error {
    _, err := w.Write(data)
    return err
}

// Simpler and equivalent
func WriteAll(w io.Writer, data []byte) error {
    _, err := w.Write(data)
    return err
}
```

### Where generics clearly win

```go
// 1. The result must keep the caller's type
func Max[T cmp.Ordered](a, b T) T {
    if a > b {
        return a
    }
    return b
}

// 2. Containers: one implementation, any element type, no assertions
type Queue[T any] struct{ items []T }

// 3. Preserving a named slice type with the ~[]E pattern
type IDs []int

func Reversed[S ~[]E, E any](s S) S {
    out := make(S, len(s))
    for i, v := range s {
        out[len(s)-1-i] = v
    }
    return out
}

ids := Reversed(IDs{1, 2, 3}) // the result is IDs, not []int
```

- Rule of thumb: if the function only calls methods on its argument, take an interface.
- If the same type appears in more than one place, as two parameters or a parameter and a result, generics keep them tied together.
- A heterogeneous list such as `[]Shape` needs an interface. Generics give you one type per instantiation.
- The `~[]E` pattern is how the `slices` package returns your own slice type instead of a plain `[]E`.

## 9. How generics are implemented

Go compiles generics with a hybrid approach called GC shape stenciling: one copy of the code per memory layout, plus a hidden dictionary of type information.

| Approach | Used by | Trade-off |
| --- | --- | --- |
| Full monomorphization: a copy per type | C++, Rust | Fastest code, larger binaries, slower builds |
| Type erasure: one copy, values boxed | Java | Small binaries, boxing and casts at runtime |
| GC shape stenciling with dictionaries | Go | In between: a copy per shape, shared across similar types |

- Types with the same underlying layout share a "GC shape". Each distinct value type, such as `int`, `float64` or a particular struct, gets its own copy.
- All pointer types share a single shape, so `Stack[*User]` and `Stack[*Order]` run the same machine code.
- The compiler passes each instantiation a dictionary describing the actual type, used where the shared code needs type-specific details.
- Everything is resolved at compile time. There is no reflection and, for value types, no boxing.

### What it means for performance

- Generic code over value types is about as fast as code written by hand for that type.
- Calling a method on a type parameter whose argument is a pointer or interface goes through the dictionary. It may not be inlined and can be slower than you expect.
- Compared with `any`, generics avoid the heap allocation of boxing and the cost of type assertions.
- Binaries grow a little with each new shape, and compile times rise slightly.
- As always, measure with a benchmark before optimizing: `go test -bench=. -benchmem`.

## 10. Interview coding problems

These four exercises cover the usual requests: a concurrent container, a data structure with an iterator, a classic algorithm, and generics combined with goroutines.

### Thread-safe generic cache

```go
type Cache[K comparable, V any] struct {
    mu sync.RWMutex
    m  map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
    return &Cache[K, V]{m: make(map[K]V)}
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.m[k]
    return v, ok
}

func (c *Cache[K, V]) Set(k K, v V) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.m[k] = v
}

users := NewCache[int, string]()
users.Set(1, "Asha")
name, ok := users.Get(1) // Asha true, with no type assertion
```

### Linked list with an iterator (Go 1.23)

```go
type node[T any] struct {
    val  T
    next *node[T]
}

type List[T any] struct {
    head, tail *node[T]
}

func (l *List[T]) Push(v T) {
    n := &node[T]{val: v}
    if l.tail == nil {
        l.head, l.tail = n, n
        return
    }
    l.tail.next = n
    l.tail = n
}

// All returns an iterator, so callers can use a plain for-range loop
func (l *List[T]) All() iter.Seq[T] {
    return func(yield func(T) bool) {
        for n := l.head; n != nil; n = n.next {
            if !yield(n.val) {
                return // the caller broke out of the loop
            }
        }
    }
}

var list List[string]
list.Push("a")
list.Push("b")
for v := range list.All() {
    fmt.Println(v) // a, then b
}
```

### Binary search over any ordered type

```go
func BinarySearch[T cmp.Ordered](xs []T, target T) (int, bool) {
    lo, hi := 0, len(xs)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2 // avoids overflow
        switch {
        case xs[mid] == target:
            return mid, true
        case xs[mid] < target:
            lo = mid + 1
        default:
            hi = mid - 1
        }
    }
    return lo, false // lo is where target would be inserted
}

fmt.Println(BinarySearch([]int{1, 3, 5, 7}, 5))       // 2 true
fmt.Println(BinarySearch([]string{"a", "c"}, "b")) // 1 false
```

### Parallel map

```go
func ParallelMap[T, U any](xs []T, f func(T) U) []U {
    out := make([]U, len(xs))
    var wg sync.WaitGroup
    for i, x := range xs {
        wg.Add(1)
        go func() {
            defer wg.Done()
            out[i] = f(x) // each goroutine writes only its own index, so no lock is needed
        }()
    }
    wg.Wait()
    return out
}

squares := ParallelMap([]int{1, 2, 3}, func(n int) int { return n * n }) // [1 4 9]
```

- For the cache, expect follow-ups on a TTL per entry and on LRU eviction.
- `ParallelMap` relies on per-iteration loop variables, so it needs Go 1.22 or later as written.
- A natural follow-up for `ParallelMap` is limiting concurrency with a semaphore or `errgroup.SetLimit`.

## 11. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| When were generics added? | Go 1.18. |
| What is a type parameter? | A placeholder for a type, declared in square brackets. |
| What is a constraint? | An interface listing the types a type parameter may be. |
| What does `any` mean as a constraint? | Every type is allowed, but no operators can be used. |
| What does `comparable` allow? | `==` and `!=`; it is required for map keys. |
| What is `cmp.Ordered`? | Integers, floats and strings: the types that support ordering. |
| What does `~int` mean? | Any type whose underlying type is `int`. |
| What is type inference? | The compiler deducing type arguments from the call's arguments. |
| When must type arguments be explicit? | When a type parameter appears only in the result, and for generic types. |
| How do you get the zero value of T? | `var zero T`. |
| Can you type-switch on a type parameter? | Only after converting the value to `any`. |
| Can methods have type parameters? | Yes, since Go 1.27; not before, and never on interface methods. |
| Can a constraint interface be a variable's type? | Not if it contains a union of types. |
| Generics or interface? | Generics for one algorithm over many types; an interface for varying behaviour. |
| Generics or `any`? | Generics: the type is checked at compile time and nothing is boxed. |
| How are generics implemented? | GC shape stenciling with dictionaries. |
| Do generics slow programs down? | Not for value types; calls through pointer shapes can cost a little. |
| Which standard packages are generic? | `slices`, `maps`, `cmp`, `iter`, plus `min`, `max` and `clear`. |
| What is the `~[]E` pattern for? | Returning the caller's named slice type instead of a plain slice. |
| Does Go support specialization? | No; one implementation serves every type argument. |
