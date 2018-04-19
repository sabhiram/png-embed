package pngembed

////////////////////////////////////////////////////////////////////////////////

import (
	"encoding/json"
	"fmt"
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
		k          string
		v          interface{}
		isErr      bool
		expectedSz int
	}{
		// Negative test cases.
		{data: []byte{1, 2, 3, 4}, k: "Fail", v: "FailValue", isErr: true},
		{data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, k: "Fail", v: "FailValue", isErr: true},
		{data: []byte{1, 2, 3, 4}, k: "Fail", v: nil, isErr: true},
		{data: []byte{1, 2, 3, 4}, k: "Fail", v: make(chan int), isErr: true},

		// Positive test cases.
		{data: bs, k: "Key", v: "Value0", isErr: false},
		{data: bs, k: "Key", v: 42, isErr: false},
		{data: bs, k: "Key", v: 42.0, isErr: false},
		{data: bs, k: "Key", v: struct{}{}, isErr: false},
	} {
		out, err := Embed(tc.data, tc.k, tc.v)
		if tc.isErr == false {
			fatalIfError(t, err)
		} else {
			if err == nil {
				t.Errorf("Expected error, got nil!\n")
			}
			continue
		}

		vlen := 0
		switch v := tc.v.(type) {
		case string:
			vlen = len(v)
		case int, uint:
			vlen = len(fmt.Sprintf("%d", v))
		case float32, float64:
			vlen = len(fmt.Sprintf("%f", v))
		default:
			bs, err := json.Marshal(v)
			fatalIfError(t, err)
			vlen = len(bs)
		}

		exp := len(tc.data) + len(tc.k) + 1 + vlen + 4 + 4 + len(`tEXt`)
		act := len(out)
		if act != exp {
			t.Errorf("Expected buffer size %d, got %d\n", exp, act)
		}

		m, err := Extract(out)
		fatalIfError(t, err)

		// We should have one key titled "Key".
		if 1 != len(m) {
			t.Errorf("Multiple keys found when extracting text records\n")
		}
	}
}

func TestEmbedFile(t *testing.T) {
	for _, tc := range []struct {
		fp, k, v   string
		isErr      bool
		expectedSz int
	}{
		// Negative test cases.
		{fp: "unknown.png", k: "Key0", v: "Value0", isErr: true},

		// Positive test cases.
		{fp: redPng, k: "Key0", v: "Value0", isErr: false},
	} {
		out, err := EmbedFile(tc.fp, tc.k, tc.v)
		if tc.isErr == false {
			fatalIfError(t, err)
		} else {
			if err == nil {
				t.Errorf("Expected error, got nil!\n")
			}
			continue
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

func TestBuildChunk(t *testing.T) {
	for _, tc := range []struct {
		ct         string
		data       []byte
		isErr      bool
		expectedSz int
	}{
		// Negative test cases.
		{ct: "1234", data: []byte{1, 2, 3}, isErr: true, expectedSz: 0},

		// Positive test cases.
		{ct: "tEXt", data: []byte{1, 2, 3}, isErr: false, expectedSz: 15},
	} {
		bs, err := buildChunk(tc.ct, tc.data)
		if tc.isErr == false {
			fatalIfError(t, err)
		} else {
			if err == nil {
				t.Errorf("Expected error, got nil!\n")
			}
		}

		if len(bs) != tc.expectedSz {
			t.Errorf("Expected size %d, got %d\n", tc.expectedSz, len(bs))
		}
	}
}

func TestSubstring(t *testing.T) {
	for _, tc := range []struct {
		s, sub string
		isErr  bool
	}{
		// Negative test cases.
		{"hello", "helllo", true},

		// Positive test cases.
		{"hello", "hel", false},
		{"12345", "123", false},
	} {
		err := errIfNotSubStr([]byte(tc.s), []byte(tc.sub))
		if tc.isErr == false {
			fatalIfError(t, err)
		} else {
			if err == nil {
				t.Errorf("Expected error, got nil!\n")
			}
		}
	}
}

func TestBadExtract(t *testing.T) {
	m, err := Extract([]byte{1, 2, 3})
	if m != nil {
		t.Errorf("Expected nil reader, got non-nil value\n")
	}
	if err == nil {
		t.Errorf("Expected error, got nil\n")
	}
}
