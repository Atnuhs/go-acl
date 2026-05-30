package acl

import (
	"cmp"
	"fmt"
	"slices"
	"strconv"
	"strings"
)

// L1 は長さ n の 1 次元配列を生成する。要素はゼロ値で初期化される。
//
// 計算量: O(n)
func L1[T any](n int) []T {
	return make([]T, n)
}

// L2 は n1 行 n2 列の 2 次元配列を生成する。要素はゼロ値で初期化される。
//
// 計算量: O(n1 * n2)
func L2[T any](n1, n2 int) [][]T {
	ret := make([][]T, n1)
	for i := range ret {
		ret[i] = make([]T, n2)
	}
	return ret
}

// L3 は n1 * n2 * n3 の 3 次元配列を生成する。要素はゼロ値で初期化される。
//
// 計算量: O(n1 * n2 * n3)
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

// F1 は 1 次元配列 vals の全要素を fill で埋める。
//
// 計算量: O(len(vals))
func F1[T any](vals []T, fill T) {
	for i := range vals {
		vals[i] = fill
	}
}

// F2 は 2 次元配列 vals の全要素を fill で埋める。
//
// 計算量: O(全要素数)
func F2[T any](vals [][]T, fill T) {
	for i := range vals {
		for j := range vals[i] {
			vals[i][j] = fill
		}
	}
}

// F3 は 3 次元配列 vals の全要素を fill で埋める。
//
// 計算量: O(全要素数)
func F3[T any](vals [][][]T, fill T) {
	for i := range vals {
		for j := range vals[i] {
			for k := range vals[i][j] {
				vals[i][j][k] = fill
			}
		}
	}
}

// Jag は長さ n の Jagged 配列 (各行が独立した可変長スライス) を生成する。
// 各行は空のスライスで、初期容量 10 で確保される。
//
// 計算量: O(n)
func Jag[T any](n int) [][]T {
	ret := make([][]T, n)
	for i := range ret {
		ret[i] = make([]T, 0, 10)
	}
	return ret
}

// Reverse は vals の要素を逆順にした新しいスライスを返す。元のスライスは変更しない。
//
// 計算量: O(len(vals))
func Reverse[T any](vals []T) []T {
	ret := make([]T, len(vals))
	for i := len(vals) - 1; i >= 0; i-- {
		ret[i] = vals[len(vals)-1-i]
	}
	return ret
}

// ReverseS は文字列 s を逆順にした文字列を返す。
// マルチバイト文字 (rune) ではなくバイト単位で反転するため、ASCII 用途を想定している。
//
// 計算量: O(len(s))
func ReverseS(s string) string {
	return string(Reverse([]byte(s)))
}

// RotateCW90 は 2 次元配列 src を時計回りに 90 度回転した新しい配列を返す。
// src が空、または最初の行が空の場合は nil を返す。
//
// 計算量: O(h * w) (h, w は src の行数・列数)
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

// Uniq は隣接する重複要素を取り除いた新しいスライスを返す。
// vals は昇順 (または降順) にソート済みであることを前提とする。
// 未ソートのスライスに対しては、隣接する重複だけが除去される点に注意。
//
// 計算量: O(len(vals))
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

// Sort は arr を昇順にソートした新しいスライスを返す。元のスライスは変更しない。
//
// 計算量: O(n log n)
func Sort[T cmp.Ordered](arr []T) []T {
	ret := slices.Clone(arr)
	slices.Sort(ret)
	return ret
}

// SortIdx は arr を昇順に並べたときのインデックス列を返す。
// すなわち、戻り値 idx は arr[idx[0]] <= arr[idx[1]] <= ... となるような順列。
// ソートは安定で、元の arr は変更しない。
//
// 計算量: O(n log n)
func SortIdx[T cmp.Ordered](arr []T) []int {
	ret := make([]int, len(arr))
	for i := range ret {
		ret[i] = i
	}
	slices.SortStableFunc(ret, func(i, j int) int { return cmp.Compare(arr[i], arr[j]) })
	return ret
}

// SortF は less の順序で arr をソートした新しいスライスを返す。
// ソートは安定で、元の arr は変更しない。
//
// 計算量: O(n log n)
func SortF[T any](arr []T, less LessFunc[T]) []T {
	ret := slices.Clone(arr)
	slices.SortStableFunc(ret, func(a, b T) int {
		switch {
		case less(a, b):
			return -1
		case less(b, a):
			return 1
		default:
			return 0
		}
	})
	return ret
}

// SortIdxF は less の順序で arr を並べたときのインデックス列を返す。
// ソートは安定で、元の arr は変更しない。
//
// 計算量: O(n log n)
func SortIdxF[T any](arr []T, less LessFunc[T]) []int {
	ret := make([]int, len(arr))
	for i := range ret {
		ret[i] = i
	}
	slices.SortStableFunc(ret, func(i, j int) int {
		switch {
		case less(arr[i], arr[j]):
			return -1
		case less(arr[j], arr[i]):
			return 1
		default:
			return 0
		}
	})
	return ret
}

// SortE は Entry のスライスをキー K の昇順でソートした新しいスライスを返す。
// ソートは安定で、元の arr は変更しない。
//
// 計算量: O(n log n)
func SortE[K cmp.Ordered, V any](arr []Entry[K, V]) []Entry[K, V] {
	ret := slices.Clone(arr)
	slices.SortStableFunc(ret, func(a, b Entry[K, V]) int { return cmp.Compare(a.K, b.K) })
	return ret
}

// SortIdxE は Entry のスライスをキー K の昇順で並べたときのインデックス列を返す。
// ソートは安定で、元の arr は変更しない。
//
// 計算量: O(n log n)
func SortIdxE[K cmp.Ordered, V any](arr []Entry[K, V]) []int {
	ret := make([]int, len(arr))
	for i := range ret {
		ret[i] = i
	}
	slices.SortStableFunc(ret, func(i, j int) int { return cmp.Compare(arr[i].K, arr[j].K) })
	return ret
}

// Key は整数列をハッシュ化して map のキーとして使うための文字列型。
// 配列・スライスは Go の map キーに使えないため、KeyInts で文字列化して利用する。
type Key string

// KeyInts は整数スライス a を空白区切りの Key に変換する。
// 同じ要素列からは同じ Key が得られるため、map のキーとして配列を扱える。
//
// 計算量: O(len(a))
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

// ToInts は KeyInts で生成された Key を整数スライスに復元する。
// Key が KeyInts 由来でないなど、整数として解釈できない要素を含む場合は panic する。
//
// 計算量: O(len(k))
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
