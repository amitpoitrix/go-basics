package main

import "fmt"

// Here in parameter ... should be left of type (i.e., ...int)
func sum(nums ...int) int {
	total := 0

	for _, num := range nums {
		total += num
	}

	return total
}

func anyTypeFunc(element ...any) string {
	fmt.Println(element...)

	return "All izz well"
}


// In Go, variadic functions are those function which can take n no. of variables as arguements
func main() {
	// fmt.Println() is famous example of variadic function
	fmt.Println(2,3,4,"hello")

	result := sum(1,2,3,4,5)
	fmt.Println(result)


	// using slice as parameter to pass as arguement to sum() variadic function
	slice := []int{2,3,4,5,6}

	result1 := sum(slice...) // while passing slice put ... at right of slice (i.e., nums...)
	fmt.Println(result1)


	// variadic function with any type
	result2 := anyTypeFunc("hello", "varidic", "function", "number", 3)
	fmt.Println(result2)
}