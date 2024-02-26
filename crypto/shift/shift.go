// Package shift implements simple cryptographic functions.
package shift

import (
	"bytes"
	"errors"
)

func Encipher(plaintext []byte, key byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	for i, b := range plaintext {
		ciphertext[i] = b + key
	}
	return ciphertext
}

func Decipher(ciphertext []byte, key byte) []byte {
	return Encipher(ciphertext, -key)
}

// Crack tries to guess the key used to produce the ciphertext. It needs a crib,
// i.e. a few bytes from the beginning of the plaintext.
func Crack(ciphertext, crib []byte) (key byte, err error) {
	for guess := 0; guess < 256; guess++ {
		result := Decipher(ciphertext[:len(crib)], byte(guess))
		if bytes.Equal(result, crib) {
			return byte(guess), nil
		}
	}
	return 0, errors.New("no key found")
}
