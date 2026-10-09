package main

import "fmt"

func main() {
	// In Go
	// we don't use brackets () while defining condition for if statement
	// Go doesn't have Ternary Operator we've to make use of if else
	
	// age := 10

	// if age >= 18 {
	// 	fmt.Println("Person is an Adult")
	// } else if age >= 12 {
	// 	fmt.Println("Person is a Teenager")
	// } else {
	// 	fmt.Println("Person is a Child")
	// }


	// Conditional Operator
	// age := 18
	// isEligible := false

	// if age >= 18 && isEligible {
	// 	fmt.Println("Yes")
	// } else {
	// 	fmt.Println("No")
	// }


	// we can define variable at the same time while defining condition for if statement
	if age := 14; age >= 18 {
		fmt.Println("Person is an Adult", age)
	} else if age >= 12 {
		fmt.Println("Person is a Teenager", age)
	} else {
		fmt.Println("Person is a Child", age)
	}
}