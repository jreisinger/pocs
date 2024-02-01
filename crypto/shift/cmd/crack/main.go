package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"shift"
)

func main() {
	crib := flag.String("crib", "", "(first part of) expected plaintext")
	flag.Parse()
	if *crib == "" {
		fmt.Fprintf(os.Stderr, "crack: please use -crib\n")
		os.Exit(1)
	}
	ciphertext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "crack: %v\n", err)
		os.Exit(1)
	}
	key, err := shift.Crack(ciphertext, []byte(*crib))
	if err != nil {
		fmt.Fprintf(os.Stderr, "crack: %v\n", err)
		os.Exit(1)
	}
	plaintext := shift.Decipher(ciphertext, byte(key))
	os.Stdout.Write(plaintext)
}
