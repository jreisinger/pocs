package shiftblock_test

import (
	"bytes"
	"errors"
	"fmt"
	shift "shift-block"
	"testing"
)

func TestNewCipher_GivesNoErrorForValidKey(t *testing.T) {
	t.Parallel()
	_, err := shift.NewCipher(make([]byte, shift.BlockSize))
	if err != nil {
		t.Errorf("want no error, got %v", err)
	}
}

func TestNewCipher_GivesErrKeySizeForInvalidKey(t *testing.T) {
	t.Parallel()
	_, err := shift.NewCipher([]byte{})
	if !errors.Is(err, shift.ErrKeySize) {
		t.Errorf("want %q, got %q", shift.ErrKeySize, err)
	}
}

var testKey = bytes.Repeat([]byte{1}, shift.BlockSize)

var testcases = []struct {
	plaintext, ciphertext []byte
}{
	{
		plaintext:  []byte{0, 1, 2, 3, 4, 5},
		ciphertext: []byte{1, 2, 3, 4, 5, 6},
	},
	{
		plaintext:  []byte("HAL"),
		ciphertext: []byte("IBM"),
	},
}

func TestEncrypt(t *testing.T) {
	t.Parallel()
	cipher, err := shift.NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testcases {
		name := fmt.Sprintf("%x + %x = %x", tc.plaintext, testKey, tc.ciphertext)
		t.Run(name, func(t *testing.T) {
			ciphertext := make([]byte, len(tc.plaintext))
			cipher.Encrypt(ciphertext, tc.plaintext)
			if !bytes.Equal(ciphertext, tc.ciphertext) {
				t.Errorf("want %x, got %x", tc.ciphertext, ciphertext)
			}
		})
	}

}

func TestDecrypt(t *testing.T) {
	t.Parallel()
	cipher, err := shift.NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range testcases {
		name := fmt.Sprintf("%x - %x = %x", tc.ciphertext, testKey, tc.plaintext)
		t.Run(name, func(t *testing.T) {
			plaintext := make([]byte, len(tc.ciphertext))
			cipher.Decrypt(plaintext, tc.ciphertext)
			if !bytes.Equal(plaintext, tc.plaintext) {
				t.Errorf("want %x, go %x", tc.plaintext, plaintext)
			}
		})
	}
}

func TestBlockSize(t *testing.T) {
	t.Parallel()
	cipher, err := shift.NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	want := shift.BlockSize
	got := cipher.BlockSize()
	if want != got {
		t.Errorf("want %v, got %v", want, got)
	}
}
