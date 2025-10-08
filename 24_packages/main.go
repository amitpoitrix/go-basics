package main

import (
	"fmt"

	"github.com/amitpoitrix/go-basics/auth"
	"github.com/amitpoitrix/go-basics/user"
	"github.com/fatih/color"
)

/*
Package
Packages helps to logically group similar codes.
If package is not there than we've to manually write different functions name every time
In Go, we've to write package name starting with small letters and whatever other functions functionality
needs to be used in another pkg we need to write that function name starting with CAPITAL letters else
it'll be private to that package


Installing external package:
search - go package
> go get <external pkg name>
> go get github.com/fatih/color

In order to fix any issue like indirect issue in go.mod or installing pkg which being used
> go mod tidy
*/

/*
Module
In Go, module is just our project name and run below command to create module in go
> go mod init github.com/amitpoitrix/go-basics
Above is the naming convention in order to create module
*/

func main() {
	auth.LoginWithCredentials("john@gmail.com", "password")
	userSession := auth.GetSession()
	fmt.Println(userSession)

	userDetails := user.User{
		Email: "john@gmail.com",
		Name: "John",
	}

	fmt.Println(userDetails.Email, userDetails.Name)

	/* Using external packages */
	color.Green(userDetails.Email)
}
