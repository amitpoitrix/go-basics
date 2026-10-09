package main

import (
	"fmt"
	"strings"
)

/*
Strings
strings in Go are read-only byte slices, you cannot easily modify a specific character or slice text
safely by index. To fix this, you convert the string into a slice of runes ([]rune)

Different types of strings
1. Interpreted String Literal(Double Quotes "") -> this accept escape sequences like
\n, \t etc
2. Raw String Literals(``)

Byte vs Rune
1. byte = 8 bits. Good for raw data and ASCII characters.
2. rune = 32 bits (alias of int32). Good for universal text, emojis, and global languages.

So rune is basically used for unicode characters

*/

func main2() {
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

/*
Go don't have pre increment/decrement but they've post increment/decrement
++a (X)
a++ this is applicable
*/