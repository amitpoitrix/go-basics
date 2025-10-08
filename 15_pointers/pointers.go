package main

import "fmt"

// In Go, parameters are pass by value
// func changeNum(num int) {
// 	num = 5

// 	fmt.Println("In changeNum", num)
// }

// Using pass by reference
// use pointer to receive the variable num memory address
func changeNumRef(num *int) {
	// as num is a pointer, in order to assign value we've use dereferencing it i.e, *num
	*num = 5

	fmt.Println("In changeNum", *num)
}

func main() {
	num := 1

	// changeNum(num)

	// Here while using pass by refernce we've to pass the variable num address i.e., &num
	changeNumRef(&num)

	fmt.Println("After changeNumRef in main", num)
}