package acl

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

func (na *NaiveArray) RangeMin(l, r int) int {
	const INF = 1 << 60
	min := INF
	for i := l; i < r; i++ {
		if na.data[i] < min {
			min = na.data[i]
		}
	}
	return min
}

func (na *NaiveArray) RangeSet(l, r, val int) {
	for i := l; i < r; i++ {
		na.data[i] = val
	}
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
			got := lst.At(i)
			if expected != got {
				t.Fatalf("Trial %d, Query %d: Get(%d) = %d, want %d",
					trial, q, i, got, expected)
			}
		}
	}
}

// テスト: 区間加算・区間最小値
func TestLazySegTreeRangeAddRangeMin(t *testing.T) {
	rand.Seed(45)

	for trial := 0; trial < 100; trial++ {
		n := rand.Intn(50) + 1
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(100) - 50
		}

		naive := NewNaiveArray(data)
		lst := NewLazySegmentTree(data, LazyMoRangeAddRangeMin())

		for q := 0; q < 100; q++ {
			l := rand.Intn(n)
			r := rand.Intn(n-l) + l + 1

			if rand.Intn(2) == 0 {
				val := rand.Intn(100) - 50
				naive.RangeAdd(l, r, val)
				lst.Apply(l, r, val)
			} else {
				expected := naive.RangeMin(l, r)
				got := lst.Query(l, r)
				if expected != got {
					t.Fatalf("Trial %d, Query %d: RangeMin(%d, %d) = %d, want %d",
						trial, q, l, r, got, expected)
				}
			}
		}
	}
}

// テスト: 区間更新・区間和
func TestLazySegTreeRangeUpdateRangeSum(t *testing.T) {
	rand.Seed(46)

	for trial := 0; trial < 100; trial++ {
		n := rand.Intn(50) + 1
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(100) - 50
		}

		naive := NewNaiveArray(data)
		lst := NewLazySegmentTree(data, LazyMoRangeUpdateRangeSum[int]())

		for q := 0; q < 100; q++ {
			l := rand.Intn(n)
			r := rand.Intn(n-l) + l + 1

			if rand.Intn(2) == 0 {
				val := rand.Intn(100) - 50
				naive.RangeSet(l, r, val)
				lst.Apply(l, r, &val)
			} else {
				expected := naive.RangeSum(l, r)
				got := lst.Query(l, r)
				if expected != got {
					t.Fatalf("Trial %d, Query %d: RangeUpdateSum(%d, %d) = %d, want %d",
						trial, q, l, r, got, expected)
				}
			}
		}
	}
}

// テスト: 区間更新・区間最大値
func TestLazySegTreeRangeUpdateRangeMax(t *testing.T) {
	rand.Seed(47)

	for trial := 0; trial < 100; trial++ {
		n := rand.Intn(50) + 1
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(100) - 50
		}

		naive := NewNaiveArray(data)
		lst := NewLazySegmentTree(data, LazyMoRangeUpdateRangeMax())

		for q := 0; q < 100; q++ {
			l := rand.Intn(n)
			r := rand.Intn(n-l) + l + 1

			if rand.Intn(2) == 0 {
				val := rand.Intn(100) - 50
				naive.RangeSet(l, r, val)
				lst.Apply(l, r, &val)
			} else {
				expected := naive.RangeMax(l, r)
				got := lst.Query(l, r)
				if expected != got {
					t.Fatalf("Trial %d, Query %d: RangeUpdateMax(%d, %d) = %d, want %d",
						trial, q, l, r, got, expected)
				}
			}
		}
	}
}

// テスト: 区間更新・区間最小値
func TestLazySegTreeRangeUpdateRangeMin(t *testing.T) {
	rand.Seed(48)

	for trial := 0; trial < 100; trial++ {
		n := rand.Intn(50) + 1
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(100) - 50
		}

		naive := NewNaiveArray(data)
		lst := NewLazySegmentTree(data, LazyMoRangeUpdateRangeMin())

		for q := 0; q < 100; q++ {
			l := rand.Intn(n)
			r := rand.Intn(n-l) + l + 1

			if rand.Intn(2) == 0 {
				val := rand.Intn(100) - 50
				naive.RangeSet(l, r, val)
				lst.Apply(l, r, &val)
			} else {
				expected := naive.RangeMin(l, r)
				got := lst.Query(l, r)
				if expected != got {
					t.Fatalf("Trial %d, Query %d: RangeUpdateMin(%d, %d) = %d, want %d",
						trial, q, l, r, got, expected)
				}
			}
		}
	}
}

// テスト: Set()一点更新
func TestLazySegTreeSet(t *testing.T) {
	rand.Seed(49)

	for trial := 0; trial < 50; trial++ {
		n := rand.Intn(30) + 1
		data := make([]int, n)
		for i := range data {
			data[i] = rand.Intn(100)
		}

		naive := NewNaiveArray(data)
		lst := NewLazySegmentTree(data, LazyMoRangeAddRangeSum[int]())

		for q := 0; q < 50; q++ {
			i := rand.Intn(n)
			val := rand.Intn(100)

			naive.data[i] = val
			lst.Set(i, val)

			// クエリで検証
			l := rand.Intn(n)
			r := rand.Intn(n-l) + l + 1
			expected := naive.RangeSum(l, r)
			got := lst.Query(l, r)
			if expected != got {
				t.Fatalf("Trial %d, Query %d: after Set(%d,%d), RangeSum(%d,%d) = %d, want %d",
					trial, q, i, val, l, r, got, expected)
			}
		}
	}
}

// テスト: MaxRight/MinLeft
func TestLazySegTreeMaxRightMinLeft(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	lst := NewLazySegmentTree(data, LazyMoRangeAddRangeSum[int]())

	for l := 0; l <= 5; l++ {
		for limit := 0; limit <= 15; limit++ {
			got := lst.MaxRight(l, func(v int) bool { return v <= limit })
			want := l
			sum := 0
			for r := l; r < 5; r++ {
				sum += data[r]
				if sum <= limit {
					want = r + 1
				} else {
					break
				}
			}
			if got != want {
				t.Errorf("MaxRight(%d, <=%d) = %d, want %d", l, limit, got, want)
			}
		}
	}

	for r := 0; r <= 5; r++ {
		for limit := 0; limit <= 15; limit++ {
			got := lst.MinLeft(r, func(v int) bool { return v <= limit })
			want := r
			sum := 0
			for l := r - 1; l >= 0; l-- {
				sum += data[l]
				if sum <= limit {
					want = l
				} else {
					break
				}
			}
			if got != want {
				t.Errorf("MinLeft(%d, <=%d) = %d, want %d", r, limit, got, want)
			}
		}
	}
}

// テスト: Apply後のMaxRight (遅延伝播の確認)
func TestLazySegTreeMaxRightAfterApply(t *testing.T) {
	data := []int{1, 1, 1, 1, 1}
	lst := NewLazySegmentTree(data, LazyMoRangeAddRangeSum[int]())
	lst.Apply(0, 5, 2) // 全部3に: [3,3,3,3,3]

	got := lst.MaxRight(0, func(v int) bool { return v <= 9 })
	if got != 3 {
		t.Errorf("MaxRight after apply: got %d, want 3", got)
	}

	got = lst.MinLeft(5, func(v int) bool { return v <= 9 })
	if got != 2 {
		t.Errorf("MinLeft after apply: got %d, want 2", got)
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
