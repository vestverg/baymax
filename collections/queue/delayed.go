package queue

import (
	"sync"
	"time"

	"github.com/vestverg/baymax/collections/heap"

	"github.com/vestverg/baymax/math"
)

type Delayed interface {
	GetDelay() int64
}

type DelayedQueue[T Delayed] struct {
	sync.RWMutex
	heap        heap.Heap[T]
	interrupted bool
}

func DelayedComparator[T Delayed](val1 T, val2 T) int {
	return math.Compare(val1.GetDelay(), val2.GetDelay())
}

func NewDelayedQueue[T Delayed]() BlockingQueue[T] {
	return &DelayedQueue[T]{
		heap: heap.NewBinaryHeap[T](DelayedComparator[T]),
	}
}

func (d *DelayedQueue[T]) Offer(t T) {
	d.Lock()
	defer d.Unlock()
	d.heap.Push(t)
}

func (d *DelayedQueue[T]) Peek() *T {
	defer d.RUnlock()
	d.RLock()
	return d.heap.Top()
}

func (d *DelayedQueue[T]) Poll() *T {
	d.Lock()
	defer d.Unlock()
	top := d.heap.Top()
	if top == nil || !((*top).GetDelay() <= 0) {
		return nil
	}
	return d.heap.Pop()
}

func (d *DelayedQueue[T]) Interrupt() {
	d.Lock()
	defer d.Unlock()
	d.interrupted = true
}

func (d *DelayedQueue[T]) Take() *T {
	d.Lock()
	d.Unlock()
	var res *T
	for res == nil || !d.interrupted {
		item := d.heap.Top()
		if item == nil {
			continue
		}
		delay := (*item).GetDelay()
		if delay <= 0 {
			res = d.heap.Pop()
			break
		}
		time.Sleep(time.Duration(delay))
	}
	return res
}

func (d *DelayedQueue[T]) TakeWithTimeout(timeout time.Duration) *T {
	d.Lock()
	d.Unlock()
	start := time.Now()
	var res *T
	for res == nil || !d.interrupted {
		item := d.heap.Top()
		if item == nil {
			continue
		}
		delay := (*item).GetDelay()
		if delay <= 0 {
			res = d.heap.Pop()
			break
		}
		if start.Add(timeout).Before(time.Now()) {
			break
		}
	}
	return res
}
