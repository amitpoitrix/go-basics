package main

import (
	"fmt"
	"time"
)

func main() {
	normalTime := time.Now()
	fmt.Println("normalTime---", normalTime)

	second := time.Now().Unix()
	fmt.Println("second---", second)

	milli := time.Now().UnixMilli()
	fmt.Println("milli---", milli)

	dateInFormat := time.Now().Format("2006-01-02")
	fmt.Println("dateInFormat:", dateInFormat)
}