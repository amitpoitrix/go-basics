# Arrays, Slices and Maps in Go

Oct 11, 2026 · @HR

## 1. At a glance

An array is a fixed-size value, a slice is a resizable view over an array, and a map is a hash table of key-value pairs.

|  | Array | Slice | Map |
| --- | --- | --- | --- |
| Size | Fixed, part of the type | Dynamic | Dynamic |
| Zero value | Every element zeroed | nil | nil |
| Assignment copies | All the elements | The header only; the array is shared | A pointer only; the table is shared |
| Comparable with `==` | Yes, if the elements are | Only against nil | Only against nil |
| Usable as a map key | Yes | No | No |
| Created with | A literal | A literal, `make`, or slicing | A literal or `make` |
| Ordered | Yes | Yes | No |

## 2. Arrays

An array has a fixed length that is part of its type, and it is a value: assigning or passing it copies every element.

```go
var a [3]int              // [0 0 0]
b := [3]int{1, 2, 3}
c := [...]int{1, 2, 3, 4} // the compiler counts: [4]int
d := [5]int{1: 10, 3: 30} // [0 10 0 30 0]

e := b // full copy
e[0] = 99
fmt.Println(b[0], e[0])           // 1 99
fmt.Println(b == [3]int{1, 2, 3}) // true: arrays are comparable
// var x [4]int = b               // compile error: [3]int and [4]int are different types
```

### Arrays and functions

```go
func double(arr [3]int) { // receives a copy
    for i := range arr {
        arr[i] *= 2
    }
}

func doublePtr(arr *[3]int) { // receives a pointer to the original
    for i := range arr {
        arr[i] *= 2
    }
}

nums := [3]int{1, 2, 3}
double(nums)
fmt.Println(nums) // [1 2 3]: unchanged
doublePtr(&nums)
fmt.Println(nums) // [2 4 6]
```

- `len(arr)` is a compile-time constant.
- Arrays suit data whose size never changes: a hash digest (`[32]byte`), a grid, a lookup table.
- An array can be a map key or a struct field that stays comparable; a slice cannot.
- In everyday code you will use slices far more often than arrays.

## 3. Slice basics

A slice is a three-word header (pointer, length, capacity) that describes a section of a backing array.

```go
s1 := []int{1, 2, 3}     // literal: len 3, cap 3
s2 := make([]int, 3)     // [0 0 0]: len 3, cap 3
s3 := make([]int, 0, 10) // []: len 0, cap 10
var s4 []int             // nil: len 0, cap 0

arr := [5]int{10, 20, 30, 40, 50}
s5 := arr[1:4]                // [20 30 40]
fmt.Println(len(s5), cap(s5)) // 3 4: cap runs from index 1 to the end of arr
```

| Field | Meaning |
| --- | --- |
| Pointer | Address of the slice's first element in the backing array |
| Length | Number of elements you can index: `len(s)` |
| Capacity | Elements from the first one to the end of the backing array: `cap(s)` |

### nil slice vs empty slice

```go
var a []int  // nil: a == nil is true
b := []int{} // empty but not nil

fmt.Println(len(a), len(b)) // 0 0
a = append(a, 1)            // works on a nil slice
// json.Marshal gives null for a nil slice and [] for an empty one
```

### The make trap

```go
s := make([]int, 3) // length 3: already holds three zeros
s = append(s, 1)
fmt.Println(s) // [0 0 0 1], not [1]

t := make([]int, 0, 3) // length 0, capacity 3: what you usually want before append
t = append(t, 1)
fmt.Println(t) // [1]
```

- Passing a slice to a function copies the 24-byte header, never the elements.
- Check for emptiness with `len(s) == 0`, which covers both nil and empty slices.
- Indexing at or beyond `len(s)` panics at runtime, even when capacity is larger.

## 4. Slicing and shared backing arrays

Slicing never copies elements: the new slice points into the same backing array, so a write through one slice is visible through the other.

| Expression | Elements | Capacity |
| --- | --- | --- |
| `s[low:high]` | From low up to, not including, high | `cap(s) - low` |
| `s[:high]` | From 0 up to high | `cap(s)` |
| `s[low:]` | From low to the end | `cap(s) - low` |
| `s[:]` | Everything | `cap(s)` |
| `s[low:high:max]` | From low up to high | `max - low` |

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3] // [2 3], len 2, cap 4

b[0] = 20      // writes into the shared array
fmt.Println(a) // [1 20 3 4 5]

b = append(b, 100) // fits within cap, so it overwrites a[3]
fmt.Println(a)     // [1 20 3 100 5]

c := a[1:3:3]      // full slice expression: cap limited to 2
c = append(c, 200) // exceeds cap, so Go allocates a new array
c[0] = 999
fmt.Println(a)     // [1 20 3 100 5]: a is no longer affected
```

### Reslicing can grow a slice

```go
s := []int{1, 2, 3, 4, 5}
t := s[:2]     // [1 2], cap 5
t = t[:4]      // [1 2 3 4]: you may extend up to capacity
// t = t[:6]   // panic: beyond capacity
```

- The interview question "what does this print?" almost always hinges on whether an `append` stayed within capacity.
- Use the three-index form `s[low:high:max]` when you hand out a sub-slice and do not want appends to overwrite your data.
- For a fully independent copy, use `slices.Clone` or `copy`.

## 5. append and capacity growth

`append` writes into spare capacity when there is room and otherwise allocates a bigger array and copies, which is why you must always use its return value.

```go
var s []int
for i := 0; i < 10; i++ {
    s = append(s, i)
    fmt.Print(cap(s), " ") // 1 2 4 4 8 8 8 8 16 16
}

s = append(s, 10, 11)   // several values at once
s = append(s, other...) // all elements of another slice

b := append([]byte("hi "), "there"...) // a string can be appended to bytes
```

### Append inside a function

```go
func addBroken(s []int) {
    s = append(s, 4) // changes the local header only
}

func addFixed(s []int) []int {
    return append(s, 4) // return the new header
}

func main() {
    s := []int{1, 2, 3}
    addBroken(s)
    fmt.Println(s) // [1 2 3]: the caller's length is still 3

    s = addFixed(s)
    fmt.Println(s) // [1 2 3 4]
}
```

### Preallocate when you know the size

```go
ids := make([]int, 0, len(users)) // one allocation instead of many
for _, u := range users {
    ids = append(ids, u.ID)
}
```

- Capacity roughly doubles while the slice is small, then grows by about 1.25 times. Exact numbers are an implementation detail.
- After a reallocation the slice points to a new array, so older slices of the same data stop seeing changes.
- A function can modify existing elements through its copy of the header, but it cannot change the caller's length.
- Appending is amortized O(1); a single append that triggers growth is O(n).

## 6. copy, common operations and the slices package

`copy` and `append` are enough for every slice operation, and the `slices` package (Go 1.21) wraps the common ones.

```go
src := []int{1, 2, 3}
dst := make([]int, 2)
n := copy(dst, src) // n is 2: copies min(len(dst), len(src)) elements

clone := slices.Clone(src) // independent copy

// Remove index i, keeping order: O(n)
s = append(s[:i], s[i+1:]...)

// Remove index i when order does not matter: O(1)
s[i] = s[len(s)-1]
s = s[:len(s)-1]

// Insert v at index i
s = slices.Insert(s, i, v)

// Stack
stack = append(stack, v)     // push
top := stack[len(stack)-1]   // peek
stack = stack[:len(stack)-1] // pop

// Queue
queue = append(queue, v) // enqueue
front := queue[0]        // peek
queue = queue[1:]        // dequeue

// Two-dimensional slice: each row is its own slice
grid := make([][]int, 3)
for i := range grid {
    grid[i] = make([]int, 4)
}
```

| Call | What it does |
| --- | --- |
| `slices.Sort(s)` | Sorts in place, ascending |
| `slices.SortFunc(s, cmp)` | Sorts in place with your comparison function |
| `slices.Contains(s, v)` | Reports whether v is present |
| `slices.Index(s, v)` | Index of v, or -1 |
| `slices.Equal(a, b)` | Compares length and elements |
| `slices.Reverse(s)` | Reverses in place |
| `slices.Max(s)`, `slices.Min(s)` | Largest and smallest; panic on an empty slice |
| `slices.BinarySearch(s, v)` | Position and found flag, on a sorted slice |
| `slices.Compact(s)` | Removes consecutive duplicates |
| `slices.Clone(s)` | Shallow copy |

```go
type Person struct {
    Name string
    Age  int
}

slices.SortFunc(people, func(a, b Person) int {
    return cmp.Compare(a.Age, b.Age) // negative, zero or positive
})
```

- `copy` never grows the destination. Copying into a nil or zero-length slice copies nothing.
- Before Go 1.21, the equivalents were `sort.Ints`, `sort.Strings` and `sort.Slice`.
- A clone is shallow: a slice of pointers still shares the pointed-to values.

## 7. Slice gotchas

These five behaviours cause most slice bugs and most "predict the output" questions.

### range gives you a copy of each element

```go
type User struct {
    Name string
    Age  int
}

users := []User{{"Asha", 30}, {"Ravi", 25}}

for _, u := range users {
    u.Age++ // modifies a copy; users is unchanged
}

for i := range users {
    users[i].Age++ // modifies the real element
}
```

### range reads the length once

```go
s := []int{1, 2, 3}
for range s {
    s = append(s, 0) // the loop still runs exactly 3 times
}
fmt.Println(len(s)) // 6
```

### Taking the address of the loop variable

```go
var ptrs []*int
for _, v := range []int{1, 2, 3} {
    ptrs = append(ptrs, &v)
}
for _, p := range ptrs {
    fmt.Print(*p, " ")
}
// Go 1.22+: 1 2 3 (a new v for each iteration)
// Before 1.22: 3 3 3 (one v shared by every iteration)
```

### Slices cannot be compared with ==

```go
a, b := []int{1, 2}, []int{1, 2}
// fmt.Println(a == b)           // compile error: slices compare only to nil
fmt.Println(slices.Equal(a, b)) // true
```

### A small slice keeps a big array alive

```go
func firstTen(data []byte) []byte {
    return data[:10] // still references the whole backing array
}

func firstTenCopy(data []byte) []byte {
    return slices.Clone(data[:10]) // only 10 bytes stay alive
}
```

- Slices are not safe for concurrent writes. Guard a shared slice with a mutex, or give each goroutine its own index range.
- A slice of pointers that you shrink should have the removed slots set to nil, so the garbage collector can free those objects.

## 8. Map basics

A map is an unordered collection of key-value pairs with average O(1) lookup, insert and delete.

```go
m := map[string]int{"a": 1, "b": 2} // literal
m2 := make(map[string]int)          // empty and ready to use
m3 := make(map[string]int, 100)     // with a size hint, avoids early regrowth

m["c"] = 3           // insert or update
v := m["zzz"]        // 0: a missing key returns the zero value
v, ok := m["zzz"]    // 0 false: comma-ok tells missing from a stored zero
delete(m, "a")       // no-op if the key is absent
fmt.Println(len(m))  // number of entries

if age, ok := ages["asha"]; ok {
    fmt.Println("found", age)
}

for k, v := range m { // order is random
    fmt.Println(k, v)
}
```

### The nil map

```go
var nm map[string]int // nil

_ = nm["x"]       // fine: returns 0
_ = len(nm)       // fine: 0
delete(nm, "x")   // fine: no-op
for range nm {}   // fine: zero iterations
nm["x"] = 1       // panic: assignment to entry in nil map
```

### Idioms built on the zero value

```go
// Counter: a missing key starts at 0
counts := map[string]int{}
counts["go"]++

// Grouping: a missing key gives a nil slice, which append accepts
groups := map[string][]string{}
groups["fruit"] = append(groups["fruit"], "apple")

// Set: struct{} occupies zero bytes
set := map[string]struct{}{}
set["go"] = struct{}{}
_, exists := set["go"] // true
```

- A nil map is readable but not writable. Always create a map with a literal or `make` before writing.
- `clear(m)` (Go 1.21) removes every entry.
- `map[string]bool` also works as a set and reads more simply; `struct{}` saves memory.

## 9. Map behaviour and gotchas

Maps have random iteration order, behave like references when passed around, and hold values you cannot modify in place.

### Iterating in sorted order

```go
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Println(k, m[k])
}

// Go 1.23+ one-liner
for _, k := range slices.Sorted(maps.Keys(m)) {
    fmt.Println(k, m[k])
}
```

### Maps are shared when passed

```go
func add(m map[string]int) {
    m["x"] = 1 // the caller sees this: both variables point at one table
}

a := map[string]int{}
b := a // not a copy
b["k"] = 1
fmt.Println(a["k"]) // 1

c := maps.Clone(a) // a real, shallow copy (Go 1.21)
```

### Values are not addressable

```go
type Point struct{ X, Y int }

m := map[string]Point{"a": {1, 2}}
// m["a"].X = 10 // compile error: cannot assign to struct field in map

p := m["a"] // copy out, change, store back
p.X = 10
m["a"] = p

mp := map[string]*Point{"a": {1, 2}}
mp["a"].X = 10 // fine: the map holds pointers
```

### Keys must be comparable

```go
visited := map[Point]bool{{1, 2}: true} // struct key
counts := map[[2]int]int{{0, 1}: 5}     // array key
// bad := map[[]int]string{}            // compile error: slices are not comparable
```

### Nested maps need their inner map created

```go
nested := map[string]map[string]int{}
// nested["a"]["b"] = 1 // panic: the inner map is nil

if nested["a"] == nil {
    nested["a"] = map[string]int{}
}
nested["a"]["b"] = 1
```

- Iteration order is randomized on purpose, so code cannot come to depend on it.
- Deleting entries during `range` is safe. An entry added during `range` may or may not be visited.
- Maps compare only against nil. Use `maps.Equal(a, b)` to compare contents.
- Allowed key types include numbers, strings, booleans, pointers, arrays and structs of comparable fields. Slices, maps and functions are not allowed.
- You cannot take the address of a map element with `&m[k]`, because entries move when the map grows.

## 10. Maps and concurrency

A map is not safe for concurrent use: if one goroutine writes while another reads or writes, the program crashes.

```go
m := map[int]int{}
for i := 0; i < 100; i++ {
    go func() { m[i] = i }() // fatal error: concurrent map writes
}
```

### Fix 1: a mutex around the map

```go
type SafeMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func NewSafeMap() *SafeMap {
    return &SafeMap{m: make(map[string]int)}
}

func (s *SafeMap) Get(k string) (int, bool) {
    s.mu.RLock() // many readers may hold this at once
    defer s.mu.RUnlock()
    v, ok := s.m[k]
    return v, ok
}

func (s *SafeMap) Set(k string, v int) {
    s.mu.Lock() // one writer, no readers
    defer s.mu.Unlock()
    s.m[k] = v
}
```

### Fix 2: sync.Map

```go
var sm sync.Map

sm.Store("a", 1)
v, ok := sm.Load("a") // v has type any: assert before use
sm.Delete("a")
sm.Range(func(k, v any) bool {
    return true // return false to stop early
})
```

| Approach | Choose it when |
| --- | --- |
| Map with `sync.RWMutex` | The default: type-safe and fast for mixed reads and writes |
| `sync.Map` | Keys are written once and read many times, or goroutines touch disjoint keys |
| Map owned by one goroutine | Other goroutines send requests over a channel |

- The crash is a fatal error, not a panic, so `recover` cannot catch it.
- Concurrent reads with no writer are safe.
- The runtime's detection is best-effort. Run tests with `-race` to find unsynchronized access reliably.

## 11. Map internals

A map variable is a pointer to a hash table managed by the runtime, which explains why passing a map shares it.

- A lookup hashes the key, uses the hash to locate a group of slots, then compares keys to find the match.
- Each map has its own random hash seed, which protects against deliberately engineered collisions.
- When the table gets too full it grows, and entries move to new locations. That is why `&m[k]` is not allowed.
- Lookup, insert and delete are O(1) on average. Iterating is O(n).
- A size hint in `make(map[K]V, n)` allocates room up front and avoids repeated growth.

| Go version | Design |
| --- | --- |
| Before 1.24 | Buckets of 8 key-value pairs with overflow buckets; the table doubled when buckets averaged more than 6.5 entries |
| 1.24 and later | Swiss table: groups of 8 slots with a control word of per-slot metadata, which lets a lookup check a whole group at once |

### How the three types are laid out

| Type | What the variable holds | Effect of assignment |
| --- | --- | --- |
| Array | The elements themselves | Copies all elements |
| Slice | Pointer, length, capacity: 24 bytes | Copies the header; elements are shared |
| Map | One pointer to the runtime table: 8 bytes | Copies the pointer; the table is shared |

A function can therefore add to a caller's map but cannot lengthen a caller's slice: the map's size lives in the shared table, while the slice's length lives in the copied header.

## 12. Interview coding problems

Most slice and map problems combine two pointers on a slice with a map used as a lookup, a counter or a set.

### Reverse a slice in place

```go
func reverse(s []int) {
    for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
        s[i], s[j] = s[j], s[i]
    }
}
// [1 2 3 4] becomes [4 3 2 1]; the caller sees it because elements are shared
```

### Remove duplicates, keeping order

```go
func dedupe(s []int) []int {
    seen := make(map[int]struct{}, len(s))
    out := make([]int, 0, len(s))
    for _, v := range s {
        if _, ok := seen[v]; !ok {
            seen[v] = struct{}{}
            out = append(out, v)
        }
    }
    return out
}
// dedupe([]int{1, 2, 2, 3, 1}) returns [1 2 3]
```

### Two sum

```go
func twoSum(nums []int, target int) (int, int, bool) {
    index := make(map[int]int, len(nums)) // value -> position
    for i, n := range nums {
        if j, ok := index[target-n]; ok {
            return j, i, true
        }
        index[n] = i
    }
    return 0, 0, false
}
// twoSum([]int{2, 7, 11, 15}, 9) returns 0, 1, true
```

### Words sorted by frequency

```go
func byFrequency(text string) []string {
    freq := map[string]int{}
    for _, w := range strings.Fields(strings.ToLower(text)) {
        freq[w]++
    }

    words := make([]string, 0, len(freq))
    for w := range freq {
        words = append(words, w)
    }

    sort.Slice(words, func(i, j int) bool {
        if freq[words[i]] != freq[words[j]] {
            return freq[words[i]] > freq[words[j]] // most frequent first
        }
        return words[i] < words[j] // tie-break, so the output is stable
    })
    return words
}
// byFrequency("go is fun and go is fast") returns [go is and fast fun]
```

### Group anagrams

```go
func groupAnagrams(words []string) [][]string {
    groups := map[[26]int][]string{} // an array of letter counts as the key
    for _, w := range words {
        var key [26]int
        for _, c := range w {
            key[c-'a']++ // assumes lowercase a to z
        }
        groups[key] = append(groups[key], w)
    }

    out := make([][]string, 0, len(groups))
    for _, g := range groups {
        out = append(out, g)
    }
    return out
}
// groupAnagrams([]string{"eat", "tea", "tan", "ate", "nat"})
// returns [[eat tea ate] [tan nat]], in some group order
```

- Each solution is O(n) time, apart from the sort in the frequency problem, which is O(n log n).
- Mention that map iteration order is random whenever your output is built by ranging over a map.
- The anagram solution works because arrays are comparable and so can be map keys; a slice could not be used.

## 13. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| Array vs slice? | An array is a fixed-size value; a slice is a resizable view over an array. |
| Is the array size part of its type? | Yes; `[3]int` and `[4]int` are different types. |
| What is inside a slice? | A pointer, a length and a capacity. |
| len vs cap? | Elements in use; elements available before a reallocation. |
| nil slice vs empty slice? | Both have length 0; only the nil one equals nil and marshals to JSON null. |
| Why assign the result of append? | It may return a header with a new length or a new backing array. |
| How does capacity grow? | Roughly doubling when small, then about 1.25 times. |
| Does slicing copy data? | No; the new slice shares the backing array. |
| What is `s[low:high:max]` for? | Limiting capacity, so an append cannot overwrite shared data. |
| How do you copy a slice? | `slices.Clone(s)` or `copy(dst, src)`. |
| Can you compare two slices with `==`? | No, only against nil. Use `slices.Equal`. |
| What does a missing map key return? | The zero value; use the comma-ok form to detect absence. |
| What happens on a write to a nil map? | A panic. Reads, `len`, `delete` and `range` are fine. |
| Why is map iteration order random? | By design, so code never depends on it. |
| What can be a map key? | Any comparable type; not slices, maps or functions. |
| Why can't you write `m[k].field = v`? | Map values are not addressable. Reassign the value or store pointers. |
| Is a map safe for concurrent use? | No; use a mutex or `sync.Map`. |
| Is deleting during range safe? | Yes. |
| Does passing a map copy it? | No; both variables refer to the same table. |
| How is a map implemented? | As a hash table; a Swiss-table design since Go 1.24. |
