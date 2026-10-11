package main

import (
	"fmt"
	"sync"
	// "time"
)

/*
	GMP Model
	G - Goroutine(Task)
	M - Machine (Actual OS Threads)
	P - Processor (Manage goroutine Queue)

	CPU cores = 4
	4 Processor
	Many Goroutines
*/

/*
	Goroutines Lifecycle
	1. Created
	2. Runnable - Waiting in queue
	3. Running - Exceuting on CPU
	4. Waiting/Blocked
	5. Dead/Finish
*/

/*
	Method 1: Goroutines with time
*/

// func Task(task string) {
// 	for i := 1; i <= 4; i++ {
// 		fmt.Println(task, i)
// 		time.Sleep(time.Millisecond * 500)
// 	}
// }

// func main() {
// 	/*
// 		Traditional Threads
// 		Heavy (~1mb)
// 		Managed by OS

// 		goroutine(lightweight threads)
// 		Size(~2kb)
// 		Managed by go runtime
// 	*/

// 	go Task("task 1")
// 	go Task("task 2")

// 	time.Sleep(time.Second * 2)
// }

/*
	Method 2: Using WaitGroup to stop the execution of main() before completion of all goroutines
*/

func Task(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println("id:", id)
}

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		go Task(i, &wg)
	}

	wg.Wait()
}