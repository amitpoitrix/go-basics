package main

import "fmt"

/*
panic is just a exception that crash the code when occur
recover helps the code to recover from panic and avoid code crash and should always invoke
inside defer function

*/

func main() {
	fmt.Println("Welcome to  Panic and Recover")

	fmt.Println("Calling divide function")
	result := divide(10, 2)
	fmt.Println("result:",result)
	fmt.Println("End of Divide() call")
}

func divide(a, b int) (result int) {
	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("Recover from panic")
			result = 0	// fallback value
		}
	}()	// defer function should be self invoked function

	if(b == 0) {
		panic("Dividing by zero")
	}

	return a / b
}