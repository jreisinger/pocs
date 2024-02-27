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
	plaintext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encipher: %v\n", err)
		os.Exit(1)
	}
	ciphertext := shift.Encipher(plaintext, byte(*key))
	os.Stdout.Write(ciphertext)
}
