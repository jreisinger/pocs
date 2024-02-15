Simple shift (Caesar) cipher that supports multibyte encryption keys (`shift` supports just single byte keys). Taken from https://github.com/bitfield/eg-crypto/.

```
echo hello | go run ./cmd/encipher -key DEADBEEF | go run ./cmd/decipher -key DEADBEEF
```

The key, in some sense, is just a single number, no matter how many bytes it takes to express it. For example, if we had a 32-byte (that is, 256-bit) key, we could express it as either a series of 32 integers (one for each byte), or as a single very large integer. But Go's `int64` can hold only 8 bytes (or 64 bits) worth of information... There's a neat and concise way to write large integers: as a string, using hexadecimal notation. For example the decimal number 3 735 928 559 can be represented as DEADBEEF (4 bytes) in hex, isn't that funny? :-) If fact, any given byte can be written as exactly two hex digits, which is convenient.

Unlike with the single-byte version, the same plaintext letter does not always produce the same ciphertext letter:

```
$ echo hello | go run ./cmd/encipher -key DEADBEEF
F*[M�
```

The "ll" is enciphered as "[M". This makes the frequency analysis a lot harder for Eve.

The longer key is harder to bruteforce but it's still possible:

```
go run ./cmd/encipher -key DEADBEEF < ../shift/testdata/tiger.txt | go run ./cmd/crack -crib 'The tiger'
```

NOTE: the crib must be at least as long as the key; it this case 4 bytes, i.e. 'The '.
