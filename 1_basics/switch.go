package main

import (
	"fmt"
	"time"
)

func main() {
	// Simple Switch
	// In Go, inside each case Go will take care of break statement
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
	// Here, interface{} is a any type
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