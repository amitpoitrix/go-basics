package main

import "fmt"

// for -> only construct in go for looping
func main() {
	// for can be used in 3 ways
	// 1. classic for loop
	// for i := 0; i <= 3; i++ {
	// 	// break
	// 	if i == 2 {
	// 		continue
	// 	}
	// 	fmt.Println(i)
	// }


	// 2. while loop - there is no while keyword in go, we implement it using for loop
	// i := 1
	// for i <= 3 {
	// 	fmt.Println(i)
	// 	i++
	// }


	// 3. Infinite loop
	// for {
	// 	fmt.Println("1")
	// }


	// 4. Go 1.22 new feature - range
	for i := range 3 {
		fmt.Println(i)
	} 
}	