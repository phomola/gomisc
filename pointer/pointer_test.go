package pointer

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestPointerTo(t *testing.T) {
	req := require.New(t)
	{
		p := To("abcd")
		req.Equal("abcd", *p)
	}
	{
		t := time.Now()
		p := To(t)
		req.Equal(t, *p)
	}
}
