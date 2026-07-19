package hasher

import (
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"io"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/ports"
)

type sha256Hasher struct {
	hasher hash.Hash
}

func NewSha256Hasher() ports.ChecksumHasher {
	return &sha256Hasher{
		hasher: sha256.New(),
	}
}

const buffer_size = 1 * 1024 * 1024 // 1MB

func (h *sha256Hasher) ComputeFromReaderToBase64(reader io.Reader) (string, error) {
	hasher := sha256.New()

	buf := make([]byte, buffer_size)
	if _, err := io.CopyBuffer(hasher, reader, buf); err != nil {
		return "", domain.ErrInternal.MessageF("failed to compute sha265 hash").Wrap(err)
	}

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func (h *sha256Hasher) CompareFromReader(
	reader io.Reader,
	hashedValue string,
) (bool, error) {
	computedHash, err := h.ComputeFromReaderToBase64(reader)
	if err != nil {
		return false, err
	}

	return computedHash == hashedValue, nil
}

func (h *sha256Hasher) ComputeToBase64() string {
	return base64.StdEncoding.EncodeToString(h.hasher.Sum(nil))
}

func (h *sha256Hasher) AsWriter() io.Writer {
	return h.hasher
}

func (h *sha256Hasher) Reset() {
	h.hasher.Reset()
}
