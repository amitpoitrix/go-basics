package main

import (
	"fmt"
	"maps"
)

// map -> hash, object, dictionary
// key-value pair in Go
func main() {
	// In Go we declare map using make()
	m := make(map[string]string)

	// setting element in map
	m["tech"] = "golang"
	m["area"] = "backend"
	// getting element in map
	fmt.Println("m['tech']:", m["tech"])
	// IMP: if key doesn't exist in map than it return zeroed value i.e, "" in case of string(value)
	fmt.Println("m['launch']:", m["launch"])

	// fmt.Println(m)
	fmt.Println("len(m):", len(m))


	m1 := make(map[string]int)
	m1["one"] = 1
	m1["two"] = 2
	fmt.Println("m1:", m1)
	// removing a key in map using delete()
	delete(m1, "one")
	fmt.Println("after deleting key 'one' in m1:", m1)
	// removing all the key-value pair from map using clear()
	clear(m1)
	fmt.Println("after clear in m1:", m1)


	// Another way of declaring & initializing map
	m2 := map[string]int{"phone": 1, "model": 2023}
	fmt.Println("m2:", m2)

	// Now checking whether a particular key exist in map or not
	// here v will contains value of given key in map & ok will contain bool value whether given 
	// key exists or not
	v, ok := m2["models"] 
	fmt.Println("v:", v)

	if ok {
		fmt.Println("key exists")
	} else {
		fmt.Println("doesn't exist")
	}


	// Checking whether two maps are equal or not using maps pkg using Equal()
	m3 := map[string]int{"company": 1, "capacity": 2000}
	m4 := map[string]int{"company": 1, "capacity": 2000}

	fmt.Println("maps.Equal(m3, m4):", maps.Equal(m3, m4))
}