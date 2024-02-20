package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	shift "shift-block"
)

func main() {
	keyHex := flag.String("key", "", "32-byte key in hexadecimal")
	debug := flag.Bool("debug", false, "print raw and padded input bytes in hex")
	flag.Parse()
	key, err := hex.DecodeString(*keyHex)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encipher: %v\n", err)
		os.Exit(1)
	}
	cipher, err := shift.NewCipher(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encipher: %v\n", err)
		os.Exit(1)
	}
	enc := shift.NewEncrypter(cipher)
	plaintext, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "encipher: %v\n", err)
		os.Exit(1)
	}
	if *debug {
		fmt.Printf("raw:\t% x\n", plaintext)
	}
	plaintext = shift.Pad(plaintext, enc.BlockSize())
	if *debug {
		fmt.Printf("padded:\t% x\n", plaintext)
	}
	ciphertext := make([]byte, len(plaintext))
	enc.CryptBlocks(ciphertext, plaintext)
	os.Stdout.Write(ciphertext)
}
