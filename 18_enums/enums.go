package main

import "fmt"

/*
In Go,
we don't have enums so we'll make use of const to achieve the same
*/

/* First create a custom type and use it for first enum field, rest will follow in incremental fashion in case
type is int, in case of string we've assign value manually
*/
// type Status int

// const (
// 	Initial Status = iota // iota is used to assign int value starting from 0
// 	New
// 	Deployed
// 	Received
// 	Delivered
// )

type Status string

const (
	Initial   Status = "initial"
	New       Status = "new"
	Deployed  Status = "deployed"
	Received  Status = "received"
	Delivered Status = "delivered"
)

func changeStatus(status Status) {
	fmt.Println("Status changed to", status)
}

func main() {
	changeStatus(Delivered)
}
