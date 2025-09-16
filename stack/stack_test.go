package stack_test

import (
	"sync"
	"testing"

	"github.com/ShudderStorm/lockfree/stack"
	"github.com/stretchr/testify/assert"
)

const N = 1000
const W = 10

func TestLIFO(t *testing.T) {
	s := stack.New[int]()
	for v := range N {
		s.Push(v)
	}

	for i := N - 1; i >= 0; i-- {
		val, ok := s.Pop()
		assert.Truef(t, ok, "element %d not found in stack", i)
		assert.Equalf(t, i, val, "expected %d, got %d", i, val)
	}
}

func TestMultiPush(t *testing.T) {
	s := stack.New[int]()
	var wg sync.WaitGroup
	for w := range W {
		wg.Add(1)
		go func(w int) {
			defer wg.Done()
			for i := w * N / W; i < (w+1)*N/W; i++ {
				s.Push(i)
			}
		}(w)
	}
	wg.Wait()

	seen := make(map[int]bool, N)
	for {
		val, ok := s.Pop()
		if !ok {
			break
		}
		if seen[val] {
			t.Fatalf("value %d was popped twice", val)
		}
		seen[val] = true
	}

	assert.Equal(t, N, len(seen), "expected %d unique elements, got %d", N, len(seen))
}

func TestMultiPop(t *testing.T) {
	s := stack.New[int]()
	for i := range N {
		s.Push(i)
	}

	results := make(chan int, N)
	var wg sync.WaitGroup
	for w := range W {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := w * N / W; i < (w+1)*N/W; i++ {
				val, ok := s.Pop()
				assert.Truef(t, ok, "element %d not found in stack", i)
				results <- val
			}
		}()
	}
	wg.Wait()
	close(results)

	seen := make(map[int]bool, N)
	for val := range results {
		if seen[val] {
			t.Fatalf("value %d was popped twice", val)
		}
		seen[val] = true
	}

	assert.Equal(t, N, len(seen), "expected %d elements, got %d", N, len(seen))
}
