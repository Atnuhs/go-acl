package main

import (
	"fmt"

	"golang.org/x/exp/constraints"
)

// lIdxは左の兄弟
// return (idx / 2) * 2
func lIdx(idx int) int {
	return idx & ^1
}

// rIdxは右の兄弟
// return lIdx(idx) + 1
func rIdx(idx int) int {
	return idx | 1
}

// pIdxは親の左兄弟
func pIdx(idx int) int {
	return ((idx >> 1) - 1) & ^1
}

// cIdxは同側ヒープの最初の子
func cIdx(idx int) int {
	// 1101 -> 110 / 1 -> 11011
	// 10 -> 1 / 0 -> 110
	return (idx & ^1)<<1 | 2 | idx&1
}

// DEPQ は二重優先度キューの実装
// 二重優先度キューは、最大値と最小値をO(logN)で取得できるデータ構造
type DEPQ[T constraints.Ordered] struct {
	values []T
}

// NewDEPQ は二重優先度キューを初期化する
func NewDEPQ[T constraints.Ordered](values ...T) *DEPQ[T] {
	pq := &DEPQ[T]{
		values: values,
	}

	// Nodeの大小関係をそろえる
	for i := lIdx(pq.Size() - 1); i >= 0; i -= 2 {
		if i+1 < pq.Size() && pq.values[i] < pq.values[i+1] {
			pq.swap(i, i+1)
		}
	}

	// 親を葉側からdownしていく
	for i := pIdx(pq.Size() - 1); i >= 0; i -= 2 {
		pq.downMax(i)
		pq.downMin(i + 1)
	}
	return pq
}

func (pq *DEPQ[T]) swap(i, j int) {
	pq.values[i], pq.values[j] = pq.values[j], pq.values[i]
}

func (pq *DEPQ[T]) up(i int) {
	for {
		l, r := lIdx(i), rIdx(i)
		if r < pq.Size() && pq.values[l] < pq.values[r] {
			pq.swap(l, r)
			i ^= 1
		}

		// 末尾で奇数番目の場合
		if i&1 == 0 && i == pq.Size()-1 && i >= 2 {
			mp := pIdx(i) + 1
			if pq.values[i] < pq.values[mp] {
				pq.swap(i, mp)
				i = mp
			}
		}

		if i < 2 {
			return
		}

		p := pIdx(i) | (i & 1)
		var ok bool
		if p&1 == 0 {
			ok = pq.values[p] > pq.values[i]
		} else {
			ok = pq.values[p] < pq.values[i]
		}
		if ok {
			return
		}
		pq.swap(i, p)
		i = p
	}
}

func (pq *DEPQ[T]) downMax(i int) {
	if i&1 == 1 {
		panic(fmt.Sprintf("i should be even: i=%d", i))
	}
	for {
		lc, rc := cIdx(i), cIdx(i)+2
		if lc >= pq.Size() {
			return
		}
		big := lc
		if rc < pq.Size() && pq.values[rc] > pq.values[lc] {
			big = rc
		}

		if pq.values[i] > pq.values[big] {
			return
		}
		pq.swap(i, big)
		bl, br := lIdx(big), rIdx(big)
		if br < pq.Size() && pq.values[bl] < pq.values[br] {
			pq.swap(bl, br)
		}
		i = bl
	}
}

func (pq *DEPQ[T]) downMin(i int) {
	if i&1 == 0 {
		panic(fmt.Sprintf("i should be odd: i=%d", i))
	}
	for {
		lc, rc := cIdx(i), cIdx(i)+2
		if lc >= pq.Size() {
			return
		}
		small := lc
		if rc < pq.Size() && pq.values[rc] < pq.values[lc] {
			small = rc
		}

		if pq.values[i] < pq.values[small] {
			return
		}
		pq.swap(i, small)
		bl, br := lIdx(small), rIdx(small)
		if br < pq.Size() && pq.values[bl] < pq.values[br] {
			pq.swap(bl, br)
		}
		i = br
	}
}

func (pq *DEPQ[T]) Size() int {
	return len(pq.values)
}

func (pq *DEPQ[T]) Empty() bool {
	return len(pq.values) == 0
}

func (pq *DEPQ[T]) Push(x T) {
	pq.values = append(pq.values, x)
	pq.up(pq.Size() - 1)
}

func (pq *DEPQ[T]) GetMax() T {
	if pq.Empty() {
		panic(ErrEmptyContainer)
	}
	return pq.values[0]
}

func (pq *DEPQ[T]) GetMin() T {
	if pq.Empty() {
		panic(ErrEmptyContainer)
	}
	if pq.Size() < 2 {
		return pq.values[0]
	}
	return pq.values[1]
}

func (pq *DEPQ[T]) PopMax() T {
	if pq.Empty() {
		panic(ErrEmptyContainer)
	}
	ret := pq.values[0]
	pq.values[0] = pq.values[pq.Size()-1]
	pq.values = pq.values[:pq.Size()-1]
	pq.downMax(0)
	return ret
}

func (pq *DEPQ[T]) PopMin() T {
	if pq.Empty() {
		panic(ErrEmptyContainer)
	}
	if pq.Size() == 1 {
		ret := pq.values[0]
		pq.values = pq.values[:0]
		return ret
	}
	ret := pq.values[1]
	pq.values[1] = pq.values[pq.Size()-1]
	pq.values = pq.values[:pq.Size()-1]
	pq.downMin(1)
	return ret
}
