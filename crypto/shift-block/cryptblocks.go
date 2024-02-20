package shiftblock

import "crypto/cipher"

// Encrypter implements crypo/cipher.BlockMode. Encrypter extends the capability
// of Cipher to stream of blocks.
type Encrypter struct {
	cipher    cipher.Block
	blockSize int
}

func NewEncrypter(c cipher.Block) Encrypter {
	return Encrypter{
		cipher:    c,
		blockSize: c.BlockSize(),
	}
}

func (e Encrypter) CryptBlocks(dst, src []byte) {
	if len(src)%e.blockSize != 0 {
		panic("encrypter: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("encrypter: output smaller than input")
	}
	// Keep chopping block-sized pieces off the plaintext
	// and enciphering them until there are no more pieces.
	for len(src) > 0 {
		e.cipher.Encrypt(dst[:e.blockSize], src[:e.blockSize])
		dst = dst[e.blockSize:]
		src = src[e.blockSize:]
	}

}

func (e Encrypter) BlockSize() int {
	return e.blockSize
}

type Decrypter struct {
	cipher    cipher.Block
	blockSize int
}

func NewDecrypter(c cipher.Block) Decrypter {
	return Decrypter{
		cipher:    c,
		blockSize: c.BlockSize(),
	}
}

func (d Decrypter) CryptBlocks(dst, src []byte) {
	if len(src)%d.blockSize != 0 {
		panic("decrypter: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("decrypter: output smaller than input")
	}
	for len(src) > 0 {
		d.cipher.Decrypt(dst[:d.blockSize], src[:d.blockSize])
		dst = dst[d.blockSize:]
		src = src[d.blockSize:]
	}
}

func (d Decrypter) BlockSize() int {
	return d.blockSize
}
