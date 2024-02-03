package shift_test

import (
	"bytes"
	"fmt"
	"shift-multibytekey"
	"testing"
)

var cases = []struct {
	key        []byte
	plaintext  []byte
	ciphertext []byte
}{
	{
		key:        []byte{1},     // {} -> literal
		plaintext:  []byte("HAL"), // () -> conversion
		ciphertext: []byte("IBM"),
	},
	{
		key:        []byte{1, 2},
		plaintext:  []byte("AAAA"),
		ciphertext: []byte("BCBC"),
	},
	{
		key:        []byte{1, 2, 3},
		plaintext:  []byte{0, 0, 0},
		ciphertext: []byte{1, 2, 3},
	},
	{
		key:        []byte{1, 2},
		plaintext:  []byte{0, 1, 2},
		ciphertext: []byte{1, 3, 3},
	},
}

func TestEncipher(t *testing.T) {
	t.Parallel()
	for _, tc := range cases {
		name := fmt.Sprintf("%s + %d = %s", tc.plaintext, tc.key, tc.ciphertext)
		t.Run(name, func(t *testing.T) {
			got := shift.Encipher(tc.plaintext, tc.key)
			if !bytes.Equal(tc.ciphertext, got) {
				t.Errorf("want %q, got %q", tc.ciphertext, got)
			}
		})
	}
}

func TestDecipher(t *testing.T) {
	t.Parallel()
	for _, tc := range cases {
		name := fmt.Sprintf("%s - %d = %s", tc.ciphertext, tc.key, tc.plaintext)
		t.Run(name, func(t *testing.T) {
			got := shift.Decipher(tc.ciphertext, tc.key)
			if !bytes.Equal(tc.plaintext, got) {
				t.Errorf("want %q, got %q", tc.plaintext, got)
			}
		})
	}
}
