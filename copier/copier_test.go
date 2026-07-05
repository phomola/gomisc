package copier

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	req := require.New(t)

	type d2 struct {
		A string
		X bool
	}
	type s2 struct {
		A string
	}
	type d struct {
		Name   string
		Age    int
		Height float64
		IP     *s2
		I      d2
		X      bool
	}
	type s struct {
		Name   string
		Age    int
		Height float64
		IP     *s2
		I      s2
	}

	var r d
	err := Copy(&r, &s{"name", 1234, 1.85, &s2{"abcd"}, s2{"efgh"}})
	req.Nil(err)
	req.Equal("name", r.Name)
	req.Equal(1234, r.Age)
	req.Equal(1.85, r.Height)
	req.Equal("abcd", r.IP.A)
	req.Equal("efgh", r.I.A)
}
