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
	"io/ioutil"
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

// buildTxtChunk builds a given text chunk based on a key and value with the
// correct CRC.
func buildTxtChunk(key string, value []byte) []byte {
	// Header
	typ := `tEXt`
	hdrb := append([]byte{}, []byte(typ)...)

	// Payload
	bb := []byte{}
	bb = append(bb, []byte(key)...)
	bb = append(bb, 0)
	bb = append(bb, value...)

	// Size
	szb := make([]byte, 4)
	binary.BigEndian.PutUint32(szb, uint32(len(bb)))

	// Prepend the header to the payload
	bb = append(hdrb, bb...)

	// CRC32
	c := make([]byte, 4)
	crcval := crc32.ChecksumIEEE(bb)
	binary.BigEndian.PutUint32(c, crcval)

	// Prepend the size now that we have the crc
	bb = append(szb, bb...)

	// Append the CRC to the new chunk
	return append(bb, c...)
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

	// Append our chunk.
	out = append(out, buildTxtChunk(k, v)...)

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
