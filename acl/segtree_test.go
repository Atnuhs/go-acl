package acl

import (
	"fmt"
	"testing"

	"github.com/Atnuhs/go-acl/acl/testlib"
)

func TestMoMax(t *testing.T) {
	mo := MoMax()

	testlib.AclAssert(t, -INF, mo.E)
	testlib.AclAssert(t, 5, mo.Op(3, 5))
	testlib.AclAssert(t, 5, mo.Op(5, 3))
	testlib.AclAssert(t, 0, mo.Op(-1, 0))
}

func TestMoMin(t *testing.T) {
	mo := MoMin()

	testlib.AclAssert(t, INF, mo.E)
	testlib.AclAssert(t, 3, mo.Op(3, 5))
	testlib.AclAssert(t, 3, mo.Op(5, 3))
	testlib.AclAssert(t, -1, mo.Op(-1, 0))
}

func TestMoSum(t *testing.T) {
	mo := MoSum[int]()

	testlib.AclAssert(t, 0, mo.E)
	testlib.AclAssert(t, 8, mo.Op(3, 5))
	testlib.AclAssert(t, -1, mo.Op(-5, 4))
}

func TestMoXOR(t *testing.T) {
	mo := MoXOR()

	testlib.AclAssert(t, 0, mo.E)
	testlib.AclAssert(t, 6, mo.Op(3, 5)) // 3 ^ 5 = 6
	testlib.AclAssert(t, 0, mo.Op(5, 5)) // 5 ^ 5 = 0
}

func TestNewSegmentTree(t *testing.T) {
	tests := map[string]struct {
		arr      []int
		monoid   *Monoid[int]
		testFunc func(*testing.T, *SegmentTree[int])
	}{
		"sum tree": {
			arr:    []int{1, 2, 3, 4, 5},
			monoid: MoSum[int](),
			testFunc: func(t *testing.T, st *SegmentTree[int]) {
				t.Log(st.dump())
				testlib.AclAssert(t, 15, st.Query(0, 5)) // Sum of all
				t.Log(st.dump())
				testlib.AclAssert(t, 5, st.Query(1, 3)) // Sum of [2, 3]
				t.Log(st.dump())
				testlib.AclAssert(t, 1, st.Query(0, 1)) // Sum of [1]
				t.Log(st.dump())
			},
		},
		"max tree": {
			arr:    []int{3, 1, 4, 1, 5},
			monoid: MoMax(),
			testFunc: func(t *testing.T, st *SegmentTree[int]) {
				testlib.AclAssert(t, 5, st.Query(0, 5)) // Max of all
				testlib.AclAssert(t, 4, st.Query(1, 3)) // Max of [1, 4]
				testlib.AclAssert(t, 3, st.Query(0, 1)) // Max of [3]
			},
		},
		"min tree": {
			arr:    []int{3, 1, 4, 1, 5},
			monoid: MoMin(),
			testFunc: func(t *testing.T, st *SegmentTree[int]) {
				testlib.AclAssert(t, 1, st.Query(0, 5)) // Min of all
				testlib.AclAssert(t, 1, st.Query(1, 3)) // Min of [1, 4]
				testlib.AclAssert(t, 3, st.Query(0, 1)) // Min of [3]
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			st := NewSegmentTree(tc.arr, tc.monoid)
			tc.testFunc(t, st)
		})
	}
}

func TestSegmentTree_Update(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	st := NewSegmentTree(arr, MoSum[int]())

	// Initial sum
	testlib.AclAssert(t, 15, st.Query(0, 5))

	// Update index 2 from 3 to 10
	st.Set(2, 10)
	testlib.AclAssert(t, 22, st.Query(0, 5)) // 1+2+10+4+5 = 22
	testlib.AclAssert(t, 12, st.Query(1, 3)) // 2+10 = 12

	// Update index 0 from 1 to 0
	st.Set(0, 0)
	testlib.AclAssert(t, 21, st.Query(0, 5)) // 0+2+10+4+5 = 21
}

func TestSegmentTree_QueryEdgeCases(t *testing.T) {
	arr := []int{1, 2, 3}
	st := NewSegmentTree(arr, MoSum[int]())

	// Empty range should return identity element
	testlib.AclAssert(t, 0, st.Query(0, 0))
	testlib.AclAssert(t, 0, st.Query(1, 1))
	testlib.AclAssert(t, 0, st.Query(3, 3))

	// Single element ranges
	testlib.AclAssert(t, 1, st.Query(0, 1))
	testlib.AclAssert(t, 2, st.Query(1, 2))
	testlib.AclAssert(t, 3, st.Query(2, 3))
}

func TestMoMODMul(t *testing.T) {
	mo := MoMODMul(7)

	testlib.AclAssert(t, 1, mo.E)
	testlib.AclAssert(t, 6, mo.Op(2, 3))   // 2*3=6
	testlib.AclAssert(t, 1, mo.Op(3, 5))   // 15 % 7 = 1
	testlib.AclAssert(t, 0, mo.Op(7, 100)) // 0 % 7 = 0
}

func TestSegmentTree_At(t *testing.T) {
	arr := []int{10, 20, 30, 40, 50}
	st := NewSegmentTree(arr, MoSum[int]())

	for i, want := range arr {
		if got := st.At(i); got != want {
			t.Errorf("At(%d) = %d, want %d", i, got, want)
		}
	}

	st.Set(2, 99)
	testlib.AclAssert(t, 99, st.At(2))
	testlib.AclAssert(t, 10, st.At(0))
}

func TestMaxRight_Max(t *testing.T) {
	seg := NewSegmentTree([]int{3, 1, 4, 1, 5, 9, 2}, MoMax())

	tests := []struct {
		name  string
		l     int
		limit int
		want  int
	}{
		{"max <= 4 from 0", 0, 4, 4},     // max(3,1,4)=4
		{"max <= 3 from 0", 0, 3, 2},     // max(3,1)=3
		{"max <= 1 from 1", 1, 1, 2},     // max(1,1,1)=1
		{"max <= 100 from 0", 0, 100, 7}, // 全体
		{"max <= 5 from 3", 3, 5, 5},     // max(1,5)=5
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := seg.MaxRight(tt.l, func(v int) bool {
				return v <= tt.limit
			})
			t.Logf("\n%s", seg.dump())
			if got != tt.want {
				t.Errorf("MaxRight(%d, max <= %d) = %d, want %d",
					tt.l, tt.limit, got, tt.want)
			}
		})
	}
}

func TestMaxRight_Panic(t *testing.T) {
	seg := NewSegmentTree([]int{1, 2, 3}, MoSum[int]())

	t.Run("negative index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative index")
			}
		}()
		seg.MaxRight(-1, func(v int) bool { return true })
	})

	t.Run("out of bounds index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for out of bounds index")
			}
		}()
		seg.MaxRight(4, func(v int) bool { return true })
	})

	t.Run("l == n returns n", func(t *testing.T) {
		got := seg.MaxRight(3, func(v int) bool { return true })
		if got != 3 {
			t.Errorf("MaxRight(3) = %d, want 3", got)
		}
	})
}

func TestMaxRight_VerifyWithQuery(t *testing.T) {
	// MaxRightの結果が正しいか、Queryを使って検証
	seg := NewSegmentTree([]int{1, 2, 3, 4, 5, 6, 7, 8}, MoSum[int]())

	for l := 0; l < 8; l++ {
		for limit := 0; limit <= 40; limit += 5 {
			r := seg.MaxRight(l, func(v int) bool { return v <= limit })

			// rが正しい位置にあることを確認
			if r > l {
				sum := seg.Query(l, r)
				if sum > limit {
					t.Errorf("MaxRight(%d, <=%d)=%d but Query(%d,%d)=%d > %d",
						l, limit, r, l, r, sum, limit)
				}
			}

			// r+1が条件を満たさないことを確認（範囲内の場合）
			if r < 8 {
				sum := seg.Query(l, r+1)
				if sum <= limit {
					t.Errorf("MaxRight(%d, <=%d)=%d but Query(%d,%d)=%d <= %d (should be larger)",
						l, limit, r, l, r+1, sum, limit)
				}
			}
		}
	}
}

func TestMaxRight_LargeArray(t *testing.T) {
	// 大きな配列でのテスト
	n := 1000
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = 1
	}
	seg := NewSegmentTree(data, MoSum[int]())

	tests := []struct {
		l     int
		limit int
		want  int
	}{
		{0, 100, 100},
		{0, 500, 500},
		{100, 200, 300},
		{999, 1, 1000},
	}

	for _, tt := range tests {
		got := seg.MaxRight(tt.l, func(v int) bool { return v <= tt.limit })
		if got != tt.want {
			t.Errorf("MaxRight(%d, <=%d) = %d, want %d",
				tt.l, tt.limit, got, tt.want)
		}
	}
}

func TestMaxRight_Detailed(t *testing.T) {
	// すべての l と様々な閾値で徹底的にテスト
	data := []int{1, 2, 3, 4, 5}
	seg := NewSegmentTree(data, MoSum[int]())

	// 期待値を愚直に計算
	for l := 0; l < 5; l++ {
		for limit := 0; limit <= 15; limit++ {
			got := seg.MaxRight(l, func(v int) bool { return v <= limit })

			// 愚直に正解を求める
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
				t.Errorf("data=%v, MaxRight(%d, <=%d) = %d, want %d",
					data, l, limit, got, want)
			}
		}
	}
}

func TestMinLeft_Detailed(t *testing.T) {
	data := []int{1, 2, 3, 4, 5}
	seg := NewSegmentTree(data, MoSum[int]())

	for r := 0; r <= 5; r++ {
		for limit := 0; limit <= 15; limit++ {
			got := seg.MinLeft(r, func(v int) bool { return v <= limit })

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
				t.Errorf("data=%v, MinLeft(%d, <=%d) = %d, want %d",
					data, r, limit, got, want)
			}
		}
	}
}

func TestMinLeft_Panic(t *testing.T) {
	seg := NewSegmentTree([]int{1, 2, 3}, MoSum[int]())

	t.Run("negative index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for negative index")
			}
		}()
		seg.MinLeft(-1, func(v int) bool { return true })
	})

	t.Run("out of bounds index", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for out of bounds index")
			}
		}()
		seg.MinLeft(4, func(v int) bool { return true })
	})

	t.Run("r == 0 returns 0", func(t *testing.T) {
		got := seg.MinLeft(0, func(v int) bool { return true })
		if got != 0 {
			t.Errorf("MinLeft(0) = %d, want 0", got)
		}
	})
}

// 基本的な使い方: 区間和モノイドで一点更新と区間和クエリを行う。
func ExampleSegmentTree() {
	arr := []int{1, 2, 3, 4, 5}
	seg := NewSegmentTree(arr, MoSum[int]())

	fmt.Println(seg.Query(0, 5)) // 全体の和
	fmt.Println(seg.Query(1, 4)) // [1, 4) の和: 2 + 3 + 4

	seg.Set(2, 100) // a[2] を 100 に更新
	fmt.Println(seg.Query(0, 5))
	fmt.Println(seg.At(2))
	// Output:
	// 15
	// 9
	// 112
	// 100
}

// モノイドを差し替えれば区間 max / min / xor などにも使える。
func ExampleSegmentTree_max() {
	seg := NewSegmentTree([]int{3, 1, 4, 1, 5, 9, 2, 6}, MoMax())

	fmt.Println(seg.Query(0, 8)) // 全体の最大値
	fmt.Println(seg.Query(0, 4)) // [0, 4) の最大値: max(3, 1, 4, 1)
	// Output:
	// 9
	// 4
}

// MaxRight はモノイド上の二分探索を O(log n) で行う。
// 区間和が初めて 10 を超える直前 (= 累積和がまだ 10 以下のままの最大の r) を求める例。
func ExampleSegmentTree_MaxRight() {
	seg := NewSegmentTree([]int{1, 2, 3, 4, 5}, MoSum[int]()) // 累積和: 1, 3, 6, 10, 15
	r := seg.MaxRight(0, func(s int) bool { return s <= 10 })
	fmt.Println(r) // Query(0, 4) = 10 はまだ ok。Query(0, 5) = 15 で初めて NG
	// Output:
	// 4
}
