package slice

import (
	"github.com/phomola/gomisc/maybe"
)

// Fmap is a functorial map.
func Fmap[T, U any](f func(T) U, l []T) []U {
	if l == nil {
		return nil
	}
	r := make([]U, len(l))
	for i, x := range l {
		r[i] = f(x)
	}
	return r
}

// SetFmap is a functorial map.
func SetFmap[T comparable, U any](f func(T) U, s map[T]struct{}) []U {
	r := make([]U, 0, len(s))
	for x := range s {
		r = append(r, f(x))
	}
	return r
}

// Bind is the monadic bind operation.
func Bind[T, U any](f func(T) []U, l []T) []U {
	if l == nil {
		return nil
	}
	var r []U
	for _, x := range l {
		r = append(r, f(x)...)
	}
	return r
}

// Join is the monadic join operation.
func Join[T any](x [][]T) []T {
	if x == nil {
		return nil
	}
	return Bind(maybe.Identity, x)
}

// FallibleFmap is a functorial map for a possibly erring function.
func FallibleFmap[T, U any](f func(T) (U, error), l []T) ([]U, error) {
	if l == nil {
		return nil, nil
	}
	r := make([]U, len(l))
	for i, x := range l {
		y, err := f(x)
		if err != nil {
			return nil, err
		}
		r[i] = y
	}
	return r, nil
}

// FallibleSetFmap is a functorial map for a possibly erring function.
func FallibleSetFmap[T comparable, U any](f func(T) (U, error), s map[T]struct{}) ([]U, error) {
	r := make([]U, 0, len(s))
	for x := range s {
		y, err := f(x)
		if err != nil {
			return nil, err
		}
		r = append(r, y)
	}
	return r, nil
}
