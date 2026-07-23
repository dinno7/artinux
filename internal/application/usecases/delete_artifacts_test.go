package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteArtifactsUC_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockStorage := ports.NewMockObjectStorage(ctrl)
	uc := NewDeleteArtifactsUC(&noopLogger{}, mockStorage)

	ctx := context.Background()

	t.Run("deletes multiple artifacts", func(t *testing.T) {
		keys := []string{"a.deb", "b.deb"}

		mockStorage.EXPECT().
			DeleteBatch(ctx, keys).
			Return(nil)

		err := uc.Execute(ctx, DeleteArtifactsInput{ObjectKeys: keys})
		require.NoError(t, err)
	})

	t.Run("storage error propagates", func(t *testing.T) {
		keys := []string{"x.deb"}

		mockStorage.EXPECT().
			DeleteBatch(ctx, keys).
			Return(errors.New("batch failed"))

		err := uc.Execute(ctx, DeleteArtifactsInput{ObjectKeys: keys})
		require.Error(t, err)
	})
}
