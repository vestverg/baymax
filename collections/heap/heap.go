package heap

import (
	"sort"

	"github.com/vestverg/baymax/generics"
)

type Comparable[T any] interface {
	Compare(T) int
}

type Heap[T any] interface {
	sort.Interface
	Push(T)
	Pop() *T
	Top() *T
	Bottom() *T
}

func NewBinaryHeap[T any](comparator generics.Comparator[T], el ...T) Heap[T] {
	h := &BinaryHeap[T]{
		arr:        el,
		comparator: comparator,
	}
	h.init()
	return h
}

func (b *BinaryHeap[T]) init() {
	n := b.Len()
	for i := n/2 - 1; i >= 0; i-- {
		b.down(i, n)
	}

}

type BinaryHeap[T any] struct {
	arr        []T
	comparator generics.Comparator[T]
}

func (b *BinaryHeap[T]) Len() int {
	return len(b.arr)
}

func (b *BinaryHeap[T]) Less(i, j int) bool {
	if i > b.Len() || j > b.Len() {
		panic("invalid argument")
	}
	return b.comparator(b.arr[i], b.arr[j]) < 0
}

func (b *BinaryHeap[T]) Swap(i, j int) {
	if i > b.Len() || j > b.Len() {
		panic("invalid argument")
	}

	b.arr[i], b.arr[j] = b.arr[j], b.arr[i]
}

func (b *BinaryHeap[T]) Push(item T) {
	b.arr = append(b.arr, item)
	b.up(b.Len() - 1)
}

func (b *BinaryHeap[T]) Pop() *T {
	if b.Len() == 0 {
		return nil
	}
	n := b.Len() - 1
	b.Swap(0, n)
	b.down(0, n)

	old_item := b.arr[b.Len()-1]
	b.arr = b.arr[0 : b.Len()-1]
	return &old_item
}

func (b *BinaryHeap[T]) Top() *T {

	if b.Len() == 0 {
		return nil
	}
	return &b.arr[0]
}

func (b *BinaryHeap[T]) Bottom() *T {
	if b.Len() == 0 {
		return nil
	}
	return &b.arr[b.Len()-1]
}

func (b *BinaryHeap[T]) up(idx int) {
	for {
		parent := (idx - 1) / 2
		if parent == idx || !b.Less(idx, parent) {
			break
		}
		b.Swap(parent, idx)
		idx = parent
	}
}

func (b *BinaryHeap[T]) down(idx int, n int) {
	for {
		left := 2*idx + 1
		if left >= n || left < 0 {
			break
		}
		next := left
		if right := next + 1; right < n && b.Less(right, left) {
			next = right
		}
		if !b.Less(next, idx) {
			break
		}
		b.Swap(idx, next)
		idx = next
	}
}
