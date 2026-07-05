package copier

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	d2 struct {
		A string
		X bool
	}
	s2 struct {
		A string
	}
	d struct {
		Name   string
		Age    int
		Height float64
		IP     *s2
		I      d2
		I2     s2
		X      bool
	}
	s struct {
		Name   string
		Age    int
		Height float64
		IP     *s2
		I      s2
		I2     s2
	}
)

func TestCopy(t *testing.T) {
	req := require.New(t)

	var r d
	err := Copy(&r, &s{"name", 1234, 1.85, &s2{"abcd"}, s2{"efgh"}, s2{"ijkl"}})
	req.Nil(err)
	req.Equal("name", r.Name)
	req.Equal(1234, r.Age)
	req.Equal(1.85, r.Height)
	req.Equal("abcd", r.IP.A)
	req.Equal("efgh", r.I.A)
}

func TestCopy2(t *testing.T) {
	req := require.New(t)

	var r d
	err := Copy2(&r, &s{"name", 1234, 1.85, &s2{"abcd"}, s2{"efgh"}, s2{"ijkl"}})
	req.Nil(err)
	req.Equal("name", r.Name)
	req.Equal(1234, r.Age)
	req.Equal(1.85, r.Height)
	req.Equal("abcd", r.IP.A)
	req.Equal("efgh", r.I.A)
	req.Equal("ijkl", r.I2.A)
}

func BenchmarkNativeCopy(b *testing.B) {
	r := make([]*d, 0, 50_000_000)
	src := &s{"name", 1234, 1.85, &s2{"abcd"}, s2{"efgh"}, s2{"ijkl"}}
	for b.Loop() {
		d := &d{
			Name:   src.Name,
			Age:    src.Age,
			Height: src.Height,
			I:      d2{A: src.I.A},
			IP:     &s2{A: src.IP.A},
		}
		r = append(r, d)
	}
}

func BenchmarkCopierCopy(b *testing.B) {
	r := make([]*d, 0, 50_000_000)
	src := &s{"name", 1234, 1.85, &s2{"abcd"}, s2{"efgh"}, s2{"ijkl"}}
	for b.Loop() {
		d, err := Copied[d](src)
		if err != nil {
			b.Fatal(err)
		}
		r = append(r, d)
	}
}

func BenchmarkCopier2Copy(b *testing.B) {
	r := make([]*d, 0, 50_000_000)
	src := &s{"name", 1234, 1.85, &s2{"abcd"}, s2{"efgh"}, s2{"ijkl"}}
	for b.Loop() {
		d, err := Copied2[d](src)
		if err != nil {
			b.Fatal(err)
		}
		r = append(r, d)
	}
}
