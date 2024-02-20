package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"shiftblock"
)

func main() {
	keyHex := flag.String("key", "", "32-byte key in hexadecimal")
	flag.Parse()
	key, err := hex.DecodeString(*keyHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decipher: %v\n", err)
		os.Exit(1)
	}
	cipher, err := shiftblock.NewCipher(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decipher: %v\n", err)
		os.Exit(1)
	}
	ciphertext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "decipher: %v\n", err)
		os.Exit(1)
	}
	plaintext := make([]byte, len(ciphertext))
	dec := shiftblock.NewDecrypter(cipher)
	dec.CryptBlocks(plaintext, ciphertext)
	plaintext = shiftblock.Unpad(plaintext, shiftblock.BlockSize)
	os.Stdout.Write(plaintext)
}
