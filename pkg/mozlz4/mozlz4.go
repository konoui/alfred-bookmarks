package mozlz4

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pierrec/lz4"
)

type Option func(r *mozlz4Config)

// WithBlockMaxSize configures size of the uncompressed data block
func WithBlockMaxSize(max int64) Option {
	return func(c *mozlz4Config) {
		c.maxSize = max
	}
}

// WithMagicHeader configures magic header. default value is `mozLz40\x00`
func WithMagicHeader(expectedHeader string) Option {
	return func(r *mozlz4Config) {
		r.magic = expectedHeader
	}
}

func NewReader(r io.Reader, opts ...Option) (io.Reader, error) {
	c := &mozlz4Config{
		maxSize: 150 * 1024 * 1024,
		magic:   "mozLz40\x00",
	}

	for _, opt := range opts {
		opt(c)
	}

	header := make([]byte, len(c.magic))
	if _, err := r.Read(header); err != nil {
		return nil, fmt.Errorf("read header error: %w", err)
	}

	if got := string(header); got != c.magic {
		return nil, fmt.Errorf("unexpected header. expected %s but got %s", c.magic, got)
	}

	var size uint32
	if err := binary.Read(r, binary.LittleEndian, &size); err != nil {
		return nil, fmt.Errorf("parse size error: %w", err)
	}

	if int64(size) > c.maxSize {
		return nil, fmt.Errorf("uncompressed size is greater than max allowed size ")
	}

	lr := io.LimitReader(r, c.maxSize)
	data, err := io.ReadAll(lr)
	if err != nil {
		return nil, fmt.Errorf("raw data read error: %w", err)
	}

	buf := make([]byte, size)
	_, err = lz4.UncompressBlock(data, buf)
	if err != nil {
		return nil, fmt.Errorf("uncompress read error: %w", err)
	}

	mr := &mozlz4Reader{
		size: int(size),
		r:    bytes.NewReader(buf),
	}
	return mr, nil
}

type mozlz4Config struct {
	maxSize int64
	magic   string
}

type mozlz4Reader struct {
	size int
	r    io.Reader
}

func (r *mozlz4Reader) Read(p []byte) (n int, err error) {
	n, err = io.ReadFull(r.r, p)
	return n, err
}
