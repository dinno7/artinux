package ports

import "io"

type ChecksumHasher interface {
	ComputeFromReader(reader io.Reader) (string, error)
	CompareFromReader(reader io.Reader, hashedValue string) (bool, error)
}
