# Structs, Pointers, Methods and Embedding in Go

## 1. Struct basics

A struct groups named fields into one value type, and it is Go's replacement for a class.

```go
type User struct {
    ID    int
    Name  string
    Email string
    admin bool // lowercase: visible only inside this package
}

var u1 User                             // zero value: every field zeroed
u2 := User{ID: 1, Name: "Asha"}         // named fields: the preferred form
u3 := User{2, "Ravi", "r@x.com", false} // positional: every field, in order
u4 := &User{ID: 3}                      // pointer to a new struct
u5 := new(User)                         // pointer to a zero-valued struct

u2.Email = "a@x.com"
fmt.Printf("%v\n", u2)  // {1 Asha a@x.com false}
fmt.Printf("%+v\n", u2) // {ID:1 Name:Asha Email:a@x.com admin:false}
```

### Anonymous structs

```go
point := struct{ X, Y int }{1, 2} // a one-off type with no name

tests := []struct {
    in   string
    want int
}{
    {"a", 1},
    {"bb", 2},
} // the usual shape of a table-driven test
```

- A struct is a value: assigning it or passing it to a function copies every field.
- The zero value is usable without a constructor; there is no nil struct.
- Uppercase fields are exported from the package; lowercase fields are private to it.
- Named-field literals survive when someone adds a field later. Positional literals break.
- `%+v` prints field names and `%#v` prints full Go syntax, which helps when debugging.

## 2. Pointers

A pointer holds the address of a value: `&` takes an address and `*` reads or writes through it.

```go
x := 10
p := &x         // p has type *int and holds the address of x
fmt.Println(*p) // 10: dereference
*p = 20
fmt.Println(x)  // 20

var np *int // the zero value of a pointer is nil
// fmt.Println(*np) // panic: nil pointer dereference

u := &User{Name: "Asha"}
u.Name = "Ravi" // shorthand for (*u).Name: Go dereferences for you
```

### Passing by value vs by pointer

```go
func rename(u User)     { u.Name = "X" } // works on a copy
func renamePtr(u *User) { u.Name = "X" } // works on the caller's struct

u := User{Name: "Asha"}
rename(u)
fmt.Println(u.Name) // Asha
renamePtr(&u)
fmt.Println(u.Name) // X
```

### Reassigning a pointer vs writing through it

```go
func resetBroken(p *User) { p = &User{} } // changes only the local copy of the pointer
func resetFixed(p *User)  { *p = User{} } // overwrites the value the pointer refers to
```

### Returning a pointer to a local variable

```go
func NewUser(name string) *User {
    u := User{Name: name}
    return &u // safe: the compiler moves u to the heap (escape analysis)
}
```

| Use a pointer when | Use a value when |
| --- | --- |
| The function or method must modify the original | The data is small and read-only, such as a point or a time |
| The struct is large and copying is costly | You want a guarantee that the callee cannot change it |
| The struct holds a mutex or must not be copied | You want to avoid nil checks |
| nil is a meaningful "not present" state | Fewer heap allocations matter |

- Go passes everything by value. Passing a pointer copies the address, which is why the callee reaches the same data.
- There is no pointer arithmetic outside the `unsafe` package.
- Dereferencing a nil pointer panics at runtime; the compiler does not catch it.
- `new(T)` and `&T{}` both return a pointer to a zeroed `T`. The literal form can also set fields.

## 3. Methods

A method is a function with a receiver, and the receiver's kind decides whether the method works on a copy or on the original.

```go
type Counter struct{ n int }

func (c Counter) Value() int { return c.n } // value receiver: reads a copy
func (c *Counter) Inc()      { c.n++ }      // pointer receiver: modifies the original
func (c Counter) IncBroken() { c.n++ }      // bug: increments a copy and discards it

c := Counter{}
c.Inc()       // Go rewrites this as (&c).Inc()
c.IncBroken() // no effect
fmt.Println(c.Value()) // 1

p := &c
fmt.Println(p.Value()) // Go rewrites this as (*p).Value()
```

### Choosing the receiver

| Use a pointer receiver when | Use a value receiver when |
| --- | --- |
| The method modifies the receiver | The type is small and immutable |
| The struct is large | The type is a map, function or channel |
| The struct contains a `sync.Mutex` or similar | The type is a basic type such as `int` or `string` |
| Any other method on the type already uses one | Every method on the type can use one |

Keep receivers consistent: if one method needs a pointer receiver, give all methods on that type a pointer receiver.

### Methods on any named type

```go
type Celsius float64

func (c Celsius) Fahrenheit() float64 { return float64(c)*9/5 + 32 }

type IntList []int

func (l IntList) Sum() int {
    total := 0
    for _, v := range l {
        total += v
    }
    return total
}

fmt.Println(Celsius(100).Fahrenheit()) // 212
fmt.Println(IntList{1, 2, 3}.Sum())    // 6
```

### Nil receivers

```go
type Node struct {
    Val  int
    Next *Node
}

func (n *Node) Len() int {
    if n == nil { // a method can be called on a nil pointer
        return 0
    }
    return 1 + n.Next.Len()
}

var head *Node
fmt.Println(head.Len()) // 0, no panic
```

### The Stringer interface

```go
func (u User) String() string { return fmt.Sprintf("%s <%s>", u.Name, u.Email) }

fmt.Println(User{Name: "Asha", Email: "a@x.com"}) // Asha <a@x.com>
```

- You can declare methods only on types defined in the same package, so not on `int` or on another package's types directly.
- Go has no method overloading: one name, one method per type.
- A method value such as `f := c.Inc` binds the receiver, so `f()` can be called later or passed as a callback.
- A method called on a nil pointer panics only when it touches a field without checking for nil.

## 4. Method sets and interfaces

A value of type `T` carries only the value-receiver methods, while a `*T` carries all of them; this decides which one satisfies an interface.

| Method declared with | In the method set of `T` | In the method set of `*T` |
| --- | --- | --- |
| Value receiver `(t T)` | Yes | Yes |
| Pointer receiver `(t *T)` | No | Yes |

```go
type Shape interface{ Area() float64 }

type Square struct{ side float64 }

func (s *Square) Area() float64 { return s.side * s.side } // pointer receiver

var a Shape = &Square{2} // ok: *Square has Area
// var b Shape = Square{2} // compile error: Square does not implement Shape
//                         // (method Area has pointer receiver)
```

### Why a plain variable still works

```go
sq := Square{2}
sq.Area() // fine: sq is addressable, so Go calls (&sq).Area()

// Square{2}.Area() // compile error: a literal is not addressable
(&Square{2}).Area() // fine

shapes := map[string]Square{"a": {2}}
// shapes["a"].Area() // compile error: map elements are not addressable
```

- Calling a pointer method on a variable works because Go can take the variable's address automatically.
- A value stored in an interface is not addressable, so Go cannot do that there. Hence the rule.
- With all value receivers, both `T` and `*T` satisfy the interface.
- A compile-time check that a type implements an interface: `var _ Shape = (*Square)(nil)`.

## 5. Constructors and functional options

Go has no constructors: by convention a function named `NewX` builds the value, sets defaults and validates input.

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
}

type Option func(*Server) // an option is a function that adjusts the server

func WithPort(p int) Option {
    return func(s *Server) { s.port = p }
}

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServer(host string, opts ...Option) (*Server, error) {
    if host == "" {
        return nil, errors.New("host is required")
    }
    s := &Server{host: host, port: 8080, timeout: 30 * time.Second} // defaults
    for _, opt := range opts {
        opt(s)
    }
    return s, nil
}

srv, err := NewServer("localhost", WithPort(9090), WithTimeout(5*time.Second))
```

- Unexported fields plus an exported `NewX` give you encapsulation: other packages must go through the constructor.
- Return `(*T, error)` when construction can fail.
- Functional options suit types with many optional settings. New options can be added without breaking callers.
- For two or three settings, a plain config struct parameter is simpler.
- If the zero value of a type is already useful, as with `sync.Mutex` and `bytes.Buffer`, it needs no constructor.

## 6. Embedding

Embedding places one type inside another without a field name, and the inner type's fields and methods are promoted to the outer type.

```go
type Animal struct{ Name string }

func (a Animal) Speak() string    { return a.Name + " makes a sound" }
func (a Animal) Describe() string { return "I say: " + a.Speak() }

type Dog struct {
    Animal // embedded: no field name
    Breed  string
}

func (d Dog) Speak() string { return d.Name + " barks" } // shadows Animal.Speak

d := Dog{Animal: Animal{Name: "Rex"}, Breed: "Lab"}

fmt.Println(d.Name)           // Rex: promoted field, same as d.Animal.Name
fmt.Println(d.Speak())        // Rex barks: the outer method wins
fmt.Println(d.Animal.Speak()) // Rex makes a sound: the inner one is still reachable
fmt.Println(d.Describe())     // I say: Rex makes a sound
```

The last line is the key interview point. `Describe` belongs to `Animal`, so inside it `a.Speak()` is always `Animal.Speak`. Go has no virtual dispatch through embedding.

### Embedding is not inheritance

```go
func feed(a Animal) {}

// feed(d)     // compile error: a Dog is not an Animal
feed(d.Animal) // pass the embedded value explicitly

type Speaker interface{ Speak() string }

var s Speaker = d // fine: polymorphism in Go comes from interfaces
```

### Embedding an interface

```go
type Store interface {
    Get(id int) (string, error)
    Save(id int, v string) error
}

// Decorator: override one method, delegate the rest automatically
type LoggingStore struct {
    Store
}

func (l LoggingStore) Save(id int, v string) error {
    log.Println("saving", id)
    return l.Store.Save(id, v)
}

// Test stub: implement only what the test needs
type stubStore struct {
    Store // nil: calling an unimplemented method panics
}

func (stubStore) Get(id int) (string, error) { return "stub", nil }
```

### Name conflicts

```go
type A struct{ ID int }
type B struct{ ID int }

type C struct {
    A
    B
}

var c C
// fmt.Println(c.ID) // compile error: ambiguous selector c.ID
fmt.Println(c.A.ID)  // fine: name the path
```

- Embedding is composition with shorter syntax: a `Dog` has an `Animal`; it is not one.
- The outer type satisfies every interface the embedded type satisfies, through the promoted methods.
- Embedding `T` promotes its pointer-receiver methods only to `*Outer`. Embedding `*T` promotes all methods to both.
- An embedded pointer is nil until you set it, and using it before then panics.
- Embedding `sync.Mutex` exposes `Lock` and `Unlock` on your type. For exported types, prefer a named field such as `mu sync.Mutex`.

## 7. Struct tags and JSON

A struct tag is a string of metadata attached to a field, which packages such as `encoding/json` read through reflection.

```go
type Product struct {
    ID       int      `json:"id"`
    Name     string   `json:"name"`
    Price    float64  `json:"price,omitempty"` // left out when zero
    Secret   string   `json:"-"`               // never encoded
    Tags     []string `json:"tags"`
    internal string   // unexported: invisible to encoding/json
}

p := Product{ID: 1, Name: "Pen", Secret: "x"}
b, _ := json.Marshal(p)
fmt.Println(string(b)) // {"id":1,"name":"Pen","tags":null}

var q Product
err := json.Unmarshal([]byte(`{"id":2,"name":"Ink","extra":true}`), &q)
fmt.Println(q.ID, q.Name, err) // 2 Ink <nil>: unknown keys are ignored
```

### Telling "missing" from "zero"

```go
type Patch struct {
    Name *string `json:"name"` // nil when the key is absent
    Age  *int    `json:"age"`
}

var p Patch
json.Unmarshal([]byte(`{"age":0}`), &p)
fmt.Println(p.Name == nil, *p.Age) // true 0: age was sent as zero, name was not sent
```

### Reading a tag yourself

```go
f, _ := reflect.TypeOf(Product{}).FieldByName("Name")
fmt.Println(f.Tag.Get("json")) // name
```

- Only exported fields are encoded or decoded. A lowercase field is silently skipped.
- `Unmarshal` needs a pointer, because it must write into your value.
- A nil slice encodes as `null` and an empty slice as `[]`.
- Fields of an embedded struct are flattened into the outer JSON object.
- Other common tag keys are `db`, `yaml`, `xml` and `validate`. A field can carry several, separated by spaces.
- Tags are plain strings checked only at runtime. `go vet` catches malformed ones.

## 8. Copying, comparing and memory layout

Assigning a struct makes a shallow copy: plain fields are duplicated, but slices, maps and pointers inside still share their data.

```go
type Team struct {
    Name    string
    Members []string
}

a := Team{Name: "core", Members: []string{"asha"}}
b := a // shallow copy

b.Name = "infra"      // independent: strings are values
b.Members[0] = "ravi" // shared: both slices point at one array

fmt.Println(a.Name, a.Members) // core [ravi]

c := a
c.Members = slices.Clone(a.Members) // deep copy: clone each reference field
```

### Comparing structs

```go
type Point struct{ X, Y int }

fmt.Println(Point{1, 2} == Point{1, 2}) // true: compared field by field

visited := map[Point]bool{{1, 2}: true} // comparable structs can be map keys

// Team{} == Team{} // compile error: the slice field makes Team non-comparable
fmt.Println(reflect.DeepEqual(a, c)) // true: slow, but works on anything
```

### The empty struct

```go
fmt.Println(unsafe.Sizeof(struct{}{})) // 0

set := map[string]struct{}{} // a set with no memory spent on values
done := make(chan struct{})  // a channel used only as a signal
```

### Field order and padding

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

fmt.Println(unsafe.Sizeof(Bad{}), unsafe.Sizeof(Good{})) // 24 16
```

- A struct is comparable only if every field is. Slice, map and function fields break `==`.
- Each field is aligned to its own size, so the compiler inserts padding. Ordering fields from largest to smallest minimizes it.
- Padding matters when you hold millions of structs; otherwise prefer the order that reads best.
- Pointers to two different zero-size values may or may not be equal, so never rely on comparing them.

## 9. Gotchas

These five mistakes account for most struct and pointer bugs, and for most "what does this print?" questions.

### A nil pointer inside an interface is not nil

```go
type MyErr struct{}

func (*MyErr) Error() string { return "boom" }

func do() error {
    var e *MyErr // nil pointer
    return e     // the interface now holds (type *MyErr, value nil)
}

fmt.Println(do() == nil) // false
```

An interface is nil only when both its type and its value are nil. Return a literal `nil` on the success path.

### Copying a struct that holds a mutex

```go
type Counter struct {
    mu sync.Mutex
    n  int
}

func (c Counter) Inc() { // value receiver: copies the mutex on every call
    c.mu.Lock()
    c.n++ // and increments the copy
    c.mu.Unlock()
}
```

Use a pointer receiver. `go vet` reports this as "passes lock by value".

### range hands you a copy

```go
users := []User{{Name: "Asha"}, {Name: "Ravi"}}

for _, u := range users {
    u.Name = "X" // changes a copy
}

for i := range users {
    users[i].Name = "X" // changes the element
}
```

A slice of pointers, `[]*User`, avoids the copy: `u.Name = "X"` then reaches the original.

### Struct values in a map cannot be modified in place

```go
m := map[string]User{"a": {Name: "Asha"}}
// m["a"].Name = "X" // compile error

u := m["a"]
u.Name = "X"
m["a"] = u // copy out, change, store back; or use map[string]*User
```

### Dereferencing nil

```go
type Order struct{ Customer *User }

o := Order{}
// fmt.Println(o.Customer.Name) // panic: nil pointer dereference

if o.Customer != nil {
    fmt.Println(o.Customer.Name)
}
```

- Pointer fields default to nil. Check them, or set them in the constructor.
- Mixing value and pointer receivers on one type leads to surprises with interfaces. Pick one kind.
- Returning a struct by value gives the caller a copy, so changes to it never reach your internal state.

## 10. Interview coding problems

These four exercises cover what interviewers check: pointer manipulation, receiver choice, thread safety, and polymorphism through interfaces.

### Linked list with reverse

```go
type Node struct {
    Val  int
    Next *Node
}

type LinkedList struct {
    head *Node
    size int
}

func (l *LinkedList) Push(v int) { // insert at the front
    l.head = &Node{Val: v, Next: l.head}
    l.size++
}

func (l *LinkedList) Reverse() {
    var prev *Node
    cur := l.head
    for cur != nil {
        next := cur.Next // remember the rest of the list
        cur.Next = prev  // point this node backwards
        prev = cur
        cur = next
    }
    l.head = prev
}

func (l *LinkedList) Values() []int {
    out := make([]int, 0, l.size)
    for n := l.head; n != nil; n = n.Next {
        out = append(out, n.Val)
    }
    return out
}

// Push 1, 2, 3 gives [3 2 1]; after Reverse the values are [1 2 3]
```

### Generic stack

```go
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

func (s *Stack[T]) Len() int { return len(s.items) }
```

### Thread-safe bank account

```go
type Account struct {
    mu      sync.Mutex
    balance int
}

func (a *Account) Deposit(n int) {
    a.mu.Lock()
    defer a.mu.Unlock()
    a.balance += n
}

func (a *Account) Withdraw(n int) error {
    a.mu.Lock()
    defer a.mu.Unlock()
    if n > a.balance {
        return errors.New("insufficient funds")
    }
    a.balance -= n
    return nil
}

func (a *Account) Balance() int {
    a.mu.Lock()
    defer a.mu.Unlock()
    return a.balance
}
```

### Polymorphism without inheritance

```go
type Shape interface {
    Area() float64
}

type Rect struct{ W, H float64 }

func (r Rect) Area() float64 { return r.W * r.H }

type Circle struct{ R float64 }

func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

func totalArea(shapes []Shape) float64 {
    total := 0.0
    for _, s := range shapes {
        total += s.Area() // dispatched on the concrete type at runtime
    }
    return total
}

fmt.Println(totalArea([]Shape{Rect{2, 3}, Circle{1}})) // about 9.14
```

- Be ready to say why each method above uses a pointer or a value receiver.
- For the account, expect the follow-up "how would you transfer between two accounts without a deadlock?" Lock the two in a fixed order, for example by account ID.
- For the list, expect "find the middle" and "detect a cycle", which both use a slow and a fast pointer.

## 11. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| Does Go have classes? | No; it has structs with methods, interfaces and embedding. |
| Is a struct a value or a reference type? | A value type; assignment copies every field. |
| `new(T)` vs `&T{}`? | Both return a pointer to a zeroed T; the literal can also set fields. |
| Is Go pass by value or by reference? | Always by value; a pointer argument is a copied address. |
| Is returning a pointer to a local variable safe? | Yes; escape analysis moves the variable to the heap. |
| Does Go have pointer arithmetic? | No, except through the `unsafe` package. |
| Value or pointer receiver? | Pointer if the method mutates, the struct is large, or it holds a lock. |
| Can you call a pointer method on a value? | Yes, if the value is addressable; Go takes its address. |
| What is the method set of T and of \*T? | T has value-receiver methods only; \*T has both kinds. |
| Can a method be called on a nil pointer? | Yes; it panics only if it touches a field without a nil check. |
| Can you add methods to `int`? | Not directly; define a named type such as `type MyInt int`. |
| Does Go support method overloading? | No. |
| Is embedding inheritance? | No; it is composition with promoted fields and methods. |
| What happens when the outer type defines the same method? | It shadows the embedded one, which stays reachable by its full path. |
| How does Go achieve polymorphism? | Through interfaces, satisfied implicitly. |
| When are two structs comparable? | When every field is comparable. |
| Is a struct copy deep or shallow? | Shallow; slices, maps and pointers inside are shared. |
| What is a struct tag? | Field metadata read through reflection, as in `json:"name"`. |
| Why is a lowercase field missing from JSON? | It is unexported, so `encoding/json` cannot see it. |
| What is the size of an empty struct? | Zero bytes. |
