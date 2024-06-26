package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"shift-multibytekey"
)

func main() {
	crib := flag.String("crib", "", "crib text")
	flag.Parse()
	if *crib == "" {
		fmt.Fprintln(os.Stderr, "Please specify a crib text with -crib")
		os.Exit(1)
	}
	ciphertext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	key, err := shift.Crack(ciphertext, []byte(*crib))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	plaintext := shift.Decipher(ciphertext, key)
	os.Stdout.Write(plaintext)
}
