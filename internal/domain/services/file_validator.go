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
	File      *os.File
}

func (v *FileValidator) ValidateAndGetExt(
	fileName string,
	fileSize int64,
) (string, error) {
	if err := validateSize(fileSize, v.maxFileSizeBytes); err != nil {
		return "", err
	}

	ext := getExtension(fileName)
	// INFO: Check file extension
	if err := validateFileExt(ext, v.allowedExtensions); err != nil {
		return "", err
	}
	return ext, nil
}

func (v *FileValidator) ValidateFile(filePath string) (*ValidationResult, error) {
	// INFO: Check file accessibility (can we read it?)
	f, err := os.Open(filePath)
	if err != nil {
		return nil, domain.ErrFileNotAccessible.MessageF("file not readable: %s", filePath)
	}

	info, err := f.Stat()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, domain.ErrFileNotExists
		}
		if os.IsPermission(err) {
			return nil, domain.ErrFileNotAccessible
		}
		return nil, domain.ErrInternal.Wrap(err)
	}

	// INFO: Check file info
	if err := validateFileInfo(info, v.maxFileSizeBytes); err != nil {
		return nil, err
	}

	// TODO: Validate file type via file's magic number

	ext := getExtension(filePath)
	// INFO: Check file extension
	if err := validateFileExt(ext, v.allowedExtensions); err != nil {
		return nil, err
	}

	return &ValidationResult{
		FilePath:  filePath,
		FileName:  filepath.Base(filePath),
		FileSize:  info.Size(),
		Extension: ext,
		File:      f,
	}, nil
}

func validateFileInfo(info os.FileInfo, maxFileSizeBytes int64) error {
	if info.IsDir() {
		return domain.ErrPathNotFile
	}
	err := validateSize(info.Size(), maxFileSizeBytes)
	if err != nil {
		return err
	}
	return nil
}

func validateSize(size int64, maxFileSizeBytes int64) error {
	if size == 0 {
		return domain.ErrFileIsEmpty
	}

	if size > maxFileSizeBytes {
		return domain.ErrFileTooLarge.MessageF(
			"file size %d bytes exceeds maximum %d bytes",
			size,
			maxFileSizeBytes,
		)
	}
	return nil
}

func validateFileExt(ext string, allowedExtensions []string) error {
	if !slices.Contains(allowedExtensions, ext) {
		return domain.ErrInvalidExtension
	}
	return nil
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
