package queue

import "sync/atomic"

type queueNode[T any] struct {
	value T
	next  atomic.Pointer[queueNode[T]]
}

type Queue[T any] struct {
	mock *queueNode[T]
	h, t atomic.Pointer[queueNode[T]]
}

func New[T any]() *Queue[T] {
	mock := &queueNode[T]{}
	head := atomic.Pointer[queueNode[T]]{}
	tail := atomic.Pointer[queueNode[T]]{}

	head.Store(mock)
	tail.Store(mock)

	return &Queue[T]{
		mock: mock, h: head, t: tail,
	}
}

func (q *Queue[T]) Enqueue(value T) bool {
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
