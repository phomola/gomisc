package slice

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFmap(t *testing.T) {
	req := require.New(t)

	req.Equal([]string{"1", "2", "3"}, Fmap(func(x int) string { return strconv.Itoa(x) }, []int{1, 2, 3}))

	x, err := FallibleFmap(func(x int) (string, error) { return strconv.Itoa(x), nil }, []int{1, 2, 3})
	req.NoError(err)
	req.Equal([]string{"1", "2", "3"}, x)

	_, err = FallibleFmap(func(x int) (string, error) { return "", errors.ErrUnsupported }, []int{1, 2, 3})
	req.Error(err)
}

func TestJoin(t *testing.T) {
	req := require.New(t)

	req.Equal([]int{1, 2, 3, 4, 5}, Join([][]int{{1, 2}, {}, {3, 4, 5}}))
}
