package list

import (
	"unique"
)

// List is a comparable linked list.
type List[T comparable] struct {
	head             unique.Handle[T]
	tail             unique.Handle[List[T]]
	hasHead, hasTail bool
}

// Cons creates a linked list.
func Cons[T comparable](head T, tail List[T]) List[T] {
	return List[T]{
		head:    unique.Make(head),
		tail:    unique.Make(tail),
		hasHead: true,
		hasTail: true,
	}
}

func (l List[T]) Head() T { return l.head.Value() }

func (l List[T]) Tail() List[T] { return l.tail.Value() }

func (l List[T]) IsEmpty() bool { return !l.hasHead }

func (l List[T]) IsSingleton() bool { return l.hasHead && !l.hasTail }

func (l List[T]) Len() int {
	if l.IsEmpty() {
		return 0
	}
	if l.IsSingleton() {
		return 1
	}
	return l.Tail().Len() + 1
}

// Unit creates a singleton list.
func Unit[T comparable](x T) List[T] {
	return List[T]{
		head:    unique.Make(x),
		hasHead: true,
	}
}

func (l List[T]) Enum() func(func(T) bool) {
	return func(yield func(T) bool) {
		if !l.IsEmpty() {
			if !yield(l.Head()) {
				return
			}
			if !l.IsSingleton() {
				for x := range l.Tail().Enum() {
					if !yield(x) {
						return
					}
				}
			}
		}
	}
}

func (l List[T]) Slice() []T {
	s := make([]T, 0, l.Len())
	for x := range l.Enum() {
		s = append(s, x)
	}
	return s
}

// FromSlice creates a linked list from a slice.
func FromSlice[T comparable](s []T) List[T] {
	switch len(s) {
	case 0:
		var l List[T]
		return l
	case 1:
		return Unit(s[0])
	}
	return Cons(s[0], FromSlice(s[1:]))
}
