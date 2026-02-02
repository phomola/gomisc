package list

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromSlice(t *testing.T) {
	req := require.New(t)

	l := FromSlice([]int{})
	req.True(l.IsEmpty())
	req.False(l.IsSingleton())

	l = FromSlice([]int{1})
	req.False(l.IsEmpty())
	req.True(l.IsSingleton())

	l = FromSlice([]int{1, 2, 3})
	req.Equal([]int{1, 2, 3}, l.Slice())
}

func TestUnit(t *testing.T) {
	req := require.New(t)

	req.Equal([]int{1234}, Unit(1234).Slice())
}

func TestLen(t *testing.T) {
	req := require.New(t)

	req.Equal(0, List[int]{}.Len())
	req.Equal(1, Unit(1234).Len())
	req.Equal(5, FromSlice([]int{1, 2, 3, 4, 5}).Len())
}

func TestEnum(t *testing.T) {
	req := require.New(t)

	var s []int
	for x := range FromSlice([]int{1, 2, 3, 4, 5}).Enum() {
		s = append(s, x)
	}
	req.Equal([]int{1, 2, 3, 4, 5}, s)
}

var gr interface{}

func BenchmarkNativeEnum(b *testing.B) {
	s := []int{1, 2, 3, 4, 5}
	var lr interface{}
	for i := 0; i < b.N; i++ {
		r := make([]int, 0, 5)
		for _, x := range s {
			r = append(r, x)
		}
		lr = r
	}
	gr = lr
}

func BenchmarkListEnum(b *testing.B) {
	l := FromSlice([]int{1, 2, 3, 4, 5})
	var lr interface{}
	for i := 0; i < b.N; i++ {
		r := make([]int, 0, 5)
		for x := range l.Enum() {
			r = append(r, x)
		}
		lr = r
	}
	gr = lr
}
