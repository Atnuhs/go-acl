package acl

type Abel[T any] struct {
	Add func(v1, v2 T) T
	Sub func(v1, v2 T) T
	E   T
}

func NewIntAbel() *Abel[int] {
	return &Abel[int]{
		Add: func(v1, v2 int) int { return v1 + v2 },
		Sub: func(v1, v2 int) int { return v1 - v2 },
		E:   0,
	}
}

func NewFloat64Abel() *Abel[float64] {
	return &Abel[float64]{
		Add: func(v1, v2 float64) float64 { return v1 + v2 },
		Sub: func(v1, v2 float64) float64 { return v1 - v2 },
		E:   0,
	}
}

func NewXorAbel() *Abel[int] {
	return &Abel[int]{
		Add: func(v1, v2 int) int { return v1 ^ v2 },
		Sub: func(v1, v2 int) int { return v1 ^ v2 }, // XORの逆元は自分自身
		E:   0,
	}
}

func NewModAbel(mod int) *Abel[int] {
	return &Abel[int]{
		Add: func(v1, v2 int) int { return (v1 + v2) % mod },
		Sub: func(v1, v2 int) int { return ((v1-v2)%mod + mod) % mod },
		E:   0,
	}
}

func (a *Abel[T]) Minus(v T) T {
	return a.Sub(a.E, v)
}

// PartialUF
type PartialUF[T any] struct {
	par []int
	dw  []T
	ab  *Abel[T]
}

// NewUnionFind は新しいUnion-Find木を生成する
func NewPartialUF[T any](n int, ab *Abel[T]) *PartialUF[T] {
	par := make([]int, n)
	dw := make([]T, n)
	for i := range par {
		par[i] = -1
		dw[i] = ab.E
	}
	return &PartialUF[T]{
		par: par,
		dw:  dw,
		ab:  ab,
	}
}

func (uf *PartialUF[T]) Root(x int) int {
	if uf.par[x] < 0 {
		return x
	} else {
		p := uf.par[x]
		uf.par[x] = uf.Root(p)
		uf.dw[x] = uf.ab.Add(uf.dw[x], uf.dw[p])
		return uf.par[x]
	}
}

func (uf *PartialUF[T]) Family(x, y int) bool {
	return uf.Root(x) == uf.Root(y)
}

func (uf *PartialUF[T]) Size(x int) int {
	return -uf.par[uf.Root(x)]
}

func (uf *PartialUF[T]) Union(x, y int, w T) {
	w = uf.ab.Sub(uf.ab.Add(w, uf.Weight(x)), uf.Weight(y))

	rx := uf.Root(x)
	ry := uf.Root(y)

	if rx == ry {
		return
	}

	// size(rx) >= size(ry)の状態にする
	if uf.Size(rx) < uf.Size(ry) {
		rx, ry = ry, rx
		w = uf.ab.Minus(w)
	}

	uf.par[rx] += uf.par[ry]
	uf.par[ry] = rx

	uf.dw[ry] = w
}

func (uf *PartialUF[T]) Weight(x int) T {
	uf.Root(x)
	return uf.dw[x]
}

func (uf *PartialUF[T]) Diff(x, y int) T {
	return uf.ab.Sub(uf.Weight(x), uf.Weight(y))
}
