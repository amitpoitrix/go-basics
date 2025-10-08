package main

import (
	"fmt"
	"slices"
)

// slice -> Dynamic Array in Go & most used construct + useful methods
func main() {
	// uninitialized slice is nil (or Null in other programming language)

	// var num []int
	// fmt.Println(num)
	// fmt.Println(num == nil)
	// fmt.Println(len(num))


	// using make() to create slice

	// this will create slice of initial size & capacity as 2 with value as 0
	// var num = make([]int, 2) // 2nd parameter is initial size, 3rd one is capacity
	// fmt.Println(num)
	// fmt.Println(num == nil)
	// fmt.Println(len(num))
	// // capacity -> maximum no. of elements that can fit in slice if this size increase
	// // than it'll dynamically double the capacity
	// fmt.Println(cap(num))

	// using make() to create slice with capacity
	// var num = make([]int, 2, 5) // because of intial size 2, 2 zeroed value will be appended
	// num = append(num, 1) // Adding elements in slice using append()
	// num = append(num, 2)
	// num = append(num, 3)
	// num = append(num, 4)

	// num[0] = 1 // Assigning value at 0th index

	// fmt.Println(num)
	// fmt.Println(len(num))
	// fmt.Println(cap(num))


	// Another way of creating slice with 0 initial size & dynamic capacity
	// nums := []int{}
	// nums = append(nums, 1)
	// nums = append(nums, 2)
	// nums = append(nums, 3)
	
	// fmt.Println(nums)
	// fmt.Println(len(nums))
	// fmt.Println(cap(nums))


	// copy slice from one to another using copy(destination[], source[])
	// var nums1 = make([]int, 0, 5)
	// nums1 = append(nums1, 2)
	// var nums2 = make([]int, len(nums1))

	// copy(nums2, nums1)

	// fmt.Println(nums1, nums2)


	// slice Operator arr[:]
	// var nums = []int{1,2,3,4,5}
	// fmt.Println(nums[0:2])
	// fmt.Println(nums[:2])
	// fmt.Println(nums[3:])


	// comparing 2 slice using slices pkg Equal()
	var num1 = []int{1,2,1}
	var num2 = []int{1,2,1}
	fmt.Println(slices.Equal(num1, num2))


	// 2-D slices
	var nums = [][]int{{1,2,3}, {4,5}}
	fmt.Println(nums)
}