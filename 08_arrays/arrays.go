package main

import "fmt"

// In Go, arrays are numbered sequence of specific length
func main() {
	// In Go, default values or Zeroed values
	// int -> 0, string -> "", bool -> false

	// int arrays
	var nums [4]int
	nums[2] = 1
	fmt.Println(len(nums))
	fmt.Println(nums)

	// bool array
	var vars [4]bool
	vars[1] = true
	fmt.Println(vars)

	// string array
	var names [4]string
	names[3] = "amit"
	fmt.Println(names)


	// Initializing and declaring array at the same time
	arr := [3]int{1,2,3}
	fmt.Println(arr)

	// 2D arrays
	twoDArr := [2][2]int{{6,7}, {2,3}}
	fmt.Println(twoDArr)


	// So in Go arrays are
	// fixed size, that is predictable
	// Memory Optimization
	// Constant time access
}