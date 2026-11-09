package circuit

import (
	"context"
	"errors"

	"github.com/sony/gobreaker/v2"
)

// NewBreaker ...
func NewBreaker[T any]() *gobreaker.CircuitBreaker[T] {
	return gobreaker.NewCircuitBreaker[T](gobreaker.Settings{
		IsSuccessful: func(err error) bool { return err == nil },
		IsExcluded:   func(err error) bool { return errors.Is(err, context.Canceled) },
	})
}
