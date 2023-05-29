package main

import (
	"fmt"
	"log"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("UID:", user.Uid)
	fmt.Println("GID:", user.Gid)
	fmt.Println("Username:", user.Username)
	fmt.Println("Name:", user.Name)
	fmt.Println("HomeDir:", user.HomeDir)
}
