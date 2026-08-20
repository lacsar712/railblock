package codec

import "fmt"

// BatchEncoder serializes multiple frames into a contiguous byte slice.
type BatchEncoder struct {
	buf []byte
	err error
}

// NewBatchEncoder creates an empty batch encoder.
func NewBatchEncoder() *BatchEncoder {
	return &BatchEncoder{}
}

// Add appends a frame to the batch, recording the first encode error.
func (b *BatchEncoder) Add(f *Frame) *BatchEncoder {
	if b.err != nil {
		return b
	}
	raw, err := Encode(f)
	if err != nil {
		b.err = err
		return b
	}
	b.buf = append(b.buf, raw...)
	return b
}

// Bytes returns the encoded batch or an error if any Add failed.
func (b *BatchEncoder) Bytes() ([]byte, error) {
	if b.err != nil {
		return nil, b.err
	}
	if len(b.buf) == 0 {
		return nil, fmt.Errorf("empty batch")
	}
	return b.buf, nil
}

// Count returns how many complete frames are in the buffer.
func (b *BatchEncoder) Count() int {
	if len(b.buf) == 0 {
		return 0
	}
	return len(b.buf) / FrameSize
}

// BatchDecoder splits a byte slice into individually validated frames.
type BatchDecoder struct {
	frames []*Frame
	err    error
}

// NewBatchDecoder parses all frames from data.
func NewBatchDecoder(data []byte) *BatchDecoder {
	frames, err := DecodeMany(data)
	return &BatchDecoder{frames: frames, err: err}
}

// Frames returns successfully parsed frames even when err is non-nil (partial parse).
func (b *BatchDecoder) Frames() []*Frame { return b.frames }

// Err returns the first decode error, if any.
func (b *BatchDecoder) Err() error { return b.err }

// Valid reports whether the entire buffer was parsed without error.
func (b *BatchDecoder) Valid() bool { return b.err == nil && len(b.frames) > 0 }
