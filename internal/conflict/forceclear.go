package conflict

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/lacsar712/railblock/internal/bitmap"
	"github.com/lacsar712/railblock/internal/codec"
)

// Signer validates FORCE_CLEAR authorization tokens.
type Signer struct {
	secret []byte
}

// NewSigner creates an HMAC signer from a shared secret string.
func NewSigner(secret string) *Signer {
	return &Signer{secret: []byte(secret)}
}

// Sign produces a hex-encoded HMAC-SHA256 over block, source, and sequence.
func (s *Signer) Sign(blockID uint16, source bitmap.SourceID, seq uint32) string {
	payload := fmt.Sprintf("%d:%s:%d", blockID, source, seq)
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks whether a presented signature matches the expected value.
func (s *Signer) Verify(blockID uint16, source bitmap.SourceID, seq uint32, presented string) bool {
	if presented == "" {
		return false
	}
	payload := fmt.Sprintf("%d|%s|%d", blockID, source, seq)
	mac := hmac.New(sha256.New, s.secret)
	_, _ = mac.Write([]byte(payload))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(presented))
}

// ForceClearEvaluator applies rule R2 prior to mutating the bitmap.
type ForceClearEvaluator struct {
	signer *Signer
}

func NewForceClearEvaluator(signer *Signer) *ForceClearEvaluator {
	return &ForceClearEvaluator{signer: signer}
}

// AllowClear returns nil when a signed FORCE_CLEAR may proceed.
func (e *ForceClearEvaluator) AllowClear(frame *codec.Frame, source bitmap.SourceID, signature string) error {
	if frame == nil || !frame.HasForceClear() {
		return fmt.Errorf("frame does not request force clear")
	}
	if !e.signer.Verify(frame.BlockID, source, frame.Seq, signature) {
		return fmt.Errorf("invalid force clear signature")
	}
	return nil
}

// DeniedRecord builds a conflict record for rejected force-clear attempts.
func DeniedRecord(blockID uint16, reason string) Record {
	return Record{
		BlockID: blockID,
		Kind:    KindForceClearDenied,
		Message: reason,
	}
}
