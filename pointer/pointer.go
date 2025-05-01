package pointer

// To returns a pointer to its argument.
func To[T any](x T) *T {
	return &x
}

// Pointer is a pointer to [T].
type Pointer[T any] interface {
	*T
}
