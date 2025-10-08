package auth

import "fmt"

/* Writing function name starting with CAPITAL letter to make it public scope */
func LoginWithCredentials(username string, password string) {
	fmt.Println("Login with username:", username, "and password:", password)
}