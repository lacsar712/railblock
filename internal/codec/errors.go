package codec

import "errors"

var (
	ErrBadMagic   = errors.New("codec: invalid frame magic")
	ErrBadVersion = errors.New("codec: unsupported frame version")
	ErrBadCRC     = errors.New("codec: CRC mismatch")
	ErrShortFrame = errors.New("codec: frame too short")
	ErrBadBlockID = errors.New("codec: block ID out of range")
	ErrBadOccupied = errors.New("codec: occupied field must be 0 or 1")
)

// DecodeError wraps a codec error with contextual detail for logging and API responses.
type DecodeError struct {
	Kind error
	Detail string
}

func (e *DecodeError) Error() string {
	if e.Detail == "" {
		return e.Kind.Error()
	}
	return e.Kind.Error() + ": " + e.Detail
}

func (e *DecodeError) Unwrap() error { return e.Kind }

func newDecodeError(kind error, detail string) error {
	return &DecodeError{Kind: kind, Detail: detail}
}
