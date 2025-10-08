package main

import "fmt"

// Closures in Go is same as in JavaScript, as it still make the reference of variables being declared 
// in outer scope

func counter() func() int {
	counter := 0

	return func() int {
		counter++

		return counter
	}
}

func main() {
	increment := counter()

	fmt.Println(increment())
	fmt.Println(increment())
}