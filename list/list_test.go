package list

import (
	"strconv"
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

func TestFmap(t *testing.T) {
	req := require.New(t)

	l := FromSlice([]int{1, 2, 3})
	l = l.Fmap(func(x int) int { return x + 1 })
	req.Equal([]int{2, 3, 4}, l.Slice())

	l = FromSlice([]int{1, 2, 3})
	l2 := l.Fmap(func(x int) string { return strconv.Itoa(x) })
	req.Equal([]string{"1", "2", "3"}, l2.Slice())
}

func TestConcat(t *testing.T) {
	req := require.New(t)

	l := FromSlice([]int{1, 2, 3})
	l2 := FromSlice([]int{4, 5, 6})
	req.Equal([]int{1, 2, 3, 4, 5, 6}, l.Concat(l2).Slice())
}

func TestBind(t *testing.T) {
	req := require.New(t)

	l := FromSlice([]int{1, 5, 10})
	req.Equal([]int{1, 2, 5, 6, 10, 11}, l.Bind(func(x int) List[int] { return FromSlice([]int{x, x + 1}) }).Slice())
}

func TestJoin(t *testing.T) {
	req := require.New(t)

	l := FromSlice([]List[int]{
		FromSlice([]int{1, 2, 3}),
		FromSlice([]int{4, 5, 6}),
		FromSlice([]int{7, 8, 9}),
	})
	req.Equal([]int{1, 2, 3, 4, 5, 6, 7, 8, 9}, Join(l).Slice())
}

func TestInsert(t *testing.T) {
	req := require.New(t)

	req.Equal([]int{1}, FromSlice([]int{}).Insert(1, func(x, y int) bool { return x < y }).Slice())
	req.Equal([]int{1, 2}, FromSlice([]int{2}).Insert(1, func(x, y int) bool { return x < y }).Slice())
	req.Equal([]int{1, 2}, FromSlice([]int{1}).Insert(2, func(x, y int) bool { return x < y }).Slice())
}

func TestSorted(t *testing.T) {
	req := require.New(t)

	req.Equal([]int{1, 2, 3}, FromSlice([]int{3, 2, 1}).Sorted(func(x, y int) bool { return x < y }).Slice())
}

func TestEnumIndexed(t *testing.T) {
	req := require.New(t)

	var r [][]int
	for i, x := range FromSlice([]int{11, 22, 33}).EnumIndexed() {
		r = append(r, []int{i, x})
	}
	req.Equal([][]int{{0, 11}, {1, 22}, {2, 33}}, r)
}

func TestAt(t *testing.T) {
	req := require.New(t)

	l := FromSlice([]int{1, 2, 3})
	req.Equal(1, l.At(0))
	req.Equal(2, l.At(1))
	req.Equal(3, l.At(2))
}

var gr any

func BenchmarkNativeEnum(b *testing.B) {
	s := []int{1, 2, 3, 4, 5}
	var lr any
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
	var lr any
	for i := 0; i < b.N; i++ {
		r := make([]int, 0, 5)
		for x := range l.Enum() {
			r = append(r, x)
		}
		lr = r
	}
	gr = lr
}
