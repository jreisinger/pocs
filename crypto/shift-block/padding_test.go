package shiftblock_test

import (
	"bytes"
	"shiftblock"
	"testing"
)

var padCases = []struct {
	name        string
	raw, padded []byte
}{
	{
		name:   "1 short of full block",
		raw:    []byte{0, 0, 0},
		padded: []byte{0, 0, 0, 1},
	},
	{
		name:   "2 short of full block",
		raw:    []byte{0, 0},
		padded: []byte{0, 0, 2, 2},
	},
	{
		name:   "full block",
		raw:    []byte{0, 0, 0, 0},
		padded: []byte{0, 0, 0, 0, 4, 4, 4, 4},
	},
	{
		name:   "empty block",
		raw:    []byte{},
		padded: []byte{4, 4, 4, 4},
	},
}

func TestPad(t *testing.T) {
	t.Parallel()
	blockSize := 4
	for _, tc := range padCases {
		t.Run(tc.name, func(t *testing.T) {
			got := shiftblock.Pad(tc.raw, blockSize)
			if !bytes.Equal(tc.padded, got) {
				t.Errorf("want %v, got %v", tc.padded, got)
			}
		})
	}
}

func TestUnpad(t *testing.T) {
	t.Parallel()
	blockSize := 4
	for _, tc := range padCases {
		got := shiftblock.Unpad(tc.padded, blockSize)
		if !bytes.Equal(tc.raw, got) {
			t.Errorf("want %v, got %v", tc.raw, got)
		}
	}
}

func TestPadWithManyBytes(t *testing.T) {
	t.Parallel()
	blockSize := 32
	raw := bytes.Repeat([]byte{0}, 52)
	padded := append(raw, bytes.Repeat([]byte{12}, 12)...)
	got := shiftblock.Pad(raw, blockSize)
	t.Logf("got: % 0x\n", got)
	if !bytes.Equal(padded, got) {
		t.Errorf("want %v, got %v", padded, got)
	}
}
