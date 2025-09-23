package hashmap

import "sync/atomic"

type table[K comparable, V any] struct {
	buckets    []bucket[K, V]
	size, mask uint64
}

type HashMap[K comparable, V any] struct {
	current, old atomic.Pointer[table[K, V]]
}
