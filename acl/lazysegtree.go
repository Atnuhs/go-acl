package acl

import (
	"fmt"
	"math/bits"
)

type LazySegmentTree[S, F any] struct {
	n      int
	size   int
	rank   int
	data   []S
	lazy   []F
	length []int
	lm     *LazyMonoid[S, F]
}

func NewLazySegmentTree[S, F any](arr []S, lm *LazyMonoid[S, F]) *LazySegmentTree[S, F] {
	n := len(arr)
	log := bits.Len(uint(n - 1))
	size := 1 << uint(log)

	data := L1[S](size << 1)
	lazy := L1[F](size)

	// 単位元で初期化
	F1(data, lm.E)
	F1(lazy, lm.Id)

	for i := range arr {
		data[size+i] = arr[i]
	}

	// その区間が請け負う範囲
	length := L1[int](size << 1)
	for i := range length[size:] {
		length[i+size] = 1
	}
	for i := size - 1; i >= 1; i-- {
		length[i] = length[i<<1] + length[i<<1|1]
	}

	t := &LazySegmentTree[S, F]{
		n:      n,
		size:   size,
		rank:   log,
		data:   data,
		lazy:   lazy,
		length: length,
		lm:     lm,
	}
	for i := size - 1; i >= 1; i-- {
		t.pull(i)
	}
	return t
}

func (t *LazySegmentTree[S, F]) Size() int {
	return t.n
}

// applyはノードkの更新と、子に適用予定の作用を遅延に蓄積
func (t *LazySegmentTree[S, F]) apply(k int, f F) {
	// dataへ直接適用
	t.data[k] = t.lm.Mapping(f, t.data[k], t.length[k])

	// 葉でない場合、lazyへ蓄積
	if k < t.size {
		t.lazy[k] = t.lm.Composition(f, t.lazy[k])
	}
}

// applyRangeは[l, r)の範囲を張る要素全体にapply
func (t *LazySegmentTree[S, F]) applyRange(l, r int, f F) {
	for l < r {
		if l&1 == 1 {
			t.apply(l, f)
			l++
		}
		if r&1 == 1 {
			r--
			t.apply(r, f)
		}
		l >>= 1
		r >>= 1
	}
}

// pullは値を子から更新
func (t *LazySegmentTree[S, F]) pull(i int) {
	t.data[i] = t.lm.Op(t.data[i<<1], t.data[i<<1|1])
}

func (t *LazySegmentTree[S, F]) pullRange(l, r int) {
	// 更新した範囲の親を更新
	for i := 1; i <= t.rank; i++ {
		if ((l >> i) << i) != l {
			t.pull(l >> i)
		}
		if ((r >> i) << i) != r {
			t.pull(r >> i)
		}
	}
}

// pushはノードkの遅延を子へ解消する
func (t *LazySegmentTree[S, F]) push(k int) {
	if k >= t.size {
		return
	}
	t.apply(k<<1, t.lazy[k])
	t.apply(k<<1|1, t.lazy[k])
	t.lazy[k] = t.lm.Id
}

// pushRangeは[l,r)の範囲で、apply予定の屋根より上で遅延を解消しきる
func (t *LazySegmentTree[S, F]) pushRange(l, r int) {
	// 必要な範囲の遅延を伝播
	for i := t.rank; i >= 1; i-- {
		// ノードlがノード(l>>i)の左端の葉ではない
		if ((l >> i) << i) != l {
			t.push(l >> i)
		}
		if ((r >> i) << i) != r {
			t.push(r >> i)
		}
	}
}

// Applyは区間[l, r)に作用素fを適用
func (t *LazySegmentTree[S, F]) Apply(l, r int, f F) {
	if l >= r {
		return
	}
	l += t.size
	r += t.size

	// 貯めている遅延をapply予定のノードへ下ろす
	t.pushRange(l, r)

	// applyする
	t.applyRange(l, r, f)

	// 正規の値を根に向かってpull upする
	t.pullRange(l, r)
}

func (t *LazySegmentTree[S, F]) Query(l, r int) S {
	if l >= r {
		return t.lm.E
	}
	l += t.size
	r += t.size

	// 貯めている遅延をapply予定のノードへ下ろす
	t.pushRange(l, r)

	sml, smr := t.lm.E, t.lm.E
	for l < r {
		if l&1 == 1 {
			sml = t.lm.Op(sml, t.data[l])
			l++
		}
		if r&1 == 1 {
			r--
			smr = t.lm.Op(t.data[r], smr)
		}
		l >>= 1
		r >>= 1
	}
	return t.lm.Op(sml, smr)
}

// Atは一点取得
func (t *LazySegmentTree[S, F]) At(i int) S {
	return t.Query(i, i+1)
}

func (t *LazySegmentTree[S, F]) ToSlice() []S {
	ret := L1[S](t.Size())
	for i := range t.Size() {
		ret[i] = t.At(i)
	}
	return ret
}

// Setは一点更新
func (t *LazySegmentTree[S, F]) Set(i int, x S) {
	i += t.size
	// 遅延のpush
	for j := t.rank; j >= 1; j-- {
		t.push(i >> j)
	}
	t.data[i] = x
	for i >>= 1; i >= 1; i >>= 1 {
		t.pull(i)
	}
}

// MaxRightはok(Query(l, r))がTrueとなるような最大のrを返す。
// ok(E)がTrueであり、okは単調 (True → False) であることを要求する。
// 0 <= l <= Size() を許容し、l == Size() のとき Size() を返す。
func (t *LazySegmentTree[S, F]) MaxRight(l int, ok func(S) bool) int {
	if l < 0 || l > t.n {
		panic(fmt.Errorf("MaxRight: l must be 0 <= l <= %d but got %d", t.n, l))
	}
	if l == t.n {
		return t.n
	}
	l += t.size
	for i := t.rank; i >= 1; i-- {
		t.push(l >> i)
	}
	sm := t.lm.E
	for {
		for l&1 == 0 {
			l >>= 1
		}
		if !ok(t.lm.Op(sm, t.data[l])) {
			for l < t.size {
				t.push(l)
				l <<= 1
				if v := t.lm.Op(sm, t.data[l]); ok(v) {
					sm = v
					l++
				}
			}
			return l - t.size
		}
		sm = t.lm.Op(sm, t.data[l])
		l++
		if l&-l == l {
			break
		}
	}
	return t.n
}

// MinLeftはok(Query(l, r))がTrueとなるような最小のlを返す。
// ok(E)がTrueであり、okは単調 (l を減らす方向で True → False) であることを要求する。
// 0 <= r <= Size() を許容し、r == 0 のとき 0 を返す。
func (t *LazySegmentTree[S, F]) MinLeft(r int, ok func(S) bool) int {
	if r < 0 || r > t.n {
		panic(fmt.Errorf("MinLeft: r must be 0 <= r <= %d but got %d", t.n, r))
	}
	if r == 0 {
		return 0
	}
	r += t.size
	for i := t.rank; i >= 1; i-- {
		t.push((r - 1) >> i)
	}
	sm := t.lm.E
	for {
		r--
		for r > 1 && (r&1) == 1 {
			r >>= 1
		}
		if !ok(t.lm.Op(t.data[r], sm)) {
			for r < t.size {
				t.push(r)
				r = r<<1 | 1
				if v := t.lm.Op(t.data[r], sm); ok(v) {
					sm = v
					r--
				}
			}
			return r + 1 - t.size
		}
		sm = t.lm.Op(t.data[r], sm)
		if r&-r == r {
			break
		}
	}
	return 0
}
