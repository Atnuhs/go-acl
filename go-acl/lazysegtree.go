package main

import "math/bits"

type (
	// 作用素同士の合成関数
	Composition[F any] func(f2, f1 F) F

	// 作用素モノイド
	Action[F any] struct {
		// 作用素のモノイド
		Composition Composition[F] // 作用素同士の合成
		Id          F              // 作用素の単位元(恒等写像)
	}
	Mapping[S, F any]    func(f F, x S, length int) S
	LazyMonoid[S, F any] struct {
		*Monoid[S]
		*Action[F]
		// 作用の定義
		Mapping Mapping[S, F] // 作用素fをデータxに適用
	}
)

func ActRangeAdd[T int | float64]() *Action[T] {
	return &Action[T]{
		Composition: func(f2, f1 T) T { return f1 + f2 },
		Id:          0,
	}
}

func ActRangeUpdate[T int | float64]() *Action[*T] {
	return &Action[*T]{
		Composition: func(f2, f1 *T) *T {
			if f2 != nil {
				return f2
			}
			return f1
		},
		Id: nil,
	}
}

func NewAction[T any](cmp Composition[T], id T) *Action[T] {
	return &Action[T]{
		Composition: cmp,
		Id:          id,
	}
}

// 区間加算・区間和の遅延セグ木用モノイド
func LazyMoRangeAddRangeSum[T int | float64]() *LazyMonoid[T, T] {
	return &LazyMonoid[T, T]{
		Monoid: MoSum[T](),
		Action: ActRangeAdd[T](),
		Mapping: func(f T, x T, length int) T {
			return x + f*T(length)
		},
	}
}

// 区間加算・区間最大値の遅延セグ木用モノイド
func LazyMoRangeAddRangeMax() *LazyMonoid[int, int] {
	return &LazyMonoid[int, int]{
		Monoid: MoMax(),
		Action: ActRangeAdd[int](),
		Mapping: func(f int, x int, length int) int {
			return x + f
		},
	}
}

// 区間加算・区間最小値の遅延セグ木用モノイド
func LazyMoRangeAddRangeMin() *LazyMonoid[int, int] {
	return &LazyMonoid[int, int]{
		Monoid: MoMin(),
		Action: ActRangeAdd[int](),
		Mapping: func(f int, x int, length int) int {
			return x + f
		},
	}
}

// 区間更新・区間和
func LazyMoRangeUpdateRangeSum[T int | float64]() *LazyMonoid[T, *T] {
	return &LazyMonoid[T, *T]{
		Monoid: MoSum[T](),
		Action: ActRangeUpdate[T](),
		Mapping: func(f *T, x T, length int) T {
			if f != nil {
				return (*f) * T(length)
			}
			return x
		},
	}
}

// 区間更新・区間最大値
func LazyMoRangeUpdateRangeMax() *LazyMonoid[int, *int] {
	return &LazyMonoid[int, *int]{
		Monoid: MoMax(),
		Action: ActRangeUpdate[int](),
		Mapping: func(f *int, x int, length int) int {
			if f != nil {
				return *f
			}
			return x
		},
	}
}

// 区間更新・区間最大値
func LazyMoRangeUpdateRangeMin() *LazyMonoid[int, *int] {
	return &LazyMonoid[int, *int]{
		Monoid: MoMin(),
		Action: ActRangeUpdate[int](),
		Mapping: func(f *int, x int, length int) int {
			if f != nil {
				return *f
			}
			return x
		},
	}
}

// カスタム用
func NewLazyMo[S, F any](mo *Monoid[S], act *Action[F], mp Mapping[S, F]) *LazyMonoid[S, F] {
	return &LazyMonoid[S, F]{
		Monoid:  mo,
		Action:  act,
		Mapping: mp,
	}
}

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

// Setは一転更新
func (t *LazySegmentTree[S, F]) Set(i int, x S) {
	i += t.size
	// 遅延をのpush
	for j := t.rank; j >= 1; j-- {
		t.push(i >> j)
	}
	t.data[i] = x
	for i >>= 1; i >= 1; i >>= 1 {
		t.pull(i)
	}
}
