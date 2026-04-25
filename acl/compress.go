package acl

import (
	"cmp"
	"sort"
)

type Compress[T cmp.Ordered] struct {
	toOrig []T       // idx -> value
	toIdx  map[T]int // value -> idx
}

func NewCompress[T cmp.Ordered](vals []T) *Compress[T] {
	v := make([]T, len(vals))
	copy(v, vals)
	sort.Slice(v, func(i, j int) bool { return v[i] < v[j] })
	v = Uniq(v)

	m := make(map[T]int, len(v))
	for i, x := range v {
		m[x] = i
	}
	return &Compress[T]{
		toOrig: v,
		toIdx:  m,
	}
}

func (c *Compress[T]) Idx(x T) int {
	if i, ok := c.toIdx[x]; ok {
		return i
	}
	return -1
}

func (c *Compress[T]) Val(i int) T {
	return c.toOrig[i]
}

func (c *Compress[T]) Size() int {
	return len(c.toOrig)
}
