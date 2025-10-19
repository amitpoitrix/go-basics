package main

import (
	"fmt"
)

// Slices -> Dynamic Array in Go. Most used construct & useful methods
// Slices are reference types
// Slices are built on top of arrays
// Slices have dynamic size
// Slices have length and capacity
// Length -> no. of elements present in slice
// Capacity -> maximum no. of elements that can be stored in slice without allocating new underlying array
// When capacity is exceeded, a new underlying array is allocated with double the capacity

func main() {
	// Creating a slice using make function
	slice1 := make([]int, 3, 5) // length 3, capacity 5
	fmt.Println("slice1:", slice1)
	fmt.Println("len(slice1):", len(slice1))
	fmt.Println("cap(slice1):", cap(slice1))

	// Creating a slice using slice literal
	slice2 := []int{1, 2, 3, 4, 5}
	fmt.Println("slice2:", slice2)
	fmt.Println("len(slice2):", len(slice2))
	fmt.Println("cap(slice2):", cap(slice2))

	// Slicing an array to create a slice 
	// Syntax: array[startIndex:endIndex] (endIndex is exclusive)
	arr := [5]int{10, 20, 30, 40, 50} // array of 5 integers
	slice3 := arr[1:4] // elements from index 1 to 3
	fmt.Println("slice3 after slicing an array:", slice3)
	fmt.Println("len(slice3):", len(slice3))
	fmt.Println("cap(slice3):", cap(slice3))

	// Appending elements to a slice
	slice4 := []int{1, 2, 3}
	slice4 = append(slice4, 4, 5)
	fmt.Println("slice4 after append:", slice4)

	// Copying slices
	slice5 := make([]int, len(slice4))
	copy(slice5, slice4)
	fmt.Println("slice5 after copy:", slice5)
	slice5[0] = 10
	fmt.Println("slice5 after modification:", slice5)
	fmt.Println("slice4 remains unchanged:", slice4)

	// Multi-dimensional slices
	slice6 := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	fmt.Println("slice6:", slice6)

	// Modifying slices
	slice6[0][0] = 10
	fmt.Println("slice6 after modification:", slice6)

	// Demonstrating that slices are reference types
	slice7 := slice6
	slice7[1][1] = 20
	fmt.Println("slice6 after modifying slice7:", slice6)
	fmt.Println("slice7:", slice7)

	// Creating a slice from another slice
	slice8 := slice4[1:4]
	fmt.Println("slice8:", slice8)

	// Length and Capacity of slice8
	fmt.Println("len(slice8):", len(slice8))
	fmt.Println("cap(slice8):", cap(slice8))	

	// Appending to slice8 and observing changes in slice4
	slice8 = append(slice8, 6, 7)
	fmt.Println("slice8 after append:", slice8)
	fmt.Println("slice4 after slice8 append:", slice4)
}