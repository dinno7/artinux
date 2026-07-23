package usecases

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/dinno7/artinux/internal/domain/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUploadArtifactUC_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockStorage := ports.NewMockObjectStorage(ctrl)
	mockHasher := ports.NewMockChecksumHasher(ctrl)

	validator := services.NewFileValidator([]string{"deb"}, 10)
	uc := NewUploadArtifactUC(&noopLogger{}, mockStorage, mockHasher, validator)

	ctx := context.Background()

	t.Run("successful upload", func(t *testing.T) {
		content := "some deb content"
		reader := strings.NewReader(content)

		mockHasher.EXPECT().
			ComputeFromReaderToBase64(gomock.Any()).
			Return("abc123", nil)

		mockStorage.EXPECT().
			Upload(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("abc123", nil)

		objectKey, err := uc.Execute(ctx, UploadArtifactInput{
			FileName:   "pkg.deb",
			FileSize:   int64(len(content)),
			FileReader: reader,
			OS:         "linux",
			Arch:       "amd64",
			Hostname:   "host-01",
			Username:   "deploy",
		})

		require.NoError(t, err)
		assert.Contains(t, objectKey, "linux/amd64/")
		assert.Contains(t, objectKey, "pkg.deb")
	})

	t.Run("invalid os/arch returns error", func(t *testing.T) {
		_, err := uc.Execute(ctx, UploadArtifactInput{
			FileName:   "pkg.deb",
			FileSize:   100,
			FileReader: strings.NewReader("data"),
			OS:         "linux",
			Arch:       "m68k",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidOsOrArch)
	})

	t.Run("invalid file extension returns error", func(t *testing.T) {
		_, err := uc.Execute(ctx, UploadArtifactInput{
			FileName:   "script.py",
			FileSize:   100,
			FileReader: strings.NewReader("print(1)"),
			OS:         "linux",
			Arch:       "amd64",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidExtension)
	})

	t.Run("checksum mismatch triggers delete", func(t *testing.T) {
		content := "some content"
		reader := strings.NewReader(content)

		mockHasher.EXPECT().
			ComputeFromReaderToBase64(gomock.Any()).
			Return("local_hash", nil)

		mockStorage.EXPECT().
			Upload(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("different_hash", nil)

		mockStorage.EXPECT().
			DeleteObject(gomock.Any(), gomock.Any()).
			Return(nil)

		_, err := uc.Execute(ctx, UploadArtifactInput{
			FileName:   "pkg.deb",
			FileSize:   int64(len(content)),
			FileReader: reader,
			OS:         "linux",
			Arch:       "amd64",
		})

		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFileUploadNotSupportIntegrity)
	})

	t.Run("storage error propagates", func(t *testing.T) {
		content := "some content"
		reader := strings.NewReader(content)

		mockHasher.EXPECT().
			ComputeFromReaderToBase64(gomock.Any()).
			Return("hash", nil)

		mockStorage.EXPECT().
			Upload(gomock.Any(), gomock.Any(), gomock.Any()).
			Return("", errors.New("storage unavailable"))

		_, err := uc.Execute(ctx, UploadArtifactInput{
			FileName:   "pkg.deb",
			FileSize:   int64(len(content)),
			FileReader: reader,
			OS:         "linux",
			Arch:       "amd64",
		})

		require.Error(t, err)
	})

	t.Run("reader that fails to seek reports internal error", func(t *testing.T) {
		content := "content"
		failSeeker := &failReadSeeker{data: strings.NewReader(content)}

		mockHasher.EXPECT().
			ComputeFromReaderToBase64(gomock.Any()).
			Return("hash", nil)

		_, err := uc.Execute(ctx, UploadArtifactInput{
			FileName:   "pkg.deb",
			FileSize:   int64(len(content)),
			FileReader: failSeeker,
			OS:         "linux",
			Arch:       "amd64",
		})

		require.Error(t, err)
	})
}

type failReadSeeker struct {
	data io.ReadSeeker
}

func (f *failReadSeeker) Read(p []byte) (int, error)  { return f.data.Read(p) }
func (f *failReadSeeker) Seek(_ int64, _ int) (int64, error) {
	return 0, errors.New("seek failed")
}
