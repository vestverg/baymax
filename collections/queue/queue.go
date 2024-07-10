package queue

import "time"

type Queue[T any] interface {
	Offer(T)
	Peek() *T
	Poll() *T
}

type BlockingQueue[T any] interface {
	Queue[T]
	Take() *T
	TakeWithTimeout(timeout time.Duration) *T
	Interrupt()
	Len() int64
}
