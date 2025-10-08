package main

import "fmt"

// func sum(a int, b int) int {
func sum(a, b int) int { // prior to b & b itself will be of type int
	return a + b
}


// In Go, function can return more than 1 values, just need to mention return type under bracket () if mentioning
// more than one return type
func languages() (string, string, bool) {
	return "golang", "FORTRAN", false
}


// In Go function are first class citizens
func returningFunction(fn func(a int) int) func(b int) int {
	// returning a function(i.e., anonymous function)
	result := fn(2)

	return func(c int) int {
		return c + result
	}
}

// below main() is the starting point and Go compiler first call this main()
func main() {
	result := sum(2, 3)
	fmt.Println(result)

	// lang1, lang2, isObjectOriented := languages()
	lang1, lang2, _ := languages() // if don't want to use any return value from function than use _
	fmt.Println(lang1, lang2)
	

	fmt.Println("Using function in Go as first class citizens")
	
	// Intializing anonymous function and storing in variable fnVar
	fnVar := func(x int) int {
		return x + 25
	}
	
	// first calling a function that returns a function
	resultFromFunc := returningFunction(fnVar)

	// And that previously returned function is called to get the final result
	finalresult := resultFromFunc(3)

	fmt.Println(finalresult)
}