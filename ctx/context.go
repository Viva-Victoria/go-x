package ctx

import (
	"context"
	"time"

	"github.com/Viva-Victoria/go-x/xmath"
)

type Context interface {
	context.Context
	WithTimeout(duration time.Duration) (Context, CancelFunc)
	WithDeadline(deadline time.Time) (Context, CancelFunc)
	WithCancel() (Context, CancelFunc)
	Values() map[any]any
}

type CancelFunc context.CancelFunc

type xContext struct {
	parent context.Context
	values map[any]any
}

func Wrap(c context.Context, values ...any) Context {
	return &xContext{
		parent: c,
		values: cloneValues(values),
	}
}

func cloneValues(values []any) map[any]any {
	var (
		valuesCount = xmath.FloorDiv(len(values), 2) * 2
		valuesMap   = make(map[any]any, valuesCount)
	)

	if len(values) < 1 {
		return valuesMap
	}

	if parentValues, ok := values[0].(map[any]any); ok {
		for k, v := range parentValues {
			valuesMap[k] = v
		}

		return valuesMap
	}

	for i := 0; i < valuesCount; i += 2 {
		valuesMap[values[i]] = values[i+1]
	}

	return valuesMap
}

func New(values ...any) Context {
	return Wrap(context.Background(), values...)
}

func (x *xContext) Deadline() (deadline time.Time, ok bool) {
	return x.parent.Deadline()
}

func (x *xContext) Done() <-chan struct{} {
	return x.parent.Done()
}

func (x *xContext) Err() error {
	return x.parent.Err()
}

func (x *xContext) Value(key any) any {
	if v, ok := x.values[key]; ok {
		return v
	}

	return x.parent.Value(key)
}

func (x *xContext) WithTimeout(duration time.Duration) (Context, CancelFunc) {
	c, f := context.WithTimeout(x, duration)
	return Wrap(c, x.values), CancelFunc(f)
}

func (x *xContext) WithDeadline(deadline time.Time) (Context, CancelFunc) {
	c, f := context.WithDeadline(x, deadline)
	return Wrap(c, x.values), CancelFunc(f)
}

func (x *xContext) WithCancel() (Context, CancelFunc) {
	c, f := context.WithCancel(x)
	return Wrap(c, x.values), CancelFunc(f)
}

func (x *xContext) Values() map[any]any {
	return x.values
}
