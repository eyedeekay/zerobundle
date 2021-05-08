package zerobundle

/*
Released under the The MIT License (MIT)
see ./LICENSE
*/

import (
	"testing"
)

func TestWriteTBZ(t *testing.T) {
	if err := Unpack(""); err != nil {
		t.Fatal(err)
	}
	t.Log("Success")
}
