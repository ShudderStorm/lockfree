package queue

import "sync/atomic"

type node[T any] struct {
	value T
	next  atomic.Pointer[node[T]]
}

type Queue[T any] struct {
	mock *node[T]
	h, t atomic.Pointer[node[T]]
}

func New[T any]() *Queue[T] {
	mock := &node[T]{}
	head := atomic.Pointer[node[T]]{}
	tail := atomic.Pointer[node[T]]{}

	head.Store(mock)
	tail.Store(mock)

	return &Queue[T]{
		mock: mock, h: head, t: tail,
	}
}

func (q *Queue[T]) Enqueue(value T) bool {
	enq := &node[T]{value: value}
	for {
		tail := q.t.Load()
		if tail.next.CompareAndSwap(nil, enq) {
			q.t.CompareAndSwap(tail, enq)
			return true
		}
		q.t.CompareAndSwap(tail, tail.next.Load())
	}
}

func (q *Queue[T]) Dequeue() (val T, ok bool) {
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
