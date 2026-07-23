package usecases

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/dinno7/artinux/internal/domain"
	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type readSeekCloser struct {
	io.ReadSeeker
}

func (r *readSeekCloser) Close() error { return nil }

func TestDownloadArtifactUC_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockStorage := ports.NewMockObjectStorage(ctrl)
	mockHasher := ports.NewMockChecksumHasher(ctrl)

	uc := NewDownloadArtifactUC(&noopLogger{}, mockStorage, mockHasher)

	ctx := context.Background()

	t.Run("successful download", func(t *testing.T) {
		artifact := &entities.Artifact{
			Name:     "pkg.deb",
			Checksum: "validhash",
		}
		reader := &readSeekCloser{strings.NewReader("data")}

		mockStorage.EXPECT().
			Download(ctx, "some/key").
			Return(reader, artifact, nil)

		mockHasher.EXPECT().
			CompareFromReader(gomock.Any(), "validhash").
			Return(true, nil)

		output, err := uc.Execute(ctx, DownloadArtifactInput{ObjectKey: "some/key"})
		require.NoError(t, err)
		require.NotNil(t, output)
		assert.Equal(t, artifact, output.Artifact)
	})

	t.Run("checksum mismatch returns error", func(t *testing.T) {
		artifact := &entities.Artifact{
			Name:     "pkg.deb",
			Checksum: "hash1",
		}
		reader := &readSeekCloser{strings.NewReader("data")}

		mockStorage.EXPECT().
			Download(ctx, "bad/key").
			Return(reader, artifact, nil)

		mockHasher.EXPECT().
			CompareFromReader(gomock.Any(), "hash1").
			Return(false, nil)

		_, err := uc.Execute(ctx, DownloadArtifactInput{ObjectKey: "bad/key"})
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrFileChecksumNotSame)
	})

	t.Run("storage error propagates", func(t *testing.T) {
		mockStorage.EXPECT().
			Download(ctx, "missing/key").
			Return(nil, nil, errors.New("object not found"))

		_, err := uc.Execute(ctx, DownloadArtifactInput{ObjectKey: "missing/key"})
		require.Error(t, err)
	})
}
