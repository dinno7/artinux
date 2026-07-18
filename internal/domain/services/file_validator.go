package services

import (
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/pkg/helper"
)

// FileValidator validates files before upload.
type FileValidator struct {
	allowedExtensions []string
	maxFileSizeBytes  int64
}

func NewFileValidator(allowedExtensions []string, maxFileSizeMB int64) *FileValidator {
	return &FileValidator{
		allowedExtensions: allowedExtensions,
		maxFileSizeBytes:  helper.MBToBytes(maxFileSizeMB),
	}
}

type ValidationResult struct {
	FilePath  string
	FileName  string
	FileSize  int64
	Extension string
}

func (v *FileValidator) Validate(filePath string) (*ValidationResult, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrFileNotExists
		}
		if os.IsPermission(err) {
			return nil, domain.ErrFileNotAccessible
		}
		return nil, domain.ErrInternal.Wrap(err)
	}

	if info.IsDir() {
		return nil, domain.ErrPathNotFile
	}

	if info.Size() == 0 {
		return nil, domain.ErrFileIsEmpty
	}

	if info.Size() > v.maxFileSizeBytes {
		return nil, domain.ErrFileTooLarge.MessageF(
			"file size %d bytes exceeds maximum %d bytes",
			info.Size(),
			v.maxFileSizeBytes,
		)
	}

	// INFO: Check file extension
	ext := getExtension(filePath)
	if !slices.Contains(v.allowedExtensions, ext) {
		return nil, domain.ErrInvalidExtension
	}

	// INFO: Check file accessibility (can we read it?)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, domain.ErrFileNotAccessible.MessageF("file not readable: %s", filePath)
	}
	f.Close()

	return &ValidationResult{
		FilePath:  filePath,
		FileName:  filepath.Base(filePath),
		FileSize:  info.Size(),
		Extension: ext,
	}, nil
}

func getExtension(filename string) string {
	base := filepath.Base(filename)
	base = strings.ToLower(base)

	// Check for compound extensions first
	compoundExts := []string{".tar.gz", ".tar.bz2", ".tar.xz", ".tar.zst"}
	for _, ext := range compoundExts {
		if strings.HasSuffix(base, ext) {
			return normalizeFileExt(ext)
		}
	}

	return normalizeFileExt(filepath.Ext(filename))
}

func normalizeFileExt(ext string) string {
	return strings.TrimPrefix(strings.ToLower(ext), ".")
}
