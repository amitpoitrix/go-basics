package main

import (
	"fmt"
	"maps"
)

// map -> hash, object, dictionary
// key-value pair in Go
func main() {
	// Method 1: In Go we declare map using make()
	m := make(map[string]string)

	// setting element in map
	m["tech"] = "golang"
	m["area"] = "backend"
	// getting element in map
	fmt.Println(m["tech"])
	// IMP: if key doesn't exist in map than it return zeroed value i.e, "" in case of string(value)
	fmt.Println(m["launch"])

	// fmt.Println(m)
	fmt.Println(len(m))


	m1 := make(map[string]int)
	m1["one"] = 1
	m1["two"] = 2
	fmt.Println(m1)

	// removing a key in map using delete()
	fmt.Println("removing a key in map using delete()")
	delete(m1, "one")
	fmt.Println(m1)

	// removing all the key-value pair from map using clear()
	fmt.Println("removing all the key-value pair from map using clear()")
	clear(m1)
	fmt.Println(m1)


	// Mathod 2: Another way of declaring & initializing map
	m2 := map[string]int{"phone": 1, "model": 2023}
	fmt.Println(m2)


	// Now checking whether a particular key exist in map or not
	// here k will contains value of given key in mao & ok will contain bool value whether given 
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

	fmt.Println(maps.Equal(m3, m4))

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
	
	fmt.Println("m6:", m6)
}