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

func (d *DelayedQueue[T]) Poll() (value *T) {
	d.Lock()
	defer d.Unlock()
	if top := d.heap.Top(); top != nil && (*top).GetDelay() <= 0 {
		value = d.heap.Pop()
	}
	return
}

func (d *DelayedQueue[T]) Interrupt() {
	d.Lock()
	defer d.Unlock()
	d.interrupted = true
}

func (d *DelayedQueue[T]) Take() *T {
	var res *T

	for res == nil || !d.interrupted {
		d.Lock()
		item := d.heap.Top()
		d.Unlock()

		if item == nil {
			continue
		}

		delay := (*item).GetDelay()
		if delay <= 0 {
			d.Lock()
			res = d.heap.Pop()
			d.Unlock()
			break
		}

		time.Sleep(time.Duration(delay))
	}

	return res
}

func (d *DelayedQueue[T]) TakeWithTimeout(timeout time.Duration) (res *T) {
	d.Lock()
	defer d.Unlock()

	start := time.Now()

	for res == nil && !d.interrupted {
		item := d.heap.Top()

		if item == nil {
			continue
		}

		delay := (*item).GetDelay()

		if delay <= 0 {
			res = d.heap.Pop()
			break
		}

		if time.Since(start) > timeout {
			break
		}
	}

	return
}
