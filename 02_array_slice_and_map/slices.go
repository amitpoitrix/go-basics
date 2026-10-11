package main

import (
	"fmt"
	// "slices"
)

/*
slice
Dynamic Array in Go & most used construct + useful methods
A slice is a dynamic view over an underlying array
Backing array

Pointer - starting index of a slice in an array
Length - number of elements in the slice
Capacity - max elements it can grow before a new array is created  
*/

func main() {
	// uninitialized slice is nil (or Null in other programming language)

	var num []int
	fmt.Println(num) // [] - see its not printing nil its just empty array []
	fmt.Println(num == nil) // true
	fmt.Println(len(num)) // 0

	num2 := []int{}
	fmt.Println(num2) // [] - see its not printing nil its just empty array []
	fmt.Println(num2 == nil) // true
	fmt.Println(len(num2)) // 0


	// Method 1: using make() to create slice

	// this will create slice of initial size & capacity as 2 with value as 0
	// var num = make([]int, 2)	// 2nd parameter is initial size, 3rd one is capacity
	// fmt.Println(num) // [0 0]
	// fmt.Println(num == nil) // false
	// fmt.Println(len(num)) // 2
	// // capacity -> maximum no. of elements that can fit in slice if this size increase
	// // than it'll dynamically double the capacity
	// fmt.Println(cap(num)) // 2


	// Method 2: using make() to create slice with capacity

	// var num = make([]int, 2, 5) // slice of length 2 with capacity 5
	// // because of initial size 2, 2 zeroed value will be appended
	
	// // Adding elements in slice using append()
	// num = append(num, 1) 
	// num = append(num, 2)
	// num = append(num, 3)
	// num = append(num, 4)

	// num[0] = 1 // Assigning value at 0th index

	// fmt.Println(num)	// [1 0 1 2 3 4]
	// fmt.Println(len(num))	// 6
	// fmt.Println(cap(num))	// 10


	// Method 3: Another way of creating slice with 0 initial size & dynamic capacity
	// nums := []int{} // without using make()
	// nums = append(nums, 1)
	// nums = append(nums, 2)
	// nums = append(nums, 3)
	
	// fmt.Println(nums)	// [1 2 3]
	// fmt.Println(len(nums))	// 3
	// fmt.Println(cap(nums))	// 4 capacity is 4 because size is keep on doubling 1, 2, 4, 8, 16....


	// copy slice from one to another using copy(destination[], source[])
	// nums1 := make([]int, 0, 5)
	// nums1 = append(nums1, 2)
	// nums2 := make([]int, len(nums1)) 
	// fmt.Println(nums2)	// [0]

	// copy(nums2, nums1) // copy all values of nums1 to nums2

	// nums2 = append(nums2, 12)

	// fmt.Println("len(num2)", len(nums2)) // 2
	// fmt.Println("cap(num2)", cap(nums2)) // 2

	// fmt.Println(nums1, nums2) // [2] [2 12]


	// slice Operator arr[:]
	// nums := []int{1,2,3,4,5}
	// fmt.Println(nums[0:2]) // [1 2]
	// fmt.Println(nums[:2]) // [1 2]
	// fmt.Println(nums[3:]) // [4 5]


	// comparing 2 slice using slices pkg Equal()
	// var num1 = []int{1,2,1}
	// var num2 = []int{1,2,1}
	// fmt.Println(slices.Equal(num1, num2))


	// // 2-D slices
	// var nums = [][]int{{1,2,3}, {4,5}}
	// fmt.Println(nums)
}