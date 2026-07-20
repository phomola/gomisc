package ephemeral

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMap(t *testing.T) {
	req := require.New(t)

	var m Map[string, int]
	m.Set("a", new(1))
	m.Set("b", new(2))
	m.Set("c", new(3))
	x, ok := m.Get("a")
	req.True(ok)
	req.Equal(1, *x)
}
