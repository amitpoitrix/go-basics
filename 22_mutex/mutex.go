package main

import (
	"fmt"
	"sync"
)

/*
In Go, mutex stands for Mutually Exclusive that helps in preventing race condition
*/

/* 
	Implementation 1
	Now this inc() is a common resource that is being accesed by multiple goroutines leading to race 
	condition, Now in order to avoid this we gonna use mutex(in struct from sync.Mutex) that'll 
	prevent this race condition
*/

// type post struct {
// 	views int
// 	mu sync.Mutex /* Defining mutex at struct */
// }

// func (p *post) inc(wg *sync.WaitGroup) {
// 	defer func(){
// 		p.mu.Unlock() // unlocking here which will run even after throwing error
// 		wg.Done()	// this will reduce the go routine count by 1
// 	}()
	
// 	/* Always apply lock on modified piece of code */
// 	p.mu.Lock()
// 	p.views++
// }

// func main() {
// 	myPost := post{
// 		views: 0,
// 	}

// 	var wg sync.WaitGroup

// 	for range 100 {
// 		wg.Add(1)
// 		go myPost.inc(&wg)
// 	}

// 	wg.Wait()
// 	fmt.Println(myPost.views)
// }


/*
	Implementation 2: Without struct
*/

var counter int

var mu sync.Mutex

func increment(wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 1; i <= 10000; i++ {
		mu.Lock()
		counter++
		mu.Unlock()
	}
}

func main() {
	var wg sync.WaitGroup

	wg.Add(8)

	go increment(&wg)
	go increment(&wg)
	go increment(&wg)
	go increment(&wg)
	go increment(&wg)
	go increment(&wg)
	go increment(&wg)
	go increment(&wg)

	wg.Wait()

	fmt.Println("Counter:", counter)
}