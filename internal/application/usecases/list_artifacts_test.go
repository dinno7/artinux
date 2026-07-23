package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/dinno7/artinux/internal/domain/entities"
	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListArtifactUC_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockStorage := ports.NewMockObjectStorage(ctrl)
	uc := NewListArtifactUC(&noopLogger{}, mockStorage)

	ctx := context.Background()

	t.Run("returns artifacts", func(t *testing.T) {
		expected := []*entities.Artifact{
			{Name: "a.deb"},
			{Name: "b.deb"},
		}

		mockStorage.EXPECT().
			ListObjects(ctx, "linux/amd64", 10).
			Return(expected, nil)

		got, err := uc.Execute(ctx, ListArtifactInput{Prefix: "linux/amd64", Limit: 10})
		require.NoError(t, err)
		assert.Len(t, got, 2)
		assert.Equal(t, expected, got)
	})

	t.Run("returns empty list", func(t *testing.T) {
		mockStorage.EXPECT().
			ListObjects(ctx, "", 10).
			Return([]*entities.Artifact{}, nil)

		got, err := uc.Execute(ctx, ListArtifactInput{Limit: 10})
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("defaults limit when zero", func(t *testing.T) {
		mockStorage.EXPECT().
			ListObjects(ctx, "", 10).
			Return([]*entities.Artifact{}, nil)

		_, err := uc.Execute(ctx, ListArtifactInput{Limit: 0})
		require.NoError(t, err)
	})

	t.Run("error from storage propagates", func(t *testing.T) {
		mockStorage.EXPECT().
			ListObjects(ctx, "", 10).
			Return(nil, errors.New("storage error"))

		_, err := uc.Execute(ctx, ListArtifactInput{Limit: 10})
		require.Error(t, err)
	})
}
