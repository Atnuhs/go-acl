package main

import "cmp"

type LessFunc[T any] func(a, b T) bool

// heaapImpl はヒープの実装
type PriorityQueue[T any] struct {
	data []T
	less LessFunc[T]
}

func NewPriorityQueueWithLessFunc[T any](less LessFunc[T]) *PriorityQueue[T] {
	return &PriorityQueue[T]{data: []T{}, less: less}
}

func NewPriorityQueue[T cmp.Ordered]() *PriorityQueue[T] {
	return &PriorityQueue[T]{data: []T{}, less: func(a, b T) bool { return a < b }}
}

func (pq *PriorityQueue[T]) Push(x T) {
	pq.data = append(pq.data, x)
	pq.up(len(pq.data) - 1)
}

func (pq *PriorityQueue[T]) Pop() T {
	n := len(pq.data)
	if n == 0 {
		panic(ErrEmptyContainer)
	}
	pq.swap(0, n-1)
	val := pq.data[n-1]
	pq.data = pq.data[:n-1]
	pq.down(0)
	return val
}

func (pq *PriorityQueue[T]) Peek() T {
	if pq.IsEmpty() {
		panic(ErrEmptyContainer)
	}
	return pq.data[0]
}

func (pq *PriorityQueue[T]) IsEmpty() bool {
	return len(pq.data) == 0
}

func (pq *PriorityQueue[T]) Size() int {
	return len(pq.data)
}

func (pq *PriorityQueue[T]) up(i int) {
	for {
		p := (i - 1) / 2
		if i == 0 || !pq.less(pq.data[i], pq.data[p]) {
			break
		}
		pq.swap(i, p)
		i = p
	}
}

func (pq *PriorityQueue[T]) down(i int) {
	n := len(pq.data)
	for {
		l, r := (i<<1)+1, (i<<1)+2
		smallest := i
		if l < n && pq.less(pq.data[l], pq.data[smallest]) {
			smallest = l
		}
		if r < n && pq.less(pq.data[r], pq.data[smallest]) {
			smallest = r
		}
		if smallest == i {
			break
		}
		pq.swap(i, smallest)
		i = smallest
	}
}

func (pq *PriorityQueue[T]) swap(i, j int) {
	pq.data[i], pq.data[j] = pq.data[j], pq.data[i]
}
