package main

import (
	"fmt"
	"time"
)

/*
In Go,
goroutines uses channels to communicate with each other so we send data to channel which is being 
received as data to another goroutines
*/

/* Receiving Channel */
func processNum(numChan <-chan int) {
	/* As here we're receiving data from channel so using <- as suffix before channel name */
	// fmt.Println("processing number", <-numChan)

	/* When there are more than one data in channel */
	for data := range numChan {
		// As we're looping over channel so we can directly access data
		fmt.Println("processing number", data)
		time.Sleep(time.Second)
	}
}

/* Sending channel */
func sum(result chan<- int, num1 int, num2 int) {
	totalSum := num1 + num2

	/* sending result to channel */
	result <- totalSum
}

/* using unbuffered channel in order make main function wait for go routines completion i.e., synchronization */
func process(isDone chan<- bool) {
	/* below code will always execute even this function throw error because of defer keyword */
	defer func() { isDone <- true }()

	fmt.Println("Processing...")
}

/* 
Sending emails which is being received over email channel i.e., emailChan 
As we're receiving 2 channels so in order to make each channel strictly adhere to receiving channel &
sending channel we can use below symbol
<-chan : this is receiving channel means we can only fetch data from this channel
chan<- : this is sending channel means we can only send data to this channel
*/
func sendingEmails(emailChan <-chan string, isDone chan<- bool) {
	defer func() { isDone <- false }()

	for email := range emailChan {
		fmt.Println("sending email to", email)
		time.Sleep(time.Second)
	}
}

func main() {
	/* creating unbuffered channel using make() */
	// messageChan := make(chan int)

	/* sending msg to messageChan channel */
	// messageChan <- "hello channels" // blocking 

	/* receiving msg from messageChn channel */
	// receivedMsg := <-messageChan // blocking
	// fmt.Println(receivedMsg)

	/*
		Because of Unbuffered Channels both sending & receiving data is blocking operation, it'll lead to
		deadlock situation
		Now in order to avoid this we'll send this channel to another goroutine i.e., another function and
		make calling function goroutine
	*/
	// go processNum(messageChan)

	/* And we've send data to channel */
	// messageChan <- "hello channels" // sending data to messageChan (non-blocking due to go routines)

	/* Stopping main function for 1 sec to allow completition of goroutines */
	// time.Sleep(time.Second)

	/* Now we'll send unlimited no. of data in the channel */
	// for {
	// 	messageChan <- rand.Intn(100)
	// }


	/* Now receiving data from another go routine to current routine */
	// result := make(chan int)

	/* passing channel as "sending channel" to go routine */
	// go sum(result, 2, 5)

	// /* receiving data from result channel */
	// response := <-result // blocking i.e., receiving channel from result channel
	// fmt.Println(response)

	/*
		In order to make go routine synchronize we've make main() wait for its completion for which we make
		use of time.Sleep() but its not a good practice
		Another way is to use sync.WaitGroup()
		OR
		We can make use of unbuffered channel to achieve the same process
	*/
	// isDone := make(chan bool)
	// go process(isDone)
	// <-isDone // blocking


	/*
		Buffered Channel - sending & receiving data from these channels are non-blocking as channel size is
		is already defined
		And also because of this non-blocking nature there won't be any Deadlock situation
	*/
	// emailChan := make(chan string, 100)
	// isDone := make(chan bool)

	// emailChan <- "1@gmail.com"
	// emailChan <- "2@gmail.com"

	// fmt.Println(<-emailChan)
	// fmt.Println(<-emailChan)

	// go sendingEmails(emailChan, isDone)

	// /* Now sending more than one email */
	// for i := range 5 {
	// 	emailChan <- fmt.Sprintf("%d@gmail.com", i)
	// }

	// fmt.Println("Done sending emails to channel")

	// /* IMP: In case of buffered channel we've to close it else it'll through error */
	// close(emailChan)

	// <-isDone // blocking

	/*
		Sending data to multiple go routines over multiple channels
	*/
	chan1 := make(chan int)
	chan2 := make(chan string)

	go func() {
		chan1 <- 25
	}()

	go func() {
		chan2 <- "pong"
	}()

	for range 2 {
		select {
		case chan1Val := <-chan1:
			fmt.Println("Received data from chan1", chan1Val)
		case chan2Val := <-chan2:
			fmt.Println("Received data from chan2", chan2Val)
		}
	}
}
