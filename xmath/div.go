package xmath

import "golang.org/x/exp/constraints"

func FloorDiv[T constraints.Integer](a, b T) T {
	d := float64(a) / float64(b)
	return T(d)
}
