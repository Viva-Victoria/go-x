package ctx

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCloneValues(t *testing.T) {
	testCases := map[string]struct {
		result map[any]any
		source []any
	}{
		"empty": {
			source: nil,
			result: make(map[any]any),
		},
		"parent": {
			source: []any{
				map[any]any{
					1:      "int-1",
					33.2:   "float-33.2",
					"text": "value",
				},
			},
			result: map[any]any{
				1:      "int-1",
				33.2:   "float-33.2",
				"text": "value",
			},
		},
		"flat-full": {
			source: []any{
				1, 2.2,
				"4", []string{"a", "b", "c"},
			},
			result: map[any]any{
				1:   2.2,
				"4": []string{"a", "b", "c"},
			},
		},
		"flat-incomplete": {
			source: []any{
				1, 2.2,
				"4",
			},
			result: map[any]any{
				1: 2.2,
			},
		},
	}

	for name, data := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.EqualValues(t, data.result, cloneValues(data.source))
		})
	}
}

func TestNew(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		xctx := New()
		assert.NotNil(t, xctx)
	})
	t.Run("values", func(t *testing.T) {
		xctx := New("key", "value")
		assert.NotNil(t, xctx)
		assert.EqualValues(t, map[any]any{
			"key": "value",
		}, xctx.Values())
	})
}

func TestXContext_Deadline(t *testing.T) {
	deadline := time.Now().Add(time.Hour * 72)

	parent, parentCancel := context.WithDeadline(context.Background(), deadline)
	defer parentCancel()

	parentDeadline, _ := parent.Deadline()
	assert.Equal(t, deadline, parentDeadline)

	child, childCancel := New().WithDeadline(deadline)
	defer childCancel()

	childDeadline, _ := child.Deadline()
	assert.Equal(t, deadline, childDeadline)
}

func TestXContext_Done(t *testing.T) {
	parent := context.Background()

	t.Run("nil", func(t *testing.T) {
		child := Wrap(parent)
		assert.Equal(t, parent.Done(), child.Done())
	})
	t.Run("with-cancel", func(t *testing.T) {
		parentCancelContext, parentCancel := context.WithCancel(parent)
		defer parentCancel()

		child, cancel := Wrap(parent).WithCancel()
		defer cancel()

		assert.NotEqual(t, parentCancelContext.Done(), child.Done())
	})
}

func TestXContext_Err(t *testing.T) {
	parent := context.Background()
	timeout := -time.Second
	deadline := time.Now().Add(timeout)

	t.Run("deadline", func(t *testing.T) {
		parentDeadline, parentCancel := context.WithDeadline(parent, deadline)
		defer parentCancel()

		assert.Error(t, context.DeadlineExceeded, parentDeadline.Err())

		childDeadline, childCancel := New().WithDeadline(deadline)
		defer childCancel()

		assert.Equal(t, context.DeadlineExceeded, childDeadline.Err())
	})
	t.Run("timeout", func(t *testing.T) {
		parentTimeout, parentCancel := context.WithTimeout(parent, timeout)
		defer parentCancel()

		assert.Error(t, context.DeadlineExceeded, parentTimeout.Err())

		childTimeout, childCancel := New().WithTimeout(timeout)
		defer childCancel()

		assert.Equal(t, context.DeadlineExceeded, childTimeout.Err())
	})
	t.Run("cancel", func(t *testing.T) {
		parentCancelContext, parentCancel := context.WithCancel(parent)
		parentCancel()

		assert.Error(t, context.Canceled, parentCancelContext.Err())

		child, childCancel := New().WithCancel()
		childCancel()

		assert.Error(t, context.Canceled, child.Err())
	})
}

func TestXContext_Value(t *testing.T) {
	var (
		key = struct {
			name string
		}{
			name: "key",
		}
		value = "value"
	)

	parent := context.WithValue(context.Background(), key, value)
	child := Wrap(parent)

	t.Run("from-parent", func(t *testing.T) {
		assert.Equal(t, value, child.Value(key))
	})
	t.Run("from-child", func(t *testing.T) {
		subChild := Wrap(child, key, value+value)
		assert.Equal(t, value+value, subChild.Value(key))
	})
}
