package mozlz4

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pierrec/lz4"
)

const magic = "mozLz40\x00"

func NewReader(r io.Reader) (io.Reader, error) {
	header := make([]byte, len(magic))
	if _, err := r.Read(header); err != nil {
		return nil, fmt.Errorf("read header error: %w", err)
	}

	if string(header) != magic {
		return nil, fmt.Errorf("unexpected header")
	}

	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, fmt.Errorf("parse size error: %w", err)
	}

	lr := io.LimitReader(r, 150*1024*1024)
	data, err := io.ReadAll(lr)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	buf := make([]byte, size)
	_, err = lz4.UncompressBlock(data, buf)
	if err != nil {
		return nil, fmt.Errorf("uncompress read error: %w", err)
	}

	mr := &mozlz4Reader{
		size: int(size),
		r:    bytes.NewReader(buf),
		nr:   0,
	}
	return mr, nil
}

type mozlz4Reader struct {
	size int
	r    io.Reader
	nr   int
}

func (r *mozlz4Reader) Read(p []byte) (n int, err error) {
	n, err = io.ReadFull(r.r, p)
	return n, err
}
