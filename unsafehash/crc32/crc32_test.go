package crc32

import (
	"fmt"
	"testing"
)

func ExampleCollision() {
	if Collision([]byte("gnu"), []byte("codding")) {
		fmt.Println("CRC32 collision!")
	}
	// Output: CRC32 collision!
}

func TestCollision(t *testing.T) {
	tests := []struct {
		msg1, msg2 []byte
		want       bool
	}{
		{[]byte(""), []byte(""), false},
		{[]byte("0"), []byte("1"), false},
	}

	for i, test := range tests {
		if got := Collision(test.msg1, test.msg2); got != test.want {
			t.Errorf("test %d: collision %t, want %t", i, got, test.want)
		}

	}
}
