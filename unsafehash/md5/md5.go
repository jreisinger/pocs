// Package md5 detects collisions in MD5 hash algorithm.
package md5

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"os"
)

// Collision returns true if msg1 and msg2 are
// different but produce the same MD5 hash value.
func Collision(msg1, msg2 []byte) bool {
	if bytes.Equal(msg1, msg2) {
		return false
	}

	h1 := md5.New()
	h2 := md5.New()
	if _, err := h1.Write(msg1); err != nil {
		fmt.Fprintf(os.Stderr, "writing hash: %v", err)
		return false
	}
	if _, err := h2.Write(msg2); err != nil {
		fmt.Fprintf(os.Stderr, "writing hash: %v", err)
		return false
	}
	return bytes.Equal(h1.Sum(nil), h2.Sum(nil))
}
