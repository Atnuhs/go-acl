package acl

import (
	"fmt"
	"strings"
)

// SegmentTree はセグメント木の実装
type SegmentTree[T any] struct {
	data []T
	n    int
	size int
	mo   *Monoid[T]
}

// NewSegmentTree はセグメント木を初期化する
func NewSegmentTree[T any](arr []T, mo *Monoid[T]) *SegmentTree[T] {
	n := len(arr)
	size := CeilPow2(n)

	data := L1[T](2 * size)
	F1(data, mo.E)

	for i := range arr {
		data[i+size] = arr[i]
	}

	for i := size - 1; i >= 1; i-- {
		data[i] = mo.Op(data[i<<1], data[i<<1|1])
	}

	return &SegmentTree[T]{
		data: data,
		n:    n,
		size: size,
		mo:   mo,
	}
}

func (t *SegmentTree[T]) Size() int {
	return t.n
}

func (t *SegmentTree[T]) Update(i int, x T) {
	i += t.size
	t.data[i] = x
	for i > 0 {
		i >>= 1
		t.data[i] = t.mo.Op(t.data[i<<1], t.data[i<<1|1])
	}
}

func (t *SegmentTree[T]) At(i int) T {
	return t.Query(i, i+1)
}

// Queryは[l, r)の範囲でクエリ処理をする
func (t *SegmentTree[T]) Query(l, r int) T {
	l += t.size
	r += t.size
	lv, rv := t.mo.E, t.mo.E
	for l < r {
		if l&1 == 1 {
			lv = t.mo.Op(lv, t.data[l])
			l++
		}
		if r&1 == 1 {
			r--
			rv = t.mo.Op(t.data[r], rv)
		}
		l >>= 1
		r >>= 1
	}
	return t.mo.Op(lv, rv)
}

func (t *SegmentTree[T]) dump() string {
	ret := strings.Builder{}
	l := 1 << 1
	for i, v := range t.data[1:] {
		ret.WriteString(fmt.Sprintf("%d:%v ", i+1, v))
		if i+2 == l {
			ret.WriteString("\n")
			l <<= 1
		}
	}
	return ret.String()
}

// MaxRightは[l, r)に対して、Query(l, r)がTrueとなるような最大のrを求める
// ok(E)がTrueを返すことを要求する
func (t *SegmentTree[T]) MaxRight(l int, ok Ok[T]) int {
	if l < 0 || l >= t.n {
		panic("MaxRight(l) l must be [0 <= l < t.Size()]")
	}
	if !ok(t.mo.E) {
		return l
	}
	if ok(t.Query(l, t.n)) {
		return t.n
	}

	l += t.size
	lv := t.mo.E
	for {
		nlv := t.mo.Op(lv, t.data[l])
		if !ok(nlv) {
			break
		}
		if (l & 1) == 1 {
			lv = nlv
			l++
		}
		l >>= 1
	}

	for l < t.size {
		nlv := t.mo.Op(lv, t.data[l<<1])
		l <<= 1
		if ok(nlv) {
			lv = nlv
			l |= 1
		}
	}
	return l - t.size
}
