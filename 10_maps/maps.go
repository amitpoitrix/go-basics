package main

import (
	"fmt"
	"maps"
)

// map -> hash, object, dictionary
// key-value pair in Go
func main() {
	var nm map[string]int // this is Nil map
	_ = nm["x"] // fine, returns 0
	// nm["x"] = 1 // panic: assignment to entry in nil mapset

	// Method 1: Map declaration without make()
	m7 := map[string]int{} // Style A: Map literal
	fmt.Println(m7)
	m7["hello"] = 1
	fmt.Println("m1", m7["hello"])

	// Method 2: In Go we declare map using make()
	m := make(map[string]string) // Style B: make function with optional hint value i.e., capacity
	// m := make(map[string]string, 100) // Initializes an empty map, but hints to Go to pre-allocate space for 100 items

	// setting element in map
	m["tech"] = "golang"
	m["area"] = "backend"
	// getting element in map
	fmt.Println(m["tech"]) // golang
	// IMP: if key doesn't exist in map than it return zeroed value i.e, "" in case of string(value)
	fmt.Println(m["launch"]) // ""

	// fmt.Println(m)
	fmt.Println(len(m)) // 2


	m1 := make(map[string]int)
	m1["one"] = 1
	m1["two"] = 2
	fmt.Println(m1) // map[one:1 two:2]

	// removing a key in map using delete()
	fmt.Println("removing a key in map using delete()")
	delete(m1, "one")
	fmt.Println(m1) // map[two:2]

	// removing all the key-value pair from map using clear()
	fmt.Println("removing all the key-value pair from map using clear()")
	clear(m1)
	fmt.Println(m1) // map[]


	// Mathod 2: Another way of declaring & initializing map
	m2 := map[string]int{"phone": 1, "model": 2023}
	fmt.Println(m2) // map[phone:1 model:2023]


	// Now checking whether a particular key exist in map or not
	// here k will contains value of given key in map & ok will contain bool value whether given 
	// key exists or not

	fmt.Println("Now checking whether a particular key exist in map or not")

	if k, ok := m2["model"]; ok {
		fmt.Println("key exists with value:", k)
	} else {
		fmt.Println("doesn't exist")
	}


	// Checking whether two maps are equal or not using maps pkg using Equal()
	m3 := map[string]int{"company": 1, "capacity": 2000}
	m4 := map[string]int{"company": 1, "capacity": 2001}

	fmt.Println(maps.Equal(m3, m4)) // false

	// Iterating over maps
	fmt.Println("Iterating over map")

	m5 := map[string]int{"amit": 29, "amazon": 30}

	for key, value := range m5 {
		fmt.Println(key, value)
	}


	// Map of struct type
	fmt.Println("Map of struct type")

	type User struct {
		Name string
		Age int
	}

	m6 := make(map[string]User)
	m6["Engineer"] = User{"Amit", 29}
	m6["Tester"] = User{"Rahul", 23}
	
	fmt.Println("m6:", m6) // map[Engineer:{Amit 29} Tester:{Rahul 23}]
}