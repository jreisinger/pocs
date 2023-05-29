package md5

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func ExampleCollision() {
	// Slighly different strings that hash to the same value.
	msg1, _ := hex.DecodeString("d131dd02c5e6eec4693d9a0698aff95c2fcab58712467eab4004583eb8fb7f8955ad340609f4b30283e488832571415a085125e8f7cdc99fd91dbdf280373c5bd8823e3156348f5bae6dacd436c919c6dd53e2b487da03fd02396306d248cda0e99f33420f577ee8ce54b67080a80d1ec69821bcb6a8839396f9652b6ff72a70")
	msg2, _ := hex.DecodeString("d131dd02c5e6eec4693d9a0698aff95c2fcab50712467eab4004583eb8fb7f8955ad340609f4b30283e4888325f1415a085125e8f7cdc99fd91dbd7280373c5bd8823e3156348f5bae6dacd436c919c6dd53e23487da03fd02396306d248cda0e99f33420f577ee8ce54b67080280d1ec69821bcb6a8839396f965ab6ff72a70")

	if Collision(msg1, msg2) {
		fmt.Println("MD5 collision!")
	}
	// Output: MD5 collision!
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
