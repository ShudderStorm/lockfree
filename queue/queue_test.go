package queue_test

import (
	"sync"
	"testing"

	"github.com/ShudderStorm/lockfree/queue"
	"github.com/stretchr/testify/assert"
)

const N = 1000
const W = 100

func TestFIFO(t *testing.T) {
	q := queue.New[int]()
	for i := range N {
		q.Enqueue(i)
	}

	for expected := range N {
		val, ok := q.Dequeue()
		assert.Truef(
			t, ok,
			"element %d not found in queue", expected,
		)
		assert.Equal(
			t, expected, val,
			"expect %d, but got %d", expected, val,
		)
	}
}

func TestMultiEnqueue(t *testing.T) {
	q := queue.New[int]()
	var wg sync.WaitGroup
	for w := range W {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := w * N / W; i < (w+1)*N/W; i++ {
				q.Enqueue(i)
			}
		}()
	}
	wg.Wait()

	seen := make(map[int]bool, N)
	for val, ok := q.Dequeue(); ok; val, ok = q.Dequeue() {
		if _, ok := seen[val]; !ok {
			seen[val] = true
		} else {
			t.Fatalf("value %d was dequeued twice", val)
		}
	}

	assert.Equal(
		t, N, len(seen),
		"expected %d elements, got %d", N, len(seen),
	)
}

func TestMultiDequeue(t *testing.T) {
	q := queue.New[int]()
	for i := range N {
		q.Enqueue(i)
	}

	results := make(chan int, N)
	var wg sync.WaitGroup
	for w := range W {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := w * N / W; i < (w+1)*N/W; i++ {
				val, ok := q.Dequeue()
				assert.Truef(t, ok, "an element was lost")
				results <- val
			}
		}()
	}
	wg.Wait()
	close(results)

	seen := make(map[int]bool, N)
	for val := range results {
		if seen[val] {
			t.Fatalf("value %d was dequeued twice", val)
		}
		seen[val] = true
	}

	assert.Equal(
		t, N, len(seen),
		"expected %d elements, got %d", N, len(seen),
	)
}
