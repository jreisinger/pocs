// Package crc32 detects collisions in CRC32 hash algorithm.
package crc32

import (
	"bytes"
	"hash/crc32"
)

// Collision returns true if msg1 and msg2 are
// different but produce the same CRC32 checksum.
func Collision(msg1, msg2 []byte) bool {
	if bytes.Equal(msg1, msg2) {
		return false
	}

	crc32q := crc32.MakeTable(crc32.IEEE)
	c1 := crc32.Checksum(msg1, crc32q)
	c2 := crc32.Checksum(msg2, crc32q)
	return c1 == c2
}
