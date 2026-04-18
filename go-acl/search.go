package main

// Enumは全探索してくれる関数型
// 全探索の要素を長さkの配列で引数とでて受け取る関数で、walkする
// 受け取る配列は参照なので、変更してはいけない。
//
// 5C3の全探索
//
//	enum := CombIdxEnum(5, 3)
//	enum(func(idx []int) bool {
//		// idxは選ばれたindex配列
//		...
//		return true
//	})
//
// A配列Aから3つ選んで合計が50未満になる選び方の件数
//
//	e := CombIdxEnum(30, 3)
//	A := make([]int, 30) // 何か値が入っているとする
//	c := e.Filter(func(idx []int) bool {
//		return A[idx[0]] + A[idx[1]] + A[idx[2]]
//	}).Count()
type Enum func(walk func([]int) bool)

// CombIdxEnum は、0...n-1個からk個選ぶ全探索を行える
//
// 5C3の全探索
//
//	enum := CombIdxEnum(5, 3)
func CombIdxEnum(n, k int) Enum {
	return func(walk func([]int) bool) {
		// 全探索の要素を保持する配列
		buf := make([]int, k)

		// pos: k個選ぶうちの、今0..k番目のどれか？
		// start: {pos}番目の要素の探索開始index
		// 5C3の場合(pos, start)
		// (0,0) -> [0,-,-]
		//   (1,1) -> [0,1,-]
		//     (2,2) -> [0,1,2]
		//       (3,3) -> @walk [0,1,2]@
		//     (2,2) -> [0,1,3]
		//       (3,4) -> @walk [0,1,3]@
		//     (2,2) -> [0,1,4]
		//       (3,5) -> @walk [0,1,4]@
		//   (1,1) -> [0,2,-]
		//     (2,2) -> [0,2,3]
		//       (3,4) -> @walk [0,2,3]@
		//     (2,2) -> [0,2,4]
		//       (3,5) -> @walk [0,2,4]@
		//   (1,1) -> [0,3,-]
		//     (2,3) -> [0,3,4]
		//       (3,5) -> @walk [0,3,4]@
		// (0,0) -> [1,-,-]
		//   (1,2) -> [1,2,-]
		//     (2,3) -> [1,2,3]
		//       (3,4) -> @walk [1,2,3]@
		//     (2,3) -> [1,2,4]
		//       (3,5) -> @walk [1,2,4]@
		//   (1,2) -> [1,3,-]
		//     (2,4) -> [1,3,4]
		//       (3,5) -> @walk [1,3,4]@
		// (0,0) -> [2,-,-]
		//   (1,3) -> [2,3,-]
		//     (2,4) -> [2,3,4]
		//       (3,5) -> @walk [2,3,4]@
		// 計10walk 30再帰
		var dfs func(pos, start int) bool
		dfs = func(pos, start int) bool {
			if pos == k {
				return walk(buf)
			}

			for i := start; i <= n-(k-pos); i++ {
				buf[pos] = i
				if !dfs(pos+1, i+1) {
					return false
				}
			}
			return true
		}
		dfs(0, 0)
	}
}

// PermIdxEnum は、0...n-1個からk個並べる並べ方の全探索をしてくれる
//
// 5P3の全探索
//
//	enum := PermIdxEnum(5, 3)
func PermIdxEnum(n, k int) Enum {
	return func(walk func([]int) bool) {
		buf := make([]int, k)
		used := make([]bool, n)
		var dfs func(pos int) bool
		dfs = func(pos int) bool {
			if pos == k {
				return walk(buf)
			}
			for i := 0; i < n; i++ {
				if used[i] {
					continue
				}
				used[i], buf[pos] = true, i
				if !dfs(pos + 1) {
					return false
				}
				used[i] = false
			}
			return true
		}
		dfs(0)
	}
}

// 実行
func (e Enum) Do(walk func([]int) bool) { e(walk) }

// Anyは全探索中にpred関数の条件に適う例が見つかればTrueを返す。
//
//	 例）[0...4]の中から3つ選んだとき、その和が15以上になるものが１つ以上存在するか
//
//		enum := CombIdxEnum(5,3)
//		enum.Any(func(idx []int) {
//			return Sum(idx) > 15
//		})
func (e Enum) Any(pred func([]int) bool) bool {
	hit := false
	e(func(idx []int) bool {
		if pred(idx) {
			hit = true
			return false
		}
		return true
	})
	return hit
}

func (e Enum) First(pred func([]int) bool) ([]int, bool) {
	var ans []int
	ok := false
	e(func(idx []int) bool {
		if pred(idx) {
			ans = append([]int(nil), idx...)
			ok = true
			return false
		}
		return true
	})
	return ans, ok
}

func (e Enum) Count() int {
	var c int
	e(func(_ []int) bool { c++; return true })
	return c
}

func (e Enum) Filter(pred func([]int) bool) Enum {
	return func(walk func([]int) bool) {
		e(func(idx []int) bool {
			if !pred(idx) {
				return true
			}
			return walk(idx)
		})
	}
}

func (e Enum) Take(k int) Enum {
	return func(walk func([]int) bool) {
		var c int
		e(func(idx []int) bool {
			if c >= k {
				return false
			}
			c++
			return walk(idx)
		})
	}
}
