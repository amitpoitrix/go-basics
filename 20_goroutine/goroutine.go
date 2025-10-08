package main

import (
	"fmt"
	"sync"
)

/*
Goroutines are light weight threads that helps to execute task concurrently

WaitGroup ensures all goroutines should be completed before moving ahead with next set of main() instructions
*/

/* accept waitgroup address as pointer */
func task(id int, w *sync.WaitGroup) {
	/* using this w.Done() is to remove that added 1 goroutine */
	defer w.Done()	// using defer keyword will execute this line before the exit of the function(even in error)
	
	fmt.Println("doing task", id)
}

func main() {
	/* first creating waitgroup variable using sync pkg WaitGroup */
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		/* Adding 1 means adding 1 goroutine */
		wg.Add(1)

		/* just use go keyword to make task non-blocking which will execute concurrently */
		/* pass waitgroup variable address */
		go task(i, &wg)

		/* we can make use of anonymous self calling function that make use of closure to achieve task(i) */
		// go func (i int) {
		// 	fmt.Println(i)
		// }(i)
	}

	/* Stopping main function to execute concurrent task but not a good practice intead use sync.WaitGroup */
	// time.Sleep(time.Second * 1)

	/* Here we'll wait for completion of Goroutines */
	wg.Wait()
}