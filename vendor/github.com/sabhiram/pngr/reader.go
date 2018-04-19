package pngr

////////////////////////////////////////////////////////////////////////////////

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
)

////////////////////////////////////////////////////////////////////////////////

var (
	pngMagic = []byte{137, 80, 78, 71, 13, 10, 26, 10}

	ErrBadCRC = errors.New("bad crc for chunk")
)

////////////////////////////////////////////////////////////////////////////////

// Chunk describes a PNG chunk.
type Chunk struct {
	// ----------------------------------------------------------------
	// |  Length    |  Chunk Type |       ... Data ...       |  CRC   |
	// ----------------------------------------------------------------
	//    4 bytes       4 bytes         `Length` bytes         4 bytes
	Length    uint32
	ChunkType string
	Data      []byte
	Crc       uint32
}

////////////////////////////////////////////////////////////////////////////////

// ReaderOptions encapsulates various PNG chunk reader options.
type ReaderOptions struct {
	IncludedChunkTypes []string
}

// Reader implements a PNG chunk reader.
type Reader struct {
	buf  *bytes.Buffer
	opts *ReaderOptions
}

// NewReader returns a PNG chunk reader if the provided data is a valid PNG byte
// stream.
func NewReader(data []byte, opts *ReaderOptions) (*Reader, error) {
	n := len(pngMagic)
	buf := bytes.NewBuffer(data)
	magic := buf.Next(n)

	if len(magic) != n {
		return nil, errors.New("missing png file header")
	}

	i := 0
	for ; i < n; i++ {
		if magic[i] != pngMagic[i] {
			break
		}
	}
	if i != n {
		return nil, errors.New("missing png file header")
	}

	return &Reader{
		buf:  buf,
		opts: opts,
	}, nil
}

// includesChunkType returns true if the reader was created with a chunk type
// filter and if the specified type matches one of the filter entries.  If the
// reader was created with no options, all chunk types are yielded by `Next()`.
func (r *Reader) includesChunkType(ct string) bool {
	if r.opts == nil {
		return true
	}
	for _, v := range r.opts.IncludedChunkTypes {
		if v == ct {
			return true
		}
	}
	return false
}

// Next yields the next PNG chunk in the reader's buffer.  It returns an error
// of io.EOF on end of data.  It returns a bad-crc error if a given chunk is
// not constructed correctly.
func (r *Reader) Next() (*Chunk, error) {
	chunk := &Chunk{}

	for r.buf.Len() > 0 {
		// V
		// ----------------------------------------------------------------
		// |  Length    |  Chunk Type |       ... Data ...       |  CRC   |
		// ----------------------------------------------------------------
		//    4 bytes       4 bytes         `Length` bytes         4 bytes

		err := binary.Read(r.buf, binary.BigEndian, &chunk.Length)
		if err != nil {
			break
		}

		//        buf = V
		// ----------------------------------------------------------------
		// |  Length    |  Chunk Type |       ... Data ...       |  CRC   |
		// ----------------------------------------------------------------
		//    4 bytes       4 bytes         `Length` bytes         4 bytes

		minLen := 4 + int(chunk.Length)
		if r.buf.Len() < minLen {
			break
		}

		ctbs := r.buf.Next(4)
		chunk.ChunkType = string(ctbs)
		chunk.Data = r.buf.Next(int(chunk.Length))

		err = binary.Read(r.buf, binary.BigEndian, &chunk.Crc)
		if err != nil {
			break
		}

		expCrc := crc32.ChecksumIEEE(append(ctbs, chunk.Data...))
		if expCrc != chunk.Crc {
			return nil, ErrBadCRC
		}

		if r.includesChunkType(chunk.ChunkType) {
			return chunk, nil
		}
	}

	return nil, io.EOF
}
