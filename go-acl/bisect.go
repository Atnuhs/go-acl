package main

import "cmp"

func Bisect(ok, ng int, pred func(mid int) bool) int {
	for Abs(ng-ok) > 1 {
		mid := (ok + ng) >> 1
		if pred(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// Geは昇順ソート済みの配列aに対して、x以上の要素の左端Indexを返す
// aのすべての要素がxより小さい場合、len(a)を返す
func Ge[T cmp.Ordered](a []T, x T) int {
	ok, ng := len(a), -1
	for ok-ng > 1 {
		m := (ok + ng) >> 1
		if x <= a[m] {
			ok = m
		} else {
			ng = m
		}
	}
	return ok
}

// Gtは昇順ソート済みの配列aに対して、xより大きい要素の左端Indexを返す
// aのすべての要素がx以下の場合、len(a)を返す
func Gt[T cmp.Ordered](a []T, x T) int {
	ok, ng := len(a), -1
	for ok-ng > 1 {
		m := (ok + ng) >> 1
		if x < a[m] {
			ok = m
		} else {
			ng = m
		}
	}
	return ok
}

// Leは昇順ソート済みの配列aに対し、x以下の要素の右端Indexを返す
// aのすべての要素がxより大きい場合、-1を返す
func Le[T cmp.Ordered](a []T, x T) int {
	ok, ng := -1, len(a)
	for ng-ok > 1 {
		m := (ok + ng) >> 1
		if x >= a[m] {
			ok = m
		} else {
			ng = m
		}
	}
	return ok
}

// Ltは昇順ソート済みの配列aに対して、xより小さい要素の右端を返す
// aのすべての要素がx以上の場合、-1を返す
func Lt[T cmp.Ordered](a []T, x T) int {
	ok, ng := -1, len(a)
	for ng-ok > 1 {
		m := (ok + ng) >> 1
		if x > a[m] {
			ok = m
		} else {
			ng = m
		}
	}
	return ok
}
