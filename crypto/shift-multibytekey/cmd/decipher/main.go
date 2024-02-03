package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"shift-multibytekey"
)

func main() {
	keyHex := flag.String("key", "01", "key in hexadecimal")
	flag.Parse()

	key, err := hex.DecodeString(*keyHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		os.Exit(1)
	}

	ciphertext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	plaintext := shift.Decipher(ciphertext, key)
	os.Stdout.Write(plaintext)
}
