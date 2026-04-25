package acl

// --- Monoid ---

type (
	Operator[T any] func(x1, x2 T) T
	Monoid[T any]   struct {
		Op Operator[T] // データ同士の積演算
		E  T           // データの単位元
	}
)

// MoMax は最大値を求めるモノイド
func MoMax() *Monoid[int] {
	return &Monoid[int]{
		Op: func(x1, x2 int) int {
			return Max(x1, x2)
		},
		E: -INF,
	}
}

// MoMin は最小値を求めるモノイド
func MoMin() *Monoid[int] {
	return &Monoid[int]{
		Op: func(x1, x2 int) int {
			return Min(x1, x2)
		},
		E: INF,
	}
}

// MoSum は和を求めるモノイド
func MoSum[T int | float64]() *Monoid[T] {
	return &Monoid[T]{
		Op: func(x1, x2 T) T {
			return x1 + x2
		},
		E: 0,
	}
}

// MoXOR はXORを求めるモノイド
func MoXOR() *Monoid[int] {
	return &Monoid[int]{
		Op: func(x1, x2 int) int {
			return x1 ^ x2
		},
		E: 0,
	}
}

// MoMODMul はmodを考慮した掛け算を求めるモノイド
func MoMODMul(mod int) *Monoid[int] {
	return &Monoid[int]{
		Op: func(x1, x2 int) int {
			return (x1 * x2) % mod
		},
		E: 1,
	}
}

func NewMo[T any](op Operator[T], e T) *Monoid[T] {
	return &Monoid[T]{
		Op: op,
		E:  e,
	}
}

// --- Action / LazyMonoid ---

type (
	// 作用素同士の合成関数
	Composition[F any] func(f2, f1 F) F

	// 作用素モノイド
	Action[F any] struct {
		Composition Composition[F] // 作用素同士の合成
		Id          F              // 作用素の単位元(恒等写像)
	}

	Mapping[S, F any]    func(f F, x S, length int) S
	LazyMonoid[S, F any] struct {
		*Monoid[S]
		*Action[F]
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

// 区間更新・区間最小値
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
