package shiftblock

import (
	"crypto/cipher"
	"errors"
	"fmt"
)

const BlockSize = 32 // bytes

var ErrKeySize = errors.New("invalid key size")

// Cipher implements crypo/cipher.Block interface.
type Cipher struct {
	key [BlockSize]byte
}

func NewCipher(key []byte) (cipher.Block, error) {
	if len(key) != BlockSize {
		return nil, fmt.Errorf("%w %d (must be %d)", ErrKeySize, len(key), BlockSize)
	}
	return &Cipher{
		key: [BlockSize]byte(key),
	}, nil
}

func (c *Cipher) Encrypt(dst, src []byte) {
	for i, b := range src {
		dst[i] = b + c.key[i]
	}
}

func (c *Cipher) Decrypt(dst, src []byte) {
	for i, b := range src {
		dst[i] = b - c.key[i]
	}
}

func (c *Cipher) BlockSize() int {
	return BlockSize
}
