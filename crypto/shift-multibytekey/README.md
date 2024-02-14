Simple shift (Caesar) cipher that supports multibyte encryption keys (`shift` supports just single byte keys). Taken from https://github.com/bitfield/eg-crypto/.

```
echo hello | go run ./cmd/encipher -key DEADBEEF | go run ./cmd/decipher -key DEADBEEF
```