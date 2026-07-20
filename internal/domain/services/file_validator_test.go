package services

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatingNewFieValidator(t *testing.T) {
	fv := NewFileValidator([]string{"deb", "gz"}, 1024)
	require.NotNil(t, fv)

	assert.Equal(t, []string{"deb", "gz"}, fv.allowedExtensions)
	assert.Equal(t, int64(1024*1024*1024), fv.maxFileSizeBytes)
}

func TestFileValidator_Validate(t *testing.T) {
	validator := NewFileValidator([]string{"deb", "rpm", "tar.gz", "bin", "log"}, 10*1024*1024)

	tests := []struct {
		testName      string
		setup         func(t *testing.T) string
		expectedError error
	}{
		{
			testName: "valid deb file",
			setup: func(t *testing.T) string {
				return createTempFile(t, "somefile.deb", "deb content")
			},
		},
		{
			testName: "valid tar.gz file",
			setup: func(t *testing.T) string {
				return createTempFile(t, "somefile.tar.gz", "tarball content")
			},
		},
		{
			testName: "valid log file",
			setup: func(t *testing.T) string {
				return createTempFile(t, "somefile.log", "log output")
			},
		},
		{
			testName: "file not found",
			setup: func(t *testing.T) string {
				return "/nonexistent/somefile.deb"
			},
			expectedError: domain.ErrFileNotExists,
		},
		{
			testName: "invalid extension",
			setup: func(t *testing.T) string {
				return createTempFile(t, "somefile.py", "python code")
			},
			expectedError: domain.ErrInvalidExtension,
		},
		{
			testName: "empty file",
			setup: func(t *testing.T) string {
				return createTempFile(t, "empty.deb", "")
			},
			expectedError: domain.ErrFileIsEmpty,
		},
		{
			testName: "directory instead of file",
			setup: func(t *testing.T) string {
				return t.TempDir()
			},
			expectedError: domain.ErrPathNotFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			path := tt.setup(t)
			result, err := validator.Validate(path)

			if tt.expectedError != nil {
				assert.Error(t, err)
				if tt.expectedError != nil {
					assert.ErrorIs(t, err, tt.expectedError)
				}
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.FileName)
				assert.Greater(t, result.FileSize, int64(0))
			}
		})
	}
}

func TestFileValidator_FileTooLarge(t *testing.T) {
	validator := NewFileValidator([]string{"deb"}, 1)

	path := createTempFile(
		t,
		"somefile.deb",
		strings.Repeat("string", 1024*1024),
	)
	_, err := validator.Validate(path)
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrFileTooLarge)
}

func TestGetExtension(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"package.deb", "deb"},
		{"package.dEb", "deb"},
		{"package.rpm", "rpm"},
		{"package.RPM", "rpm"},
		{"archive.tar.gz", "tar.gz"},
		{"archive.tar.bz2", "tar.bz2"},
		{"archive.tar.xz", "tar.xz"},
		{"binary.bin", "bin"},
		{"FILE.DEB", "deb"},
		{"log.log", "log"},
		{"archive.tar.zst", "tar.zst"},
		{"file.tgz", "tgz"},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			assert.Equal(t, tt.expected, getExtension(tt.filename))
		})
	}
}

func createTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0o644)
	require.NoError(t, err)
	return path
}
