package ports

import "io"

type ChecksumHasher interface {
	ComputeFromReaderToBase64(reader io.Reader) (string, error)
	CompareFromReader(reader io.Reader, hashedValue string) (bool, error)

	AsWriter() io.Writer
	ComputeToBase64() string
	Reset()
}
