// Package unsafehash detects collisions in hashing algorithms. Algorithms in
// this package are either cryptographically broken (MD5) or don't posses
// sufficient collision resistance (CRC32). Either way they are not suitable for
// security purposes.
package unsafehash
