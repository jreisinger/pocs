package main

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: bcrypt <username>\n")
		os.Exit(1)
	}
	username := os.Args[1]

	fmt.Print("Enter password: ")
	var password string
	fmt.Scanln(&password)

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	filename := "passwd"

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s:%s\n", username, hash)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("username and hashed password added to %s\n", filename)
}
