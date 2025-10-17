package main

import "fmt"

// In Go, arrays are numbered sequence of specific length
// and specific type of elements
// Arrays are value types
func main() {
	// In Go, below are default values or Zeroed values and array contains them by default
	// int -> 0, string -> "", bool -> false
	// when we declare an array
	// So if we declare an array of int, it will contain 0s
	

	// Declaring and initializing
	// An array of 4 integers
	var arr1 [4]int = [4]int{1,2,3,4}
	fmt.Println("arr1:" , arr1)

	// Another way of declaring and initializing
	var arr2 = [4]int{1,2,3,4}
	fmt.Println("arr2:" , arr2)

	// Yet another way of declaring and initializing
	arr3 := [4]int{1,2,3,4}
	fmt.Println("arr3:" , arr3)

	// Declaring an array of size 4
	// and initializing only 2 elements
	// rest will be zeroed values
	arr4 := [4]int{1,2}
	fmt.Println("arr4:" , arr4)

	// If we don't know the size of array at compile time
	// we can use ... to let the compiler decide the size
	arr5 := [...]int{1,2,3,4,5}
	fmt.Println("arr5:" , arr5)

	// Accessing elements of array
	fmt.Println("arr5[0]:", arr5[0])

	// Modifying elements of array
	arr5[0] = 10
	fmt.Println("arr5:" , arr5)

	// Length of array
	fmt.Println("len(arr5):", len(arr5))

	// float array
	var floatArr [4]float64
	floatArr[0] = 1.1
	floatArr[1] = 2.2
	fmt.Println("floatArr:", floatArr)

	// bool array
	var boolArr [4]bool
	boolArr[1] = true
	fmt.Println("boolArr:", boolArr)

	// string array
	var names [4]string
	names[3] = "amit"
	fmt.Println("names:", names)

	// 2D arrays
	twoDArr := [2][2]int{{6,7}, {2,3}}
	fmt.Println("twoDArr:", twoDArr)

	// Arrays are value types
	// Deep Copy
	arr6 := [3]int{1,2,3}
	arr7 := arr6 // arr7 is a deep copy of arr6
	arr7[0] = 10 // modifying arr7 won't affect arr6
	fmt.Println("arr6:", arr6)
	fmt.Println("arr7:", arr7)

	// Shallow Copy
	// If we want to modify the original array
	// we can use pointers
	arr8 := [3]int{1,2,3}
	arr9 := &arr8
	arr9[0] = 10
	fmt.Println("arr8:", arr8)
	fmt.Println("*arr9:", *arr9)

	// So in Go arrays are
	// fixed size, that is predictable
	// Memory Optimization
	// Constant time access
}   