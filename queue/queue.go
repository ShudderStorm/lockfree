package lockfree

import "sync/atomic"

type queueNode[T any] struct {
	value T
	next  atomic.Pointer[queueNode[T]]
}

type Queue[T any] struct {
	dummy *queueNode[T]
	h, t  atomic.Pointer[queueNode[T]]
}

func NewQueue[T any]() *Queue[T] {
	dummy := &queueNode[T]{}
	head := atomic.Pointer[queueNode[T]]{}
	tail := atomic.Pointer[queueNode[T]]{}

	head.Store(dummy)
	tail.Store(dummy)

	return &Queue[T]{
		dummy: dummy, h: head, t: tail,
	}
}

func (q *Queue[T]) Enq(value T) bool {
	node := &queueNode[T]{value: value}
	for {
		tail := q.t.Load()
		if tail.next.CompareAndSwap(nil, node) {
			q.t.CompareAndSwap(tail, node)
			return true
		}
		q.t.CompareAndSwap(tail, tail.next.Load())
	}
}

func (q *Queue[T]) Deq() (val T, ok bool) {
	for {
		head, tail := q.h.Load(), q.t.Load()
		next := head.next.Load()

		if head == tail {
			if next == nil {
				return
			}
			q.t.CompareAndSwap(tail, next)
		} else {
			val = next.value
			if q.h.CompareAndSwap(head, next) {
				return val, true
			}
		}
	}
}
