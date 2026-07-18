package hasher

import (
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/ports"
)

type sha256Hasher struct{}

func NewSha256Hasher() ports.ChecksumHasher {
	return &sha256Hasher{}
}

func (h *sha256Hasher) ComputeFromReader(reader io.Reader) (string, error) {
	hasher := sha256.New()

	if _, err := io.Copy(hasher, reader); err != nil {
		return "", domain.ErrInternal.MessageF("failed to compute sha265 hash").Wrap(err)
	}

	return base64.StdEncoding.EncodeToString(hasher.Sum(nil)), nil
}

func (h *sha256Hasher) CompareFromReader(
	reader io.Reader,
	hashedValue string,
) (bool, error) {
	computedHash, err := h.ComputeFromReader(reader)
	if err != nil {
		return false, err
	}

	return computedHash == hashedValue, nil
}
