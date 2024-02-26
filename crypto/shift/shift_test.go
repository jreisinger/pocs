package shift_test

import (
	"bytes"
	"shift"
	"testing"
)

var cases = []struct {
	key                   byte
	plaintext, ciphertext []byte
}{
	{
		plaintext:  []byte("HAL"),
		ciphertext: []byte("IBM"),
		key:        1,
	},
}

func TestEncrypt(t *testing.T) {
	for _, tc := range cases {
		got := shift.Encipher(tc.plaintext, tc.key)
		if !bytes.Equal(tc.ciphertext, got) {
			t.Errorf("want %v, got %v", tc.ciphertext, got)
		}
	}
}

func TestDecrypt(t *testing.T) {
	for _, tc := range cases {
		got := shift.Decipher(tc.ciphertext, tc.key)
		if !bytes.Equal(tc.plaintext, got) {
			t.Errorf("want %v, got %v", tc.plaintext, got)
		}
	}
}

func TestCrack(t *testing.T) {
	t.Parallel()
	for _, tc := range cases {
		crib := tc.plaintext[:3]
		got, err := shift.Crack(tc.ciphertext, crib)
		if err != nil {
			t.Fatal(err)
		}
		if tc.key != got {
			t.Errorf("want %q, got %q", tc.key, got)
		}
	}
}

func TestCrackReturnsErrorWhenKeyNotFound(t *testing.T) {
	t.Parallel()
	_, err := shift.Crack([]byte("no good"), []byte("bogus"))
	if err == nil {
		t.Errorf("want error when key not found, got nil")
	}
}
