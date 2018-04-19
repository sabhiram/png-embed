// Package pngembed embeds key-value data into a png image.
package pngembed

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"strings"

	"github.com/sabhiram/pngr"
)

////////////////////////////////////////////////////////////////////////////////

var (
	pngMagic = []byte{137, 80, 78, 71, 13, 10, 26, 10}
)

////////////////////////////////////////////////////////////////////////////////

// Returns nil if sub is contained in s, an error otherwise.
func errIfNotSubStr(s, sub []byte) error {
	if len(sub) > len(s) {
		return errors.New("substring larger than parent")
	}
	for i, d := range sub {
		if d != s[i] {
			return errors.New("byte mismatch with sub")
		}
	}
	return nil
}

func isValidChunkType(ct string) bool {
	for _, v := range []string{
		// Critical chunks.
		"IHDR", "PLTE", "IDAT", "IEND",

		// Ancillary chunks.
		"bKGD", "cHRM", "dSIG", "eXIf", "gAMA", "hIST", "iCCP", "iTXt", "pHYs",
		"sBIT", "sPLT", "sRGB", "sTER", "tEXt", "tIME", "tRNS", "zTXt",
	} {
		if v == ct {
			return true
		}
	}
	return false
}

// buildChunk encodes the specified chunk type and data into a png chunk.  If
// the chunk type is invalid, it is rejected.
func buildChunk(ct string, data []byte) ([]byte, error) {
	// -------------------------------------------------------------------
	// |  Length    |  Chunk Type |       ... Data ...       |    CRC    |
	// -------------------------------------------------------------------
	// |  4 bytes   |   4 bytes   |     `Length` bytes       |  4 bytes  |
	//              |-------------- CRC32'd -----------------|
	if !isValidChunkType(ct) {
		return nil, fmt.Errorf("invalid chunk type (%s)", ct)
	}

	szbs := make([]byte, 4)
	binary.BigEndian.PutUint32(szbs, uint32(len(data)))

	bb := append([]byte(ct), data...)

	crcbs := make([]byte, 4)
	binary.BigEndian.PutUint32(crcbs, crc32.ChecksumIEEE(bb))

	bb = append(bb, crcbs...)

	// Prepend the length to the payload.
	return append(szbs, bb...), nil
}

func buildTextChunk(data []byte) []byte {
	bs, _ := buildChunk(`tEXt`, data)
	return bs
}

// embed verifies that the input data slice actually describes a PNG image, and
// appends the respective (key, value) pair into its `tExt` section(s).
func embed(data []byte, k string, v []byte) ([]byte, error) {
	out := []byte{}
	buf := bytes.NewBuffer(data)

	// Magic number.
	d := buf.Next(len(pngMagic))
	out = append(out, d...)
	err := errIfNotSubStr(pngMagic, d)
	if err != nil {
		return nil, err
	}

	// Extract header length, the header type should always be the first, we
	// inject our custom text data right after this.
	d = buf.Next(4)
	out = append(out, d...)
	sz := binary.BigEndian.Uint32(d)

	// Extract the header tag, data, and CRC (for the header).
	d = buf.Next(int(sz + 8))
	out = append(out, d...)

	// Append tEXt chunk.
	out = append(out, buildTextChunk(append(append([]byte(k), 0), v...))...)

	// Add the rest of the actual palette and data info.
	return append(out, buf.Bytes()...), nil
}

////////////////////////////////////////////////////////////////////////////////

// Embed accepts a stream of bytes which represent the raw PNG image data, and
// the `key` to store the interface `v` under.  `v` is treated as JSON which
// when Marshal'd will result in either a JSON string representing a map, or
// the serialized value of a primitive type (int, string, float etc). Returns
// the raw bytes that represent the modified PNG data.
func Embed(data []byte, k string, v interface{}) ([]byte, error) {
	var (
		err error
		val []byte
	)

	switch vt := v.(type) {
	case int, uint:
		val = []byte(fmt.Sprintf("%d", vt))
	case float32, float64:
		val = []byte(fmt.Sprintf("%f", vt))
	case string:
		val = []byte(vt)
	default:
		val, err = json.Marshal(v)
	}

	if err != nil {
		return nil, err
	}
	return embed(data, k, val)
}

// EmbedFile is like `Embed` but accepts the path to a PNG file instead of the
// raw png data.
func EmbedFile(fp, k string, v interface{}) ([]byte, error) {
	data, err := ioutil.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	return Embed(data, k, v)
}

////////////////////////////////////////////////////////////////////////////////

func Extract(data []byte) (map[string][]byte, error) {
	ret := map[string][]byte{}

	r, err := pngr.NewReader(data, &pngr.ReaderOptions{
		IncludedChunkTypes: []string{`tEXt`},
	})
	if err != nil {
		return nil, err
	}

	c, err := r.Next()
	for ; err == nil; c, err = r.Next() {
		sz := len(c.Data)
		pt := strings.Index(string(c.Data), string(0))
		if pt < sz {
			ret[string(c.Data[:pt])] = c.Data[pt+1:]
		}
	}
	if err == io.EOF {
		err = nil
	}

	return ret, err
}
