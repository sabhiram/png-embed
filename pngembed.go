package pngembed

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io/ioutil"
)

var (
	pngMagic = []byte{137, 80, 78, 71, 13, 10, 26, 10}
)

// Returns nil if sub is contained in s
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
func buildTxtChunk(key, value string) []byte {
	// Add our own!
	typ := `tEXt`
	hdrb := append([]byte{}, []byte(typ)...)

	bb := []byte{}
	bb = append(bb, []byte(key)...)
	bb = append(bb, 0)
	bb = append(bb, []byte(value)...)

	// The size should include the key, the value and the null space
	szb := make([]byte, 4)
	binary.BigEndian.PutUint32(szb, uint32(len(bb)))

	// The chunk header goes before the payload
	bb = append(hdrb, bb...)

	// The CRC is calculated for the header + payload
	c := make([]byte, 4)
	crcval := crc32.ChecksumIEEE(bb)
	binary.BigEndian.PutUint32(c, crcval)

	// Prepend the size now that we have the crc
	bb = append(szb, bb...)
	// Append the CRC to the new chunk
	bb = append(bb, c...)

	return bb
}

// Embed returns a embedded png image's data stream into the file specified
// by `fpath'.  Returns error if something goes wrong!
func Embed(fpath, key, value string) ([]byte, error) {
	out := []byte{}

	// Read the image if possible
	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)

	// Magic number
	d := buf.Next(len(pngMagic))
	out = append(out, d...)
	err = errIfNotSubStr(pngMagic, d)
	if err != nil {
		return nil, err
	}

	// Extract header length, the header type should always be the first, we
	// inject our data right after this.
	d = buf.Next(4)
	out = append(out, d...)
	sz := binary.BigEndian.Uint32(d)

	// Extract the header tag, data, and CRC (for the header)
	d = buf.Next(int(sz + 8))
	out = append(out, d...)

	// Append our chunk
	out = append(out, buildTxtChunk(key, value)...)

	// Add the rest of the actual palette and data info
	out = append(out, buf.Bytes()...)

	return out, nil
}
