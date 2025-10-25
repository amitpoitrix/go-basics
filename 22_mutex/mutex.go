package main

import (
	"fmt"
	"sync"
)

/*
In Go, mutex stands for Mutually Exclusive that helps in preventing race condition
*/

type post struct {
	views int
	mu sync.Mutex /* Defining mutex at struct */
}

/* 
Now this inc() is a common resource that is being accesed by multiple go routines leading to race condition 
Now in order to avoid this we gonna use mutex(in struct from sync.Mutex) that'll preveent this race condition
*/
func (p *post) inc(wg *sync.WaitGroup) {
	defer func(){
		p.mu.Unlock() // unlocking here which will run even after throwing error
		wg.Done()	// this will reduce the go routine count by 1
	}()
	
	/* Always apply lock on modified piece of code */
	p.mu.Lock()
	p.views++
}

func main() {
	myPost := post{
		views: 0,
	}

	var wg sync.WaitGroup

	for range 100 {
		wg.Add(1)
		go myPost.inc(&wg)
	}

	wg.Wait()
	fmt.Println(myPost.views)
}