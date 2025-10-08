package main

import "fmt"

func main() {
	/* string */
	// var name string = "Amit"
	// var name = "Amit" // this will infer the string type

	// Shorthand syntax - when we want to declare and assign variable at the same time
	name := "Amit" 

	fmt.Println(name)


	/* bool */
	// var isAdult bool = true
	isAdult := false
	fmt.Println(isAdult)


	/* int */
	// var num int = 25
	num := 40
	fmt.Println(num)


	/* float */
	// var no float32 = 34.5
	no := 50.9
	fmt.Println(no)


	/* Reason for using type - just Declaration of variable but assigning any values */
	var age int
	age = 28
	fmt.Println(age)
}