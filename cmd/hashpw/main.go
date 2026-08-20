// Command hashpw prints a bcrypt hash for a password, for seeding an admin
// user row directly in the users table (password column).
//
// Usage: go run ./cmd/hashpw "yourpassword"
package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/hashpw <password>")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(os.Args[1]), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error hashing password:", err)
		os.Exit(1)
	}

	fmt.Println(string(hash))
}
