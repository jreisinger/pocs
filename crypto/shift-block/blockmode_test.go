package shiftblock_test

import (
	"bytes"
	"shiftblock"
	"testing"
)

func TestEncrypterEnciphersBlockAlignedMessage(t *testing.T) {
	t.Parallel()
	cipher, err := shiftblock.NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	enc := shiftblock.NewEncrypter(cipher)
	plaintext := []byte("This message is exactly 32 bytes")
	want := []byte("Uijt!nfttbhf!jt!fybdumz!43!czuft")
	got := make([]byte, 32)
	enc.CryptBlocks(got, plaintext)
	if !bytes.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestEncrypterCorrectlyReportsCipherBlockSize(t *testing.T) {
	t.Parallel()
	cipher, err := shiftblock.NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	enc := shiftblock.NewEncrypter(cipher)
	want := shiftblock.BlockSize
	got := enc.BlockSize()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}

func TestDecrypterDeciphersBlockAlignedMessage(t *testing.T) {
	t.Parallel()
	cipher, err := shiftblock.NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	dec := shiftblock.NewDecrypter(cipher)
	ciphertext := []byte("Uijt!nfttbhf!jt!fybdumz!43!czuft")
	want := []byte("This message is exactly 32 bytes")
	got := make([]byte, 32)
	dec.CryptBlocks(got, ciphertext)
	if !bytes.Equal(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestDecrypterCorrectlyReportsCipherBlockSize(t *testing.T) {
	t.Parallel()
	cipher, err := shiftblock.NewCipher(testKey)
	if err != nil {
		t.Fatal(err)
	}
	dec := shiftblock.NewDecrypter(cipher)
	want := shiftblock.BlockSize
	got := dec.BlockSize()
	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
}
