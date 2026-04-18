package main

import (
	"cmp"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// L1 は長さnの配列を生成する
func L1[T any](n int) []T {
	return make([]T, n)
}

// L2 はn1行n2列の配列を生成する
func L2[T any](n1, n2 int) [][]T {
	ret := make([][]T, n1)
	for i := range ret {
		ret[i] = make([]T, n2)
	}
	return ret
}

func L3[T any](n1, n2, n3 int) [][][]T {
	ret := make([][][]T, n1)
	for i := range ret {
		ret[i] = make([][]T, n2)
		for j := range ret[i] {
			ret[i][j] = make([]T, n3)
		}
	}
	return ret
}

func F1[T any](vals []T, fill T) {
	for i := range vals {
		vals[i] = fill
	}
}

func F2[T any](vals [][]T, fill T) {
	for i := range vals {
		for j := range vals[i] {
			vals[i][j] = fill
		}
	}
}

func F3[T any](vals [][][]T, fill T) {
	for i := range vals {
		for j := range vals[i] {
			for k := range vals[i][j] {
				vals[i][j][k] = fill
			}
		}
	}
}

// Jag は長さnのJagged配列を生成する
func Jag[T any](n int) [][]T {
	ret := make([][]T, n)
	for i := range ret {
		ret[i] = make([]T, 0, 10)
	}
	return ret
}

func Reverse[T any](vals []T) []T {
	ret := make([]T, len(vals))
	for i := len(vals) - 1; i >= 0; i-- {
		ret[i] = vals[len(vals)-1-i]
	}
	return ret
}

func ReverseS(s string) string {
	return string(Reverse([]byte(s)))
}

func RotateCW90[T any](src [][]T) [][]T {
	if len(src) == 0 || len(src[0]) == 0 {
		return nil
	}
	h, w := len(src), len(src[0])

	ret := make([][]T, w)
	for i := range ret {
		ret[i] = make([]T, h)
	}
	for ih := range src {
		for iw := range src[ih] {
			jh, jw := iw, h-1-ih
			ret[jh][jw] = src[ih][iw]
		}
	}
	return ret
}

func Uniq[T cmp.Ordered](vals []T) []T {
	ret := make([]T, 0, len(vals))
	if len(vals) == 0 {
		return ret
	}
	ret = append(ret, vals[0])
	for _, v := range vals[1:] {
		if ret[len(ret)-1] != v {
			ret = append(ret, v)
		}
	}
	return ret
}

func Sort[T cmp.Ordered](arr []T) []T {
	ret := append([]T(nil), arr...)
	sort.Slice(ret, func(i int, j int) bool { return ret[i] < ret[j] })
	return ret
}

func SortIdx[T cmp.Ordered](arr []T) []int {
	ret := make([]int, len(arr))
	for i := range ret {
		ret[i] = i
	}
	sort.SliceStable(ret, func(i, j int) bool { return arr[ret[i]] < arr[ret[j]] })
	return ret
}

func SortF[T any](arr []T, less LessFunc[T]) []T {
	ret := append([]T(nil), arr...)
	sort.Slice(ret, func(i, j int) bool { return less(ret[i], ret[j]) })
	return ret
}

func SortIdxF[T any](arr []T, less LessFunc[T]) []int {
	ret := make([]int, len(arr))
	for i := range ret {
		ret[i] = i
	}
	sort.Slice(ret, func(i, j int) bool { return less(arr[ret[i]], arr[ret[j]]) })
	return ret
}

func SortE[K cmp.Ordered, V any](arr []Entry[K, V]) []Entry[K, V] {
	ret := append([]Entry[K, V](nil), arr...)
	sort.Slice(ret, func(i, j int) bool { return ret[i].K < ret[j].K })
	return ret
}

func SortIdxE[K cmp.Ordered, V any](arr []Entry[K, V]) []int {
	ret := make([]int, len(arr))
	for i := range ret {
		ret[i] = i
	}
	sort.Slice(ret, func(i, j int) bool { return arr[ret[i]].K < arr[ret[j]].K })
	return ret
}

type Key string

func KeyInts(a []int) Key {
	if len(a) == 0 {
		return ""
	}
	var b strings.Builder

	b.Grow(len(a) * 3)
	b.WriteString(strconv.Itoa(a[0]))
	for i := 1; i < len(a); i++ {
		b.WriteByte(' ')
		b.WriteString(strconv.Itoa(a[i]))
	}
	return Key(b.String())
}

func (k Key) ToInts() []int {
	toks := strings.Fields(string(k))
	ret := make([]int, 0, len(toks))
	for _, s := range toks {
		x, err := strconv.Atoi(s)
		if err != nil {
			panic(fmt.Errorf("failed to parse int %s: element of %v %w", s, k, err))
		}
		ret = append(ret, x)
	}
	return ret
}
