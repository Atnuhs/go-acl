package acl

// Abelはアーベル群(可換群)。PartialUFの重みの代数構造として使う
type Abel[T any] struct {
	Add func(v1, v2 T) T
	Sub func(v1, v2 T) T
	E   T
}

// NewIntAbelは整数の加法アーベル群
func NewIntAbel() *Abel[int] {
	return &Abel[int]{
		Add: func(v1, v2 int) int { return v1 + v2 },
		Sub: func(v1, v2 int) int { return v1 - v2 },
		E:   0,
	}
}

// NewFloat64Abelはfloat64の加法アーベル群
func NewFloat64Abel() *Abel[float64] {
	return &Abel[float64]{
		Add: func(v1, v2 float64) float64 { return v1 + v2 },
		Sub: func(v1, v2 float64) float64 { return v1 - v2 },
		E:   0,
	}
}

// NewXorAbelはXOR群。逆元は自分自身
func NewXorAbel() *Abel[int] {
	return &Abel[int]{
		Add: func(v1, v2 int) int { return v1 ^ v2 },
		Sub: func(v1, v2 int) int { return v1 ^ v2 }, // XORの逆元は自分自身
		E:   0,
	}
}

// NewModAbelはmod modの加法アーベル群
func NewModAbel(mod int) *Abel[int] {
	return &Abel[int]{
		Add: func(v1, v2 int) int { return (v1 + v2) % mod },
		Sub: func(v1, v2 int) int { return ((v1-v2)%mod + mod) % mod },
		E:   0,
	}
}

// Minusはvの逆元(E - v)を返す
func (a *Abel[T]) Minus(v T) T {
	return a.Sub(a.E, v)
}

// PartialUFは重み付きUnion-Find。
// 各要素にアーベル群Tの重みを持たせ、同じ連結成分内で「重みの差」をO(α(N))で取得できる。
type PartialUF[T any] struct {
	par []int // par[x]<0 なら根(値は -size)、それ以外は親
	dw  []T   // 親との重み差 dw[x] = weight[x] - weight[par[x]]
	ab  *Abel[T]
}

// NewPartialUFはn要素の重み付きUnion-Findを生成する
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

// Rootはxの根を返す。経路圧縮しつつdwも根からの累積差に更新する
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

// Familyはxとyが同じ連結成分に属するかを返す
func (uf *PartialUF[T]) Family(x, y int) bool {
	return uf.Root(x) == uf.Root(y)
}

// Sizeはxが属する連結成分のサイズを返す
func (uf *PartialUF[T]) Size(x int) int {
	return -uf.par[uf.Root(x)]
}

// Unionはxとyを連結し、weight[y] - weight[x] = w を制約として課す。
// 既に同じ成分なら何もしない(矛盾検出は呼び出し側でDiffを比較する)。
func (uf *PartialUF[T]) Union(x, y int, w T) {
	// 根同士の重み差に変換: w' = w + weight[x] - weight[y]
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

// Weightはxの根からの累積重みを返す
func (uf *PartialUF[T]) Weight(x int) T {
	uf.Root(x)
	return uf.dw[x]
}

// Diffは Weight(x) - Weight(y) を返す。
// xとyが同じ成分でなければ結果は未定義(呼び出し側でFamilyを確認すること)。
func (uf *PartialUF[T]) Diff(x, y int) T {
	return uf.ab.Sub(uf.Weight(x), uf.Weight(y))
}
