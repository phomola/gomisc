package pointer

import "time"

// To returns a pointer to its argument.
func To[T string | int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64 | complex64 | complex128 | time.Time](x T) *T {
	return &x
}

// Pointer is a pointer to [T].
type Pointer[T any] interface {
	*T
}
