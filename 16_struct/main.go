package main

import (
	"encoding/json"
	"fmt"
)

/*
Converting struct to json data type
*/

type User struct {
	ID       int    `json:"id"` // Don't add any space after colon(:)
	Name     string `json:"name"`
	Password string `json:"-"` // As we need to hide the password so I've give dash(-)
}

func main() {
	user := User{
		ID:       1,
		Name:     "Amit",
		Password: "12345",
	}

	// Now coverting user struct to json
	jsonData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(jsonData))
}
