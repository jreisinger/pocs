Simple shift (Caesar) cipher. Adapted from https://github.com/bitfield/eg-crypto/.

```
echo hello world | go run ./cmd/encipher -key 10 | go run ./cmd/decipher -key 10
echo hello world | go run ./cmd/encipher -key 10 | go run ./cmd/crack -crib hell

go run ./cmd/encipher -key 253 < ./testdata/tiger.txt > ./testdata/tiger.bin
go run ./cmd/crack -crib 'The tiger' < ./testdata/tiger.bin

go run ./cmd/encipher -key 99 < ./testdata/devil.png > ./testdata/devil.bin
go run ./cmd/crack --crib $(printf '\x89PNG') < ./testdata/devil.bin > ./testdata/devil.png
```