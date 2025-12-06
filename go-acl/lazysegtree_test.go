package main

import (
	"math/rand"
	"testing"
)

// ナイーブな実装（正解モデル）
type NaiveArray struct {
	data []int
}

func NewNaiveArray(arr []int) *NaiveArray {
	data := make([]int, len(arr))
	copy(data, arr)
	return &NaiveArray{data: data}
}

func (na *NaiveArray) RangeAdd(l, r, val int) {
	for i := l; i < r; i++ {
		na.data[i] += val
	}
}

func (na *NaiveArray) RangeSum(l, r int) int {
	sum := 0
	for i := l; i < r; i++ {
		sum += na.data[i]
	}
	return sum
}

func (na *NaiveArray) RangeMax(l, r int) int {
	const INF = 1 << 60
	max := -INF
	for i := l; i < r; i++ {
		if na.data[i] > max {
			max = na.data[i]
		}
	}
	return max
}

func (na *NaiveArray) Get(i int) int {
	return na.data[i]
}

// テスト: 区間加算・区間和
func TestLazySegTreeRangeAddRangeSum(t *testing.T) {
	rand.Seed(42)

	for trial := 0; trial < 100; trial++ {
		n := rand.Intn(50) + 1
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(100) - 50
		}

		naive := NewNaiveArray(data)
		lst := NewLazySegmentTree(data, LazyMoRangeAddRangeSum[int]())

		// ランダムなクエリを実行
		for q := 0; q < 100; q++ {
			l := rand.Intn(n)
			r := rand.Intn(n-l) + l + 1

			if rand.Intn(2) == 0 {
				// 区間加算
				val := rand.Intn(100) - 50
				naive.RangeAdd(l, r, val)
				lst.Apply(l, r, val)
			} else {
				// 区間和クエリ
				expected := naive.RangeSum(l, r)
				got := lst.Query(l, r)
				if expected != got {
					t.Fatalf("Trial %d, Query %d: RangeSum(%d, %d) = %d, want %d",
						trial, q, l, r, got, expected)
				}
			}
		}
	}
}

// テスト: 区間加算・区間最大値
func TestLazySegTreeRangeAddRangeMax(t *testing.T) {
	rand.Seed(43)

	for trial := 0; trial < 100; trial++ {
		n := rand.Intn(50) + 1
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(100) - 50
		}

		naive := NewNaiveArray(data)
		lst := NewLazySegmentTree(data, LazyMoRangeAddRangeMax())

		// ランダムなクエリを実行
		for q := 0; q < 100; q++ {
			l := rand.Intn(n)
			r := rand.Intn(n-l) + l + 1

			if rand.Intn(2) == 0 {
				// 区間加算
				val := rand.Intn(100) - 50
				naive.RangeAdd(l, r, val)
				lst.Apply(l, r, val)
			} else {
				// 区間最大値クエリ
				expected := naive.RangeMax(l, r)
				got := lst.Query(l, r)
				if expected != got {
					t.Fatalf("Trial %d, Query %d: RangeMax(%d, %d) = %d, want %d",
						trial, q, l, r, got, expected)
				}
			}
		}
	}
}

// テスト: 一点取得
func TestLazySegTreeGet(t *testing.T) {
	rand.Seed(44)

	for trial := 0; trial < 50; trial++ {
		n := rand.Intn(30) + 1
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(100)
		}

		naive := NewNaiveArray(data)
		lst := NewLazySegmentTree(data, LazyMoRangeAddRangeSum[int]())

		for q := 0; q < 50; q++ {
			l := rand.Intn(n)
			r := rand.Intn(n-l) + l + 1
			val := rand.Intn(100) - 50

			naive.RangeAdd(l, r, val)
			lst.Apply(l, r, val)

			// ランダムな位置をチェック
			i := rand.Intn(n)
			expected := naive.Get(i)
			got := lst.Get(i)
			if expected != got {
				t.Fatalf("Trial %d, Query %d: Get(%d) = %d, want %d",
					trial, q, i, got, expected)
			}
		}
	}
}

// テスト: エッジケース
func TestLazySegTreeEdgeCases(t *testing.T) {
	// サイズ1
	lst := NewLazySegmentTree([]int{10}, LazyMoRangeAddRangeSum[int]())
	lst.Apply(0, 1, 5)
	if got := lst.Query(0, 1); got != 15 {
		t.Errorf("Size 1: got %d, want 15", got)
	}

	// 全区間更新
	data := []int{1, 2, 3, 4, 5}
	lst = NewLazySegmentTree(data, LazyMoRangeAddRangeSum[int]())
	lst.Apply(0, 5, 10)
	if got := lst.Query(0, 5); got != 65 {
		t.Errorf("Full range: got %d, want 65", got)
	}

	// 空区間
	if got := lst.Query(2, 2); got != 0 {
		t.Errorf("Empty range: got %d, want 0", got)
	}
}
