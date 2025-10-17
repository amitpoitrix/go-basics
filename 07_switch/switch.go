package main

import (
	"fmt"
	"time"
)

func main() {
	// Simple Switch
	// switch is like multiple if-else-if ladder
	// In Go, we don't need to write break statement at the end of each case
	// default is optional

	// i := 3

	// switch i {
	// case 1:
	// 	fmt.Println("One")
	// case 2: 
	// 	fmt.Println("Two")
	// case 3: 
	// 	fmt.Println("Three")
	// default:
	// 	fmt.Println("Other")
	// }


	// Multiple Condition Switch
	switch time.Now().Weekday() {
	case time.Saturday, time.Sunday:
		fmt.Println("Its weekend")
	default:
		fmt.Println("Its a Workday")
	}


	// type switch
	// In Go, function are first class citizens
	// We can pass function as arguments to other functions
	// We can return function from other functions
	// We can assign function to variables
	// We can use empty interface to accept any type of variable
	// interface{} - empty interface - can hold any type of value
	// It is similar to any in TypeScript or Object in JavaScript or object in Python

	// We can use type switch to find the type of variable at runtime
	// i.(type) - is a special construct which can be used only inside type switch
	// It is used to find the type of variable at runtime
	// It is similar to typeof in JavaScript or type() in Python
	// It is used when we want to do different things based on the type of variable
	whoAmI := func(i interface{}) {
		switch i.(type) {
		case int:
			fmt.Println("its a integer")
		case string:
			fmt.Println("its a string")
		case bool:
			fmt.Println("its a boolean")
		default:
			fmt.Println("Others type")
		}
	}

	whoAmI("23")
}