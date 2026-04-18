package main

import "math/bits"

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

// Getは一点取得
func (t *LazySegmentTree[S, F]) Get(i int) S {
	return t.Query(i, i+1)
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
