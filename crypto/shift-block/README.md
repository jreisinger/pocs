Package shiftblock is a block version of the shift (Caesar) cipher. Block ciphers work on chunks of data (stream ciphers operate on data one bit at a time). The chunk size is called the block size.

In order to work with messages that are not block-aligned, shiftblock always uses a padding. Both the number and the value of padded bytes is equal to the difference from the nearest multiple of block size. If the message size is aligned with the block size, the number and the value of padded bytes is equal to the block size. In our case the block size is 32 bytes.

```
$ export KEY=0101010101010101010101010101010101010101010101010101010101010101

# Padding of 32 bytes with bytes value 32 (20 in hex).
$ echo "This is 32 bytes, including EOF" | go run ./cmd/encipher -key $KEY -debug
raw:    54 68 69 73 20 69 73 20 33 32 20 62 79 74 65 73 2c 20 69 6e 63 6c 75 64 69 6e 67 20 45 4f 46 0a
padded: 54 68 69 73 20 69 73 20 33 32 20 62 79 74 65 73 2c 20 69 6e 63 6c 75 64 69 6e 67 20 45 4f 46 0a 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20 20
Uijt!jt!43!czuft-!jodmvejoh!FPG
                               !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

# Padding of 1 byte with byte value 1 (01 in hex).
$ echo "This is 31 long, including EOF" | go run ./cmd/encipher -key $KEY -debug
raw:    54 68 69 73 20 69 73 20 33 31 20 6c 6f 6e 67 2c 20 69 6e 63 6c 75 64 69 6e 67 20 45 4f 46 0a
padded: 54 68 69 73 20 69 73 20 33 31 20 6c 6f 6e 67 2c 20 69 6e 63 6c 75 64 69 6e 67 20 45 4f 46 0a 01
Uijt!jt!42!mpoh-!jodmvejoh!FPG
```

```
$ go run ./cmd/encipher/ -key $KEY < ../shift/testdata/tiger.txt | go run ./cmd/decipher/ -key $KEY 
The tiger appears at its own pleasure. When we become very silent at that
place, with no expectation of the tiger, that is when he chooses to appear...
When we stand at the edge of the river waiting for the tiger, it seems that the
silence takes on a quality of its own. The mind comes to a stop. In the Indian
tradition that is the moment when the teacher says, “You are that. You are that
silence. You are that.”
--Francis Lucille, “The Perfume of Silence”
```