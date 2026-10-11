package main

import "fmt"

func init() {
	// Modifying Slices Shares the Backing Array
	// When you slice an existing slice (e.g., s2 := s1[1:3]), no data is copied. 
	// Both slices point to the exact same backing array.
	s1 := []int{10, 20, 30, 40, 50, 60}
	s2 := s1[2:3] // s2 is [20, 30]
	fmt.Println("s2", s2, len(s2), cap(s2)) // [20 30] 2 5 -> now capacity depends on how much element left 
	// from starting point like here starting point is idx = 1 so from 1 there are 4 more elements so total 5 
	// hence 5 capacity, doesn't matter what's the ending idx

	s2[0] = 99    // Modifying s2 changes the backing array!
	fmt.Println(s1) // Output: [10, 99, 30, 40]
	// fmt.Println("s2", s2, len(s2), cap(s2))
	fmt.Println("s1", s1, len(s1), cap(s1)) // [10 20 99 40 50 60] 

	s3 := s1[1:3:3] // full slice expression: cap limited to 2
	// The reason the capacity of c is 2 is because Go uses a three-index slice expression ([low:high:max]) 
	// to explicitly control the resulting slice's length and capacity.
	// The math Go uses to calculate capacity from a three-index slice is:
	// Capacity = max - low

	fmt.Println("s3", s3, len(s3), cap(s3))
}

func main() {
	// 1. Create a slice with length 2 and capacity 2
	s := make([]int, 2, 2)
	s[0] = 100
	s[1] = 200

	fmt.Println("--- INITIAL SLICE ---")
	// %p prints the memory address of the first element in the backing array
	fmt.Printf("Slice s:   Value=%v, Len=%d, Cap=%d, Pointer=%p\n\n", s, len(s), cap(s), s)

	// 2. This append forces a RESIZE because capacity (2) is full
	fmt.Println("--- APPENDING (EXCEEDING CAPACITY) ---")
	s2 := append(s, 300)
	fmt.Printf("Original s: Value=%v, Len=%d, Cap=%d, Pointer=%p (Unchanged)\n", s, len(s), cap(s), s)
	fmt.Printf("New s2:     Value=%v, Len=%d, Cap=%d, Pointer=%p (CHANGED!)\n\n", s2, len(s2), cap(s2), s2)

	// 3. This append stays WITHIN CAPACITY because s2 now has room (capacity is 4)
	fmt.Println("--- APPENDING AGAIN (WITHIN CAPACITY) ---")
	s3 := append(s2, 400)
	fmt.Printf("Original s2: Value=%v, Len=%d, Cap=%d, Pointer=%p\n", s2, len(s2), cap(s2), s2)
	fmt.Printf("New s3:      Value=%v, Len=%d, Cap=%d, Pointer=%p (SAME POINTER!)\n\n", s3, len(s3), cap(s3), s3)

	// 4. Proof of shared memory: mutating s3 affects s2, but NOT s
	fmt.Println("--- PROOF OF SHARED MEMORY ---")
	s3[0] = 999
	fmt.Printf("s (Isolated):  %v\n", s)
	fmt.Printf("s2 (Mutated): %v\n", s2)
	fmt.Printf("s3 (Mutated): %v\n", s3)
}
