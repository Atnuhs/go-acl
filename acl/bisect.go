package acl

import "cmp"

// FirstTrue は [lo, hi) の中で ok(i) が True となる最小の i を返す。
// ok は (False → True) の方向で単調であることを要求する。
// すべて False のときは hi を返す。
func FirstTrue(lo, hi int, ok func(i int) bool) int {
	for hi-lo > 0 {
		m := int(uint(lo+hi) >> 1)
		if ok(m) {
			hi = m
		} else {
			lo = m + 1
		}
	}
	return lo
}

// LastTrue は [lo, hi) の中で ok(i) が True となる最大の i を返す。
// ok は (True → False) の方向で単調であることを要求する。
// すべて False のときは lo-1 を返す。
func LastTrue(lo, hi int, ok func(i int) bool) int {
	for hi-lo > 0 {
		m := int(uint(lo+hi) >> 1)
		if ok(m) {
			lo = m + 1
		} else {
			hi = m
		}
	}
	return lo - 1
}

// FirstGe は昇順ソート済みの配列 a に対し、x 以上の要素の左端 index を返す。
// すべての要素が x より小さい場合、len(a) を返す。
func FirstGe[T cmp.Ordered](a []T, x T) int {
	return FirstTrue(0, len(a), func(i int) bool { return a[i] >= x })
}

// FirstGt は昇順ソート済みの配列 a に対し、x より大きい要素の左端 index を返す。
// すべての要素が x 以下の場合、len(a) を返す。
func FirstGt[T cmp.Ordered](a []T, x T) int {
	return FirstTrue(0, len(a), func(i int) bool { return a[i] > x })
}

// LastLe は昇順ソート済みの配列 a に対し、x 以下の要素の右端 index を返す。
// すべての要素が x より大きい場合、-1 を返す。
func LastLe[T cmp.Ordered](a []T, x T) int {
	return LastTrue(0, len(a), func(i int) bool { return a[i] <= x })
}

// LastLt は昇順ソート済みの配列 a に対し、x より小さい要素の右端 index を返す。
// すべての要素が x 以上の場合、-1 を返す。
func LastLt[T cmp.Ordered](a []T, x T) int {
	return LastTrue(0, len(a), func(i int) bool { return a[i] < x })
}
