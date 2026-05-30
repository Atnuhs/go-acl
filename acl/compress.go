package acl

import (
	"slices"
	"cmp"
)

// Compress は座標圧縮を行うデータ構造。
// 値域が大きい列を、出現する要素を昇順に並べたときのインデックス (0..k-1) に対応付ける。
// 大きな値域を配列の添字として扱いたい場面 (BIT、SegmentTree、imos 法など) で使う。
type Compress[T cmp.Ordered] struct {
	toOrig []T       // idx -> value
	toIdx  map[T]int // value -> idx
}

// NewCompress は vals を座標圧縮し、対応表を構築する。
// 重複は自動的に取り除かれ、昇順で 0 始まりのインデックスが割り当てられる。
// 入力スライス vals は変更しない。
//
// 計算量: 時間 O(n log n)、空間 O(n) (n = len(vals))
func NewCompress[T cmp.Ordered](vals []T) *Compress[T] {
	v := make([]T, len(vals))
	copy(v, vals)
	slices.Sort(v)
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

// Idx は値 x に対応する圧縮後インデックスを返す。
// 構築時の vals に x が含まれていなければ -1 を返す。
//
// 計算量: 平均 O(1) (map アクセス)
func (c *Compress[T]) Idx(x T) int {
	if i, ok := c.toIdx[x]; ok {
		return i
	}
	return -1
}

// Val は圧縮後インデックス i に対応する元の値を返す。
// i が [0, Size()) の範囲外の場合は panic する。
//
// 計算量: O(1)
func (c *Compress[T]) Val(i int) T {
	return c.toOrig[i]
}

// Size は圧縮後の異なる値の個数 (= 有効なインデックス数) を返す。
//
// 計算量: O(1)
func (c *Compress[T]) Size() int {
	return len(c.toOrig)
}
