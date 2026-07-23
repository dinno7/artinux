package domain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreatingDomainError(t *testing.T) {
	domErr := NewDomainError("test", "test error")
	require.Error(t, domErr)
}

func TestSameError(t *testing.T) {
	// TEST: Same code
	domErr := NewDomainError("test", "test error")
	domErr2 := NewDomainError("test", "some other error")
	assert.True(t, domErr.Is(domErr2))
	assert.True(t, domErr2.Is(domErr))

	// TEST: Not same
	internalErr := errors.New("this is just a testing purpose error")
	domErr = NewDomainError("test1", "test error")
	domErr.Wrap(internalErr)
	domErr2 = NewDomainError("test2", "some other error")
	domErr2.Wrap(internalErr)
	assert.False(t, domErr.Is(domErr2))
	assert.False(t, domErr2.Is(domErr))

	// TEST: Target is normal error
	internalErr = errors.New("this is just a testing purpose error")
	domErr = NewDomainError("test1", "test error")
	domErr.Wrap(internalErr)
	assert.True(t, domErr.Is(internalErr))
}

func TestUnWraping(t *testing.T) {
	internalErr := errors.New("test error")
	domErr := NewDomainError("test", "test error")
	domErr.Wrap(internalErr)
	require.Error(t, domErr)
	assert.ErrorIs(t, domErr, internalErr)
	assert.NotErrorIs(t, internalErr, domErr)
	assert.ErrorIs(t, domErr, domErr.UnWrap())
	assert.ErrorIs(t, internalErr, domErr.UnWrap())
	assert.ErrorIs(t, domErr.UnWrap(), internalErr)
}

func TestDynamicMessage(t *testing.T) {
	domErr := NewDomainError("test", "test error")
	assert.Equal(t, "test error", domErr.Message())

	domErr.MessageF("some %s error failed", "test")
	assert.Equal(t, fmt.Sprintf("some %s error failed", "test"), domErr.Message())
}

func TestErrWrapItself(t *testing.T) {
	require.NotPanics(t, func() {
		err := ErrStorageUnavailable.Wrap(ErrStorageUnavailable)
		require.NotNil(t, err)
		msg := err.Error()
		require.NotEmpty(t, msg)
		require.Contains(t, msg, "storage unavailable")
	})
}
