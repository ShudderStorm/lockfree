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

func (b *bucket[K, V]) get(key K) (value V, ok bool) {
	for node := b.head.Load(); node != nil; node = node.next.Load() {
		if node.mark.Load() {
			continue
		}

		if node.key == key {
			value, ok = node.value.Load().(V)
			if !ok {
				return
			}
			return
		}
	}
	return
}

func (b *bucket[K, V]) insert(key K, value V) {
	for {
		var prev *entry[K, V]
		current := b.head.Load()

		for current != nil && current.key != key {
			prev = current
			current = current.next.Load()
		}

		if current != nil && !current.mark.Load() && current.key == key {
			current.value.Store(value)
			return
		}

		node := &entry[K, V]{key: key}
		node.value.Store(value)
		node.next.Store(current)

		if prev == nil {
			if b.head.CompareAndSwap(current, node) {
				return
			}
		} else {
			if prev.next.CompareAndSwap(current, node) {
				return
			}
		}
	}
}

func (b *bucket[K, V]) Delete(key K) bool {
	for {
		prev := (*entry[K, V])(nil)
		current := b.head.Load()

		for current != nil && current.key != key {
			prev = current
			current = current.next.Load()
		}

		if current == nil || current.mark.Load() {
			return false
		}

		current.mark.Store(true)
		next := current.next.Load()

		if prev == nil {
			if b.head.CompareAndSwap(current, next) {
				return true
			}
		} else {
			if prev.next.CompareAndSwap(current, next) {
				return true
			}
		}
	}
}
