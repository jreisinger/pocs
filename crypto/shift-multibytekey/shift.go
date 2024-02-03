// Package shift implements simple cryptographic functions.
package shift

func Encipher(plaintext []byte, key []byte) []byte {
	ciphertext := make([]byte, len(plaintext))
	for i, b := range plaintext {
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
