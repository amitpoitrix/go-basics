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

	// array declaration
	// var arr [3]int

	// Initializing and declaring array at the same time
	var arr2 = [3]int{6, 6 ,6}
	fmt.Println("arr2:", arr2)

	// Short Declaration
	arr := [3]int{1,2,3}
	fmt.Println("arr:", arr)

	// Ellipsis[...] - array will automatically takes the size no need to explicitly define
	arr1 := [...]int{2, 3, 4, 5}
	fmt.Println("arr1:", arr1)

	// 2D arrays
	twoDArr := [2][2]int{{6,7}, {2,3}}
	fmt.Println(twoDArr)

	// So in Go arrays are
	// fixed size, that is predictable
	// Memory Optimization
	// Constant time access


	// for loop in array using range
	for index, value := range arr {
		fmt.Println(index, value)
	}
}