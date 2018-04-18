package pngembed

////////////////////////////////////////////////////////////////////////////////

import (
	"io/ioutil"
	"os"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////

const (
	redPng = "./fixtures/red.png"
)

////////////////////////////////////////////////////////////////////////////////

func fatalIfError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Fatal error: %s\n", err.Error())
	}
}

func TestEmbed(t *testing.T) {
	bs, err := ioutil.ReadFile(redPng)
	fatalIfError(t, err)

	for _, tc := range []struct {
		data       []byte
		k, v       string
		isErr      bool
		expectedSz int
	}{
		{data: bs, k: "Key0", v: "Value0", isErr: false},
	} {
		out, err := Embed(tc.data, tc.k, tc.v)
		if tc.isErr == false {
			fatalIfError(t, err)
		} else {
			if err == nil {
				t.Errorf("Expected error, got nil!\n")
			}
		}

		exp := len(tc.data) + len(tc.k) + 1 + len(tc.v) + 4 + 4 + len(`tEXt`)
		act := len(out)
		if act != exp {
			t.Errorf("Expected buffer size %d, got %d\n", exp, act)
		}
	}
}

func TestEmbedFile(t *testing.T) {
	for _, tc := range []struct {
		fp, k, v   string
		isErr      bool
		expectedSz int
	}{
		{fp: redPng, k: "Key0", v: "Value0", isErr: false},
	} {
		out, err := EmbedFile(tc.fp, tc.k, tc.v)
		if tc.isErr == false {
			fatalIfError(t, err)
		} else {
			if err == nil {
				t.Errorf("Expected error, got nil!\n")
			}
		}

		fi, err := os.Stat(tc.fp)
		fatalIfError(t, err)

		exp := int(fi.Size()) + len(tc.k) + 1 + len(tc.v) + 4 + 4 + len(`tEXt`)
		act := len(out)
		if act != exp {
			t.Errorf("Expected buffer size %d, got %d\n", exp, act)
		}
	}
}
