package list

import (
	"unique"

	"github.com/fealsamh/go-utils/function"
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

func (l List[T]) Tail() List[T] {
	if l.hasTail {
		return l.tail.Value()
	}
	return List[T]{}
}

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

func (l List[T]) EnumIndexed() func(func(int, T) bool) {
	return l.enumIndexed(0)
}

func (l List[T]) enumIndexed(i int) func(func(int, T) bool) {
	return func(yield func(int, T) bool) {
		if !l.IsEmpty() {
			if !yield(i, l.Head()) {
				return
			}
			if !l.IsSingleton() {
				for i, x := range l.Tail().enumIndexed(i + 1) {
					if !yield(i, x) {
						return
					}
				}
			}
		}
	}
}

// Slice returns the list as a slice.
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

// Fmap ...
func (l List[T]) Fmap[U comparable](f func(T) U) List[U] {
	switch {
	case l.IsEmpty():
		return List[U]{}
	case l.IsSingleton():
		return Unit(f(l.Head()))
	default:
		return Cons(f(l.Head()), l.Tail().Fmap(f))
	}
}

// Concat ...
func (l List[T]) Concat(l2 List[T]) List[T] {
	switch {
	case l.IsEmpty():
		return l2
	case l.IsSingleton():
		return Cons(l.Head(), l2)
	default:
		return Cons(l.Head(), l.Tail().Concat(l2))
	}
}

// Join ...
func Join[T comparable](l List[List[T]]) List[T] {
	// switch {
	// case l.IsEmpty():
	// 	return List[T]{}
	// case l.IsSingleton():
	// 	return l.Head()
	// default:
	// 	return l.Head().Concat(Join(l.Tail()))
	// }
	return l.Bind(function.Identity)
}

// Bind ...
func (l List[T]) Bind[U comparable](f func(T) List[U]) List[U] {
	// return Join(l.Fmap(f))
	switch {
	case l.IsEmpty():
		return List[U]{}
	case l.IsSingleton():
		return f(l.Head())
	default:
		return f(l.Head()).Concat(l.Tail().Bind(f))
	}
}

// Insert returns a new list with the inserted element.
func (l List[T]) Insert(x T, less func(T, T) bool) List[T] {
	if l.IsEmpty() {
		return Unit(x)
	}
	if less(x, l.Head()) || x == l.Head() {
		return Cons(x, l)
	}
	return Cons(l.Head(), l.Tail().Insert(x, less))
}

// Sort returns the list sorted.
func (l List[T]) Sorted(less func(T, T) bool) List[T] {
	if l.IsEmpty() || l.IsSingleton() {
		return l
	}
	var l2 List[T]
	for x := range l.Enum() {
		l2 = l2.Insert(x, less)
	}
	return l2
}
