package ephemeral

import (
	"runtime"
	"sync"
	"weak"
)

// Map ...
type Map[K any, V any] struct {
	m sync.Map
}

// Set ...
func (m *Map[K, V]) Set(k K, v *V) {
	runtime.AddCleanup(v, func(k K) { m.m.Delete(k) }, k)
	m.m.Store(k, weak.Make(v))
}

// Get ...
func (m *Map[K, V]) Get(k K) (*V, bool) {
	v, ok := m.m.Load(k)
	if !ok {
		return nil, false
	}
	x := v.(weak.Pointer[V]).Value()
	return x, x != nil
}
