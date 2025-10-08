package main

import "fmt"

// range is used to iterates over data structure
func main() {
	// Iterating over slice using classic for loop
	nums := []int{2, 6, 4}
	sum := 0

	for i := 0; i < len(nums); i++ {
		fmt.Println(nums[i])
		sum += nums[i]
	}
	fmt.Println(sum)


	// Now using range to do the same thing as above
	fmt.Println("Using Range")
	
	// below range returns 2 values first is index & second is value
	// Note: If you don't want to use first variable value i.e., index than give variable name as "_" else
	// Go will through compilation error
	for _, num := range nums {
		fmt.Println(num)
	}


	// Iterating over map using range
	fmt.Println()
	fmt.Println("Iterating over map using for range")

	m := map[string]string {"company":"abc", "branch": "engineering"}

	// Here, k is key & v is value
	// if only one variable is declare than its just a key for the map
	for k, v := range m {
		fmt.Println(k, v)
	}


	// Iterating over string using for range
	fmt.Println()
	fmt.Println("Iterating over string using for range")

	str := "golang"
	
	// unicode code point rune
	// starting byte of rune
	// 300 -> uses 2 bytes and character till value 0 to 255 uses 1 byte
	// So here i is not index but starting byte of each and every character
	for i, c := range str {
		fmt.Println(i, c, string(c))
	}
}