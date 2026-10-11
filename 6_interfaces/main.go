package main

import (
	"fmt"
	"time"
)

/* Interface In Go - Another Example */

type Logger interface {
	log(message string)
}


type ConsoleLogger struct {}

func (c ConsoleLogger) log(message string) {
	fmt.Println("Console logger:", message)
}

type FileLogger struct {}

func (f FileLogger) log(message string) {
	fmt.Println("File Logger:", message)
}

type TimeLogger struct {}

func (t TimeLogger) log(message string) {
	time := time.Now().Format("2006-01-02")
	fmt.Println(time, " - ", message)
}


func ProcessLogger(logger Logger) {
	logger.log("Task Started")
	logger.log("Process")
	logger.log("Task Completed")
}

func main() {
	consoleLogger := ConsoleLogger{}
	fileLogger := FileLogger{}
	timeLogger := TimeLogger{}

	fmt.Println("Console Logger---")
	ProcessLogger(consoleLogger)

	fmt.Println("File Logger---")
	ProcessLogger(fileLogger)

	fmt.Println("Time Logger---")
	ProcessLogger(timeLogger)
}