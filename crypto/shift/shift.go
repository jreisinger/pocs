// Package shift implements simple cryptographic functions.
package shift

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
