package main

import "fmt"

/*
Pointer in Go
A Pointer is a variable that stores the memory address of another variable

& - address of variable
* - value at address(dereferencing operator)
*/

// Lesson 1: Basics of Pointer
// func main() {
// 	fmt.Println("Lesson 1: Basics of Pointer...")
// 	x := 10

// 	ptr := &x

// 	fmt.Println("Value of x---", x)
// 	fmt.Println("Address of x---", &x)
// 	fmt.Println("Value stored in pointer---", ptr)
// 	fmt.Println("Value at ptr---", *ptr)
// }


// Lesson 2: Pointer to struct - used to changes struct field values
// type User struct {
// 	Name string
// 	Age int
// }

// // Receiving parameter as pointer of user
// func UpdateAge(u *User) {
// 	u.Age = 40
// }

// func main() {
// 	user := User {
// 		Name: "Amit",
// 		Age: 29,
// 	}

// 	// Here passing arguement as address of user
// 	UpdateAge(&user)

// 	fmt.Println("user.Name:", user.Name)
// 	fmt.Println("user.Age:", user.Age)

// 	// Nil Pointer
// 	fmt.Println("Nil Pointer...")
// 	var ptr *int
// 	fmt.Println(ptr)
// }


// Lesson 3: Pointer as Function Parameter
// In Go, parameters are pass by value
func changeNum(num int) {
	fmt.Println("In changeNum", num)
}

// Using pass by reference
// use pointer to receive the variable num memory address
func changeNumRef(num *int) {
	// as num is a pointer, in order to assign value we've to use dereferencing operator it i.e, *num
	*num = *num + 1
	fmt.Println("In changeNumRef", *num)
}

func main() {
	num := 1

	changeNum(num)

	// Here while using pass by refernce we've to pass the variable num address i.e., &num
	changeNumRef(&num)

	fmt.Println("After changeNumRef in main num:", num)
}