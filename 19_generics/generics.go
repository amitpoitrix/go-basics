package main

import "fmt"

/*
In Go,
Generics are the newly introduced in Go 1.18 and works as same as other programming language

In order to use generics in function just use [T int | string] in front of function name followed by type as T 
we can use any or interface{} also instead of string | int but that is bad practice so instead of that
we can use comparable
*/

// func printSlice[T int | string | bool](nums []T) {
func printSlice[T comparable](nums []T) {
	for _, num := range nums {
		fmt.Println(num)
	}
}


/* Using generics with struct */
type stack[T comparable] struct {
	elements []T
}

func main() {
	// nums := []int{1,2,3}
	// names := []string{"amit", "kumar", "behera"}
	isAllowed := []bool{true, false}

	printSlice(isAllowed)


	// Specify type in [] while using Generics struct during instance creation
	myCustomStack := stack[bool] {
		elements: isAllowed,
	}

	fmt.Println(myCustomStack)
}