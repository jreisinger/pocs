Simple shift (Caesar) cipher that supports multibyte encryption keys (`shift` supports just single byte keys). Taken from https://github.com/bitfield/eg-crypto/.

```
$ echo hello | go run ./cmd/encipher -key DEADBEEF
F*[M�
```

Unlike with the single-byte version, the same plaintext letter does not always produce the same ciphertext letter. The "ll" is enciphered as "[M". This makes the frequency analysis a lot harder for Eve.

The longer key is harder to bruteforce but it's still possible:

```
go run ./cmd/encipher -key DEADBEEF < ../shift/testdata/tiger.txt | go run ./cmd/crack -crib 'The tiger'
```

NOTE: the crib must be at least as long as the key; it this case 4 bytes, i.e. 'The '.
