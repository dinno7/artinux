package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/dinno7/artinux/internal/domain/ports"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestDeleteArtifactUC_Execute(t *testing.T) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockStorage := ports.NewMockObjectStorage(ctrl)
	uc := NewDeleteArtifactUC(&noopLogger{}, mockStorage)

	ctx := context.Background()

	t.Run("deletes single artifact", func(t *testing.T) {
		mockStorage.EXPECT().
			DeleteObject(ctx, "linux/amd64/pkg.deb").
			Return(nil)

		err := uc.Execute(ctx, DeleteArtifactInput{ObjectKey: "linux/amd64/pkg.deb"})
		require.NoError(t, err)
	})

	t.Run("storage error propagates", func(t *testing.T) {
		mockStorage.EXPECT().
			DeleteObject(ctx, "missing").
			Return(errors.New("not found"))

		err := uc.Execute(ctx, DeleteArtifactInput{ObjectKey: "missing"})
		require.Error(t, err)
	})
}
