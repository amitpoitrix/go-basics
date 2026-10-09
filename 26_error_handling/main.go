// Error Handling in Go

/*
error interface

type Error interface {
	Error() string
}
*/

package main

import (
	"errors"
	"fmt"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("Cannot divide by zero")
	}

	return a/b, nil
}

func findUser(id int) (error) {
	return fmt.Errorf("User with id %d not found", id)
}

func main() {
	// result, err := divide(10, 2)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }

	// fmt.Println("Result:", result)

	err := findUser(2)
	fmt.Println(err)
}