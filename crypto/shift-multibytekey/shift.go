// Package shift implements simple cryptographic functions.
package shift

import (
	"bytes"
	"errors"
)

func Encipher(plaintext []byte, key []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	for i, b := range plaintext {
		// A mod B is the remainder that's left after dividing A by B as many
		// times as you can. E.g. 5 mod 2 = 1. Modular arithmetic is sometimes
		// called "clock arithmentic" because it wraps around like an analog
		// clock; 12 hours later than 5 o'clock can't be 17 o'clock, it's 5
		// o'clock again. To put it another way, 17 mod 12 = 5.
		ciphertext[i] = b + key[i%len(key)]
	}
	return ciphertext
}

func Decipher(ciphertext []byte, key []byte) []byte {
	plaintext := make([]byte, len(ciphertext))
	for i, b := range ciphertext {
		plaintext[i] = b - key[i%len(key)]
	}
	return plaintext
}

const MaxKeyLen = 32 // bytes

func Crack(ciphertext, crib []byte) (key []byte, err error) {
	for k := range min(MaxKeyLen, len(ciphertext)) {
		for guess := range 256 {
			result := ciphertext[k] - byte(guess)
			if result == crib[k] {
				key = append(key, byte(guess))
				break
			}
		}
		if bytes.Equal(crib, Decipher(ciphertext[:len(crib)], key)) {
			return key, nil
		}
	}
	return nil, errors.New("no key found")
}
