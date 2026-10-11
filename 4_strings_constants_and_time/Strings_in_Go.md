# Strings in Go

Oct 11, 2026 · @HR

## 1. What a string is

A Go string is an immutable sequence of bytes that usually holds UTF-8 text.

```go
var s string        // zero value is "", never nil
s = "hello"
fmt.Println(len(s)) // 5: len returns the number of bytes
```

- Internally a string is a two-word header: a pointer to the bytes and a length, 16 bytes on a 64-bit machine.
- Assigning a string or passing it to a function copies only that header, never the bytes.
- A string can hold any bytes; nothing forces valid UTF-8. Literals are UTF-8 because Go source files are.
- There is no nil string. Check for empty with `s == ""` or `len(s) == 0`.
- Strings are comparable, so they work with `==` and as map keys.

## 2. String literals

Go has two literal forms: interpreted literals in double quotes process escapes, and raw literals in backticks do not.

```go
a := "line1\nline2" // interpreted: \n becomes a newline
b := `C:\temp\new`  // raw: backslashes stay as written
c := `SELECT id
FROM users`         // raw literals can span lines

fmt.Println(len("\n"), len(`\n`)) // 1 2
```

- Use raw literals for regular expressions, SQL, JSON and file paths.
- A raw literal cannot contain a backtick.
- Single quotes make a rune, not a string: `'a'` has type `rune` and value 97.

## 3. Bytes, runes and UTF-8

Indexing a string gives bytes and ranging over it gives runes; this difference is the most asked string question in Go interviews.

```go
s := "Go 世界"

fmt.Println(len(s))                    // 9 bytes: G, o, space, then 3 bytes each for 世 and 界
fmt.Println(utf8.RuneCountInString(s)) // 5 characters
fmt.Println(s[3])                      // 228: the first byte of 世, not the character

for i, r := range s { // i is the byte offset, r is the decoded rune
    fmt.Printf("%d:%c ", i, r) // 0:G 1:o 2:  3:世 6:界
}

for i := 0; i < len(s); i++ { // a classic index loop walks bytes
    fmt.Printf("%x ", s[i]) // 47 6f 20 e4 b8 96 e7 95 8c
}
```

| Term | Type | Meaning |
| --- | --- | --- |
| byte | Alias for `uint8` | One raw byte of the string |
| rune | Alias for `int32` | One Unicode code point |
| UTF-8 | Encoding | Stores a code point in 1 to 4 bytes; ASCII takes 1 |

### Invalid UTF-8

```go
bad := "a\xffb"
fmt.Println(utf8.ValidString(bad)) // false

for _, r := range bad {
    fmt.Printf("%U ", r) // U+0061 U+FFFD U+0062
}
```

- `range` never fails on bad input. Each invalid byte decodes to `utf8.RuneError`, which is U+FFFD.
- To get the nth character, convert to `[]rune` first. `s[n]` is the nth byte.
- One user-visible character can be several runes, as with an emoji plus a skin-tone modifier. Rune count is not always what a user would call length.

## 4. Immutability and conversions

You cannot change a string in place: convert it to `[]byte` or `[]rune`, modify the copy, and convert back.

```go
s := "hello"
// s[0] = 'H' // compile error: cannot assign to s[0]

b := []byte(s)  // copies the bytes
b[0] = 'H'
s2 := string(b) // copies again
fmt.Println(s, s2) // hello Hello: the original is untouched

r := []rune("héllo") // one element per character
r[1] = 'e'
fmt.Println(string(r)) // hello

fmt.Println(string(rune(65))) // A: a rune becomes its character
fmt.Println(strconv.Itoa(65)) // 65: a number becomes its digits
```

| Conversion | Cost | Use it when |
| --- | --- | --- |
| `[]byte(s)` | Copies every byte | You need to modify bytes or call an API that takes bytes |
| `[]rune(s)` | Decodes and allocates 4 bytes per character | You need to index or reverse by character |
| `string(b)` | Copies every byte | You are done modifying |

- Immutability is what makes strings safe to share between goroutines without locks.
- It also lets substrings share memory and makes strings valid map keys.
- `string(65)` compiles but gives "A", and `go vet` flags it. Use `strconv.Itoa` for digits.
- The compiler skips the copy in a few safe cases, such as a map lookup with `m[string(b)]`.

## 5. Concatenation and strings.Builder

Use `+` for a few fixed pieces and `strings.Builder` inside loops, because every `+` allocates a new string.

```go
// Fine: a handful of pieces
full := first + " " + last

// Slow: allocates and copies on every iteration, O(n²) overall
s := ""
for i := 0; i < 1000; i++ {
    s += "x"
}

// Fast: one growing buffer
var sb strings.Builder
sb.Grow(1000) // optional: reserve capacity up front
for i := 0; i < 1000; i++ {
    sb.WriteString("x")
}
out := sb.String() // no extra copy

// Joining a slice you already have
csv := strings.Join([]string{"a", "b", "c"}, ",") // a,b,c

// Formatting mixed types
msg := fmt.Sprintf("%s is %d years old", name, age)
```

| Method | Best for |
| --- | --- |
| `+` | Two or three pieces in one expression |
| `strings.Builder` | Loops and unknown numbers of pieces |
| `strings.Join` | A slice of strings with a separator |
| `fmt.Sprintf` | Mixing types and formatting; slowest of the four |
| `bytes.Buffer` | When you also need to read from the buffer or pass an `io.Reader` |

- A `Builder` has `WriteString`, `WriteByte`, `WriteRune` and `Write`, and it implements `io.Writer`.
- Do not copy a `Builder` after first use; pass a pointer. Copying one and writing to the copy panics.
- The zero value is ready to use: `var sb strings.Builder` needs no constructor.

## 6. The strings package

The `strings` package covers searching, splitting, trimming and case changes; these are the functions worth knowing by heart.

| Call | Result |
| --- | --- |
| `strings.Contains("golang", "lang")` | true |
| `strings.HasPrefix("golang", "go")` | true |
| `strings.HasSuffix("main.go", ".go")` | true |
| `strings.Index("golang", "lang")` | 2, or -1 when absent |
| `strings.Count("cheese", "e")` | 3 |
| `strings.Split("a,b,c", ",")` | \[a b c\] |
| `strings.Fields("  a   b ")` | \[a b\], split on any whitespace |
| `strings.Join([]string{"a", "b"}, "-")` | a-b |
| `strings.ReplaceAll("aaa", "a", "b")` | bbb |
| `strings.ToUpper("go")` | GO |
| `strings.TrimSpace("  go \n")` | go |
| `strings.TrimPrefix("v1.2", "v")` | 1.2 |
| `strings.Repeat("ab", 3)` | ababab |
| `strings.EqualFold("Go", "GO")` | true, case-insensitive |
| `strings.Cut("key=value", "=")` | key, value, true |

### Gotchas

```go
// Trim takes a SET of characters, not a prefix or suffix
fmt.Println(strings.TrimLeft("abcabc-x", "abc"))   // -x
fmt.Println(strings.TrimPrefix("abcabc-x", "abc")) // abc-x

// Splitting an empty string gives one empty element, not zero
fmt.Println(len(strings.Split("", ","))) // 1
fmt.Println(len(strings.Fields("")))     // 0

// Cut (Go 1.18) replaces the Index-then-slice pattern
key, value, found := strings.Cut("host=localhost", "=")
fmt.Println(key, value, found) // host localhost true
```

- Every function returns a new string; the input is never modified.
- `strings.NewReader(s)` turns a string into an `io.Reader` without copying.
- `strings.NewReplacer` applies several replacements in a single pass.
- The `bytes` package mirrors most of these functions for `[]byte`.

## 7. Comparison, substrings and memory

Strings compare byte by byte, and slicing a string is an O(1) operation that shares memory with the original.

### Comparison

```go
fmt.Println("apple" == "apple") // true: same bytes
fmt.Println("apple" < "banana") // true: lexicographic order by byte
fmt.Println("Zebra" < "apple")  // true: 'Z' is 90 and 'a' is 97

fmt.Println(strings.EqualFold("Go", "GO")) // true: case-insensitive
```

Use `strings.EqualFold` for case-insensitive checks rather than lowercasing both sides, which allocates two new strings.

### Substrings

```go
s := "hello, world"
sub := s[7:12]            // world: offsets are in bytes
fmt.Println(s[:5], s[7:]) // hello world

u := "héllo"
fmt.Println(u[1:2])                 // half of é: invalid UTF-8
fmt.Println(string([]rune(u)[1:3])) // él: slice by characters instead
```

### Memory sharing

```go
big := readHugeFile()          // say 100 MB
key := big[:10]                // 10 bytes, but keeps all 100 MB alive
key = strings.Clone(big[:10]) // an independent copy, so big can be collected
```

- A substring is a new header pointing into the same bytes. Nothing is copied.
- That sharing is safe because strings are immutable.
- A small substring keeps the whole original alive. `strings.Clone` (Go 1.18) breaks the link.
- Slicing past the end panics at runtime with "slice bounds out of range".

## 8. strconv and fmt verbs

Use `strconv` to convert between strings and numbers, and `fmt` verbs to control how strings print.

```go
n, err := strconv.Atoi("42")           // 42, nil
_, err = strconv.Atoi("4x")            // error: invalid syntax
i, _ := strconv.ParseInt("ff", 16, 64) // 255: base 16, result fits in 64 bits
f, _ := strconv.ParseFloat("3.14", 64) // 3.14
b, _ := strconv.ParseBool("true")      // true

s := strconv.Itoa(42)            // "42"
bin := strconv.FormatInt(255, 2) // "11111111"
```

| Verb | Example value | Prints | Meaning |
| --- | --- | --- | --- |
| `%s` | "hi" | hi | Plain string |
| `%q` | "hi" | "hi" | Quoted, with escapes shown |
| `%v` | "hi" | hi | Default format for any value |
| `%x` | "hi" | 6869 | Hex of each byte |
| `%c` | 'A' | A | The character for a rune |
| `%d` | 'A' | 65 | The rune as a number |
| `%U` | '世' | U+4E16 | Unicode code point |
| `%T` | "hi" | string | The value's type |
| `%-10s` | "hi" | hi plus 8 spaces | Left-aligned in a width of 10 |

- Always check the error from `Atoi` and the `Parse` functions; bad input is common.
- `strconv.Itoa` is faster than `fmt.Sprintf("%d", n)` because it skips reflection.
- `%q` is useful in logs and test failures because it reveals hidden whitespace.

## 9. Interview coding problems

Most string problems reduce to three tools: a `[]rune` conversion, a frequency map, and two pointers.

### Reverse a string

```go
func reverse(s string) string {
    r := []rune(s) // runes, so multi-byte characters survive
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}
// reverse("héllo") returns "olléh"
```

### Check for a palindrome

```go
func isPalindrome(s string) bool {
    r := []rune(strings.ToLower(s))
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        if r[i] != r[j] {
            return false
        }
    }
    return true
}
// isPalindrome("Level") returns true
```

### Check for an anagram

```go
func isAnagram(a, b string) bool {
    if len(a) != len(b) {
        return false
    }
    count := map[rune]int{}
    for _, r := range a {
        count[r]++
    }
    for _, r := range b {
        count[r]--
        if count[r] < 0 {
            return false
        }
    }
    return true
}
// isAnagram("listen", "silent") returns true
```

### First non-repeating character

```go
func firstUnique(s string) (rune, bool) {
    count := map[rune]int{}
    for _, r := range s {
        count[r]++
    }
    for _, r := range s { // second pass keeps the original order
        if count[r] == 1 {
            return r, true
        }
    }
    return 0, false
}
// firstUnique("swiss") returns 'w', true
```

### Longest substring without repeating characters

```go
func longestUnique(s string) int {
    last := map[rune]int{} // rune -> index just after its last position
    best, start := 0, 0
    for i, r := range []rune(s) {
        if p, ok := last[r]; ok && p > start {
            start = p // move the window past the previous occurrence
        }
        if i-start+1 > best {
            best = i - start + 1
        }
        last[r] = i + 1
    }
    return best
}
// longestUnique("abcabcbb") returns 3
```

- Say out loud whether the input is ASCII-only. If it is, a `[26]int` or `[256]int` array beats a map.
- State the complexity: every solution above is O(n) time.
- If asked to work in place, convert to `[]byte` or `[]rune` first, because strings cannot be modified.

## 10. Rapid-fire questions

One-line answers for a quick revision before an interview.

| Question | Short answer |
| --- | --- |
| What is a string in Go? | An immutable sequence of bytes with a pointer and a length. |
| What does `len(s)` return? | The number of bytes, not characters. |
| How do you count characters? | `utf8.RuneCountInString(s)` or `len([]rune(s))`. |
| byte vs rune? | `byte` is `uint8`, one raw byte; `rune` is `int32`, one Unicode code point. |
| What does `s[i]` return? | The byte at offset i. |
| What does `range` over a string yield? | The byte offset and the decoded rune. |
| Can a string be nil? | No; its zero value is the empty string. |
| Why are strings immutable? | Safe sharing across goroutines, cheap substrings, and valid map keys. |
| How do you modify a string? | Convert to `[]byte` or `[]rune`, change it, convert back. |
| Does `[]byte(s)` copy? | Yes, apart from a few compiler-optimized cases. |
| Best way to concatenate in a loop? | `strings.Builder`. |
| Interpreted vs raw literal? | Double quotes process escapes; backticks keep text as written. |
| What does `string(65)` give? | "A"; use `strconv.Itoa(65)` for "65". |
| How do you compare ignoring case? | `strings.EqualFold(a, b)`. |
| Does a substring copy memory? | No; it shares the original's bytes. Use `strings.Clone` to detach. |
| What happens with invalid UTF-8 in `range`? | Each bad byte yields U+FFFD, the replacement character. |
