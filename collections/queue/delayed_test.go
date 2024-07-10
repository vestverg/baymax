package queue

import (
	"sync"
	"testing"
	"time"
)

type testDelayedItem struct {
	delay int64
	value string
}

func (t *testDelayedItem) GetDelay() int64 {
	return t.delay
}

func TestDelayedQueue_Offer(t *testing.T) {

	t.Run("Offer", func(t *testing.T) {
		q := NewDelayedQueue[*testDelayedItem]()
		q.Offer(&testDelayedItem{delay: 10, value: "test1"})
		q.Offer(&testDelayedItem{delay: 5, value: "test2"})

		if q.Len() != 2 {
			t.Errorf("Expected len 2, got %d", q.Len())
		}

	})
}

func TestDelayedQueue_Peek(t *testing.T) {

	t.Run("Peek", func(t *testing.T) {
		q := NewDelayedQueue[*testDelayedItem]()
		q.Offer(&testDelayedItem{delay: 10, value: "test1"})
		q.Offer(&testDelayedItem{delay: 5, value: "test2"})

		peeked := q.Peek()
		if peeked == nil || (*peeked).value != "test2" {
			t.Errorf("Expected test2, got %v", peeked)
		}

		if q.Len() != 2 {
			t.Errorf("Expected len 2, got %d", q.Len())
		}
	})

}

func TestDelayedQueue_Poll(t *testing.T) {
	t.Run("Poll", func(t *testing.T) {

		q := NewDelayedQueue[*testDelayedItem]()
		q.Offer(&testDelayedItem{delay: -1, value: "test1"}) //immediately available
		q.Offer(&testDelayedItem{delay: 5, value: "test2"})  //delayed

		polled := q.Poll()
		if polled == nil || (*polled).value != "test1" {
			t.Errorf("Expected test1, got %v", polled)
		}

		if q.Len() != 1 {
			t.Errorf("Expected len 1, got %d", q.Len())
		}

		polled = q.Poll() //test2 is not ready yet
		if polled != nil {
			t.Errorf("Expected nil, got %v", polled)
		}

		q.Offer(&testDelayedItem{delay: 0, value: "test3"}) //immediately available
		polled = q.Poll()
		if polled == nil || (*polled).value != "test3" {
			t.Errorf("Expected test3, got %v", polled)
		}
	})
}

func TestDelayedQueue_Take(t *testing.T) {
	t.Run("Take", func(t *testing.T) {

		q := NewDelayedQueue[*testDelayedItem]()
		q.Offer(&testDelayedItem{delay: 0, value: "test1"})

		taken := q.Take()

		if taken == nil || (*taken).value != "test1" {
			t.Error("Expected test1")
		}

		q = NewDelayedQueue[*testDelayedItem]()

		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			defer wg.Done()
			taken := q.Take()
			if taken == nil || (*taken).value != "test1" {
				t.Error("Expected test1")
			}
		}()

		time.Sleep(time.Millisecond * 10)
		q.Offer(&testDelayedItem{delay: 0, value: "test1"})
		wg.Wait()

		q = NewDelayedQueue[*testDelayedItem]()
		wg.Add(1)
		go func() {
			defer wg.Done()
			taken := q.Take()
			if taken != nil {
				t.Errorf("Expected nil due to interrupt, but got %+v", taken)
			}

		}()
		time.Sleep(time.Millisecond * 10)
		q.Interrupt()
		wg.Wait()
	})

}

func TestDelayedQueue_TakeWithTimeout(t *testing.T) {

	t.Run("TakeWithTimeout", func(t *testing.T) {
		q := NewDelayedQueue[*testDelayedItem]()
		q.Offer(&testDelayedItem{delay: 100, value: "test1"})

		taken := q.TakeWithTimeout(10 * time.Millisecond)

		if taken != nil {
			t.Errorf("Expected nil, got %v", taken)
		}

		q = NewDelayedQueue[*testDelayedItem]()
		q.Offer(&testDelayedItem{delay: 0, value: "test1"})

		taken = q.TakeWithTimeout(10 * time.Millisecond)

		if taken == nil || (*taken).value != "test1" {
			t.Error("Expected test1")
		}

		q = NewDelayedQueue[*testDelayedItem]()
		q.Offer(&testDelayedItem{delay: 100, value: "test1"})

		go func() {
			time.Sleep(5 * time.Millisecond)
			q.Interrupt()
		}()

		taken = q.TakeWithTimeout(100 * time.Millisecond)

		if taken != nil {
			t.Errorf("Expected nil due to interrupt")
		}
	})

}
