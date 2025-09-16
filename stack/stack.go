package stack

import "sync/atomic"

type node[T any] struct {
	value T
	ptr   *node[T]
}

type head[T any] struct {
	ptr *node[T]
	tag uint64
}

type Stack[T any] struct {
	header atomic.Pointer[head[T]]
}

func New[T any]() *Stack[T] {
	var header atomic.Pointer[head[T]]
	header.Store(&head[T]{})
	return &Stack[T]{header: header}
}

func (s *Stack[T]) Push(value T) bool {
	push := &node[T]{value: value}

	for {
		oldHead := s.header.Load()
		push.ptr = oldHead.ptr
		newHead := &head[T]{
			ptr: push,
			tag: oldHead.tag + 1,
		}

		if s.header.CompareAndSwap(oldHead, newHead) {
			return true
		}
	}
}

func (s *Stack[T]) Pop() (val T, ok bool) {
	var zero T

	for {
		oldHead := s.header.Load()
		if oldHead.ptr == nil {
			return zero, false
		}

		top := oldHead.ptr
		newHead := &head[T]{
			ptr: top.ptr,
			tag: oldHead.tag + 1,
		}

		if s.header.CompareAndSwap(oldHead, newHead) {
			return top.value, true
		}
	}
}
