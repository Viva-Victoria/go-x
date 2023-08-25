package xmath

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFloorDiv(t *testing.T) {
	t.Parallel()

	t.Run("not-panic-divide-by-zero", func(t *testing.T) {
		require.NotPanics(t, func() {
			_ = FloorDiv(1, 0)
		})
	})
	t.Run("valid-0/1", func(t *testing.T) {
		assert.Equal(t, 0, FloorDiv(0, 1))
	})
	t.Run("valid-4/3", func(t *testing.T) {
		assert.Equal(t, 1, FloorDiv(4, 3))
	})
	t.Run("valid-434/123", func(t *testing.T) {
		assert.Equal(t, 3, FloorDiv(434, 123))
	})
}
