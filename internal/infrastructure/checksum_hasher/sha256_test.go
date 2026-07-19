package hasher

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComputeSha256Hash(t *testing.T) {
	hasher := NewSha256Hasher()

	reader := strings.NewReader("the simple string for test")
	expectedHash := "uODfdoUPvd3Hk9wM5ZYzgPuZMkpIZK6y5B4yUy7oPc0="

	actualHashed, err := hasher.ComputeFromReaderToBase64(reader)

	require.NoError(t, err)
	assert.Equal(t, expectedHash, actualHashed)
}

func TestCompareSha256Hash(t *testing.T) {
	hasher := NewSha256Hasher()

	reader := strings.NewReader("the simple string for test")
	expectedHash := "uODfdoUPvd3Hk9wM5ZYzgPuZMkpIZK6y5B4yUy7oPc0="

	res, err := hasher.CompareFromReader(reader, expectedHash)
	require.NoError(t, err)
	assert.True(t, res)

	expectedHash += "spoiled"
	res, err = hasher.CompareFromReader(reader, expectedHash)
	require.NoError(t, err)
	assert.False(t, res)
}

func TestAsWriter(t *testing.T) {
	hasher := NewSha256Hasher()

	reader := strings.NewReader("the simple string for test")
	expectedHash := "uODfdoUPvd3Hk9wM5ZYzgPuZMkpIZK6y5B4yUy7oPc0="

	_, _ = io.Copy(hasher.AsWriter(), reader)

	assert.Equal(t, expectedHash, hasher.ComputeToBase64())
}
