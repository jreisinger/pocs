package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"shift"
)

func main() {
	key := flag.Int("key", 1, "shift value")
	flag.Parse()
	ciphertext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decipher: %v\n", err)
		os.Exit(1)
	}
	plaintext := shift.Decipher(ciphertext, byte(*key))
	os.Stdout.Write(plaintext)
}
