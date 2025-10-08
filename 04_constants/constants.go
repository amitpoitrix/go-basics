package main

import "fmt"

// We can use const outside the main() also
const age = 28

// we can't use := outside the main()
// age := 28 // through compilation error

func main() {
	// const language string = "golang"
	const language = "golang"
	// language := "golang" // It is not constant as its var instead of const
	
	fmt.Println(language)

	fmt.Println(age)


	/* Constants Grouping - Multiple constants defining at once */
	const (
		port = 5000
		host = "localhost"
	)

	fmt.Println(port, host)
}
  