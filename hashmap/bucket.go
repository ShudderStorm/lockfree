package hashmap

import "sync/atomic"

type entry[K comparable, V any] struct {
	key   K
	value atomic.Value

	next atomic.Pointer[entry[K, V]]
	mark atomic.Bool
}

type bucket[K comparable, V any] struct {
	head atomic.Pointer[entry[K, V]]
}
