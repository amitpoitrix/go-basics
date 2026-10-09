// package main

// import (
// 	"fmt"
// 	"time"
// )

// func Worker(ch chan int) {
// 	fmt.Println("Worker Started...")

// 	time.Sleep(time.Second * 2)

// 	ch <- 100	// Sender

// 	fmt.Println("Worker value sent")
// }

// func main() {
// 	// Buffered Channel
// 	fmt.Println("Buffered Channel")
// 	ch := make(chan int)

// 	go Worker(ch)

// 	fmt.Println("Waiting for value...")

// 	value := <-ch	// Receiver

// 	fmt.Println("Received:", value)

// 	// Unbuffered Channel
// 	fmt.Println("Unbuffered Channel")

// 	ch1 := make(chan int, 3)

// 	ch1 <- 13
// 	ch1 <- 6
// 	ch1 <- 9

// 	fmt.Println(<-ch1)
// 	fmt.Println(<-ch1)
// 	fmt.Println(<-ch1)
// }

package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// first goroutine(using Anonymous Function)
	go func() {
		time.Sleep(time.Second * 1)
		ch1 <- "Channel 1 completed"
	}()

	// second goroutine
	go func() {
		time.Sleep(time.Second * 1)
		ch2 <- "Channel 2 completed"
	}()
	
	// Whichever channel is ready first that will be executed
	select {
	case msg1 := <-ch1:
		fmt.Println(msg1)
	case msg2 := <-ch2:
		fmt.Println(msg2)
	}
}