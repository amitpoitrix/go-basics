package main

import (
	"fmt"
	"strings"
)

/*
Strings
In Go, a string is a read-only slice of bytes. Unlike many languages where strings are arrays of characters, 
Go strings are fundamentally treated as raw, immutable binary data, typically encoded in UTF-8.
Under the hood, a string header is a tiny struct (only 16 bytes on 64-bit systems) containing a pointer to a 
backing byte array and a length.

1. Key Characteristics of Go Strings
• Immutable: You cannot change the bytes inside a string. Once created, its backing array cannot be modified.
• UTF-8 by Default: Source code in Go is UTF-8 text, so string literals are naturally interpreted as UTF-8 sequences.
• Zero Value: The zero value of a string is the empty string "". It does not point to any memory (nil), and its length is 0.


2. The Difference: Bytes vs. Runes (Characters)
Because strings are byte slices, indexing a string returns bytes, not characters. To represent a Unicode 
character, Go uses the concept of a rune (which is just an alias for int32).
If your string contains non-ASCII characters (like emojis or accents), a single character will take up multiple bytes:

s := "Hello 👋"

// len() counts BYTES, not characters!
fmt.Println(len(s)) // Output: 10 (Hello = 5 bytes, Space = 1 byte, 👋 = 4 bytes)

// Indexing returns the raw byte byte at that position
fmt.Println(s[0]) // Output: 72 (ASCII code for 'H')

Byte vs Rune
1. byte = 8 bits (byte is an alias of uint8) Good for raw data and ASCII characters.
2. rune = 32 bits (rune is an alias of int32). Good for universal text, emojis, and global languages.

So rune is basically used for unicode characters

3. How to Properly Loop Over Strings
Because of the byte vs. rune distinction, how you loop over a string changes what you get back:
• Using a standard for loop: Iterates byte-by-byte. This will break multi-byte Unicode characters.
• Using a for range loop: Go automatically decodes the UTF-8 bytes into runes and provides the byte 
index where each character starts.

s := "Go" // 'G' is 1 byte, 'o' is 2 bytes

for _, runeValue := range s {
    fmt.Printf("%c ", runeValue) // Output: G o
}


Different types of strings
1. Interpreted String Literal(Double Quotes "") -> this accept escape sequences like \n, \t etc
2. Raw String Literals(``)

*/

func main() {
	a := "Hello \n world"
	fmt.Println("A---", a)

	b := `Hello \n world`
	fmt.Println("B---", b)

	s := "hello"
	r := strings.ToUpper(s)
	fmt.Println("r---", r)

	v := strings.Contains(s, "ll")
	fmt.Println("v---", v)

	split := strings.Split("a,b,c", ",")
	fmt.Println("Split---", split)

	rep := strings.Replace(b, "world", "you", 1)
	fmt.Println("rep---", rep)

	trim := strings.TrimSpace("    Amit    ")
	fmt.Println("Trim---", trim)
}

func init() {
	fmt.Println("Hello I'm from init() and I always runs first before main()")
}

/*
Go don't have pre increment/decrement but they've post increment/decrement
++a (X)
a++ this is applicable
*/