package copier

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	req := require.New(t)

	type d struct {
		Name string
		Age  int
		X    bool
	}
	type s struct {
		Name string
		Age  int
	}

	var r d
	err := Copy(&r, &s{"name", 1234})
	req.Nil(err)
	req.Equal("name", r.Name)
	req.Equal(1234, r.Age)
}
