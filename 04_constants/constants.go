package main

import "fmt"

// We can use const outside the main() also
const age = 28

// we can't use := outside the main()
// age := 28 // through compilation error

func main() {
	// Constants - Immutable values which cannot be changed
	// const - keyword to declare constant
	// const pi = 3.14
	// const pi float32 = 3.14
	// const language string = "golang"

	// var is used to declare variable which can be changed
	// language := "golang" // It is not constant as its var instead of const
	
	const language = "golang"
	
	fmt.Println(language)

	fmt.Println(age)


	/* Constants Grouping - Multiple constants defining at once */
	const (
		port = 5000
		host = "localhost"
	)

	fmt.Println(port, host)
}
  