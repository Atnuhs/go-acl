package main

type Pos struct {
	H, W int
}

func NewPos(h, w int) Pos {
	return Pos{h, w}
}

func (p Pos) Add(dh, dw int) Pos {
	return NewPos(p.H+dh, p.W+dw)
}

func (p Pos) Mul(x int) Pos {
	return NewPos(p.H*x, p.W*x)
}

func (p Pos) Neighbors4() []Pos {
	return []Pos{
		p.Add(-1, 0),
		p.Add(0, -1),
		p.Add(1, 0),
		p.Add(0, 1),
	}
}

func (p Pos) Neighbors8() []Pos {
	return []Pos{
		p.Add(-1, 0),
		p.Add(0, -1),
		p.Add(1, 0),
		p.Add(0, 1),
		p.Add(-1, -1),
		p.Add(-1, 1),
		p.Add(1, -1),
		p.Add(1, 1),
	}
}

func (p Pos) Neighbors4In(h, w int) []Pos {
	ret := make([]Pos, 0, 4)
	for _, p := range p.Neighbors4() {
		if InGrid(p, h, w) {
			ret = append(ret, p)
		}
	}
	return ret
}

func (p Pos) Neighbors8In(h, w int) []Pos {
	ret := make([]Pos, 0, 8)
	for _, p := range p.Neighbors8() {
		if InGrid(p, h, w) {
			ret = append(ret, p)
		}
	}
	return ret
}

func InGrid(p Pos, h, w int) bool {
	return InArea(p, 0, h, 0, w)
}

func InArea(p Pos, hl, hr, wl, wr int) bool {
	return InRange(p.H, hl, hr) && InRange(p.W, wl, wr)
}

func GridAt[T any](p Pos, g [][]T) T {
	return g[p.H][p.W]
}

// FuncSort はPosの配列を、高さ->幅の順でソートする関数を返す
//
// 具体的には、p.Hが違う同士はp.Hが小さい順に並び、p.Hが同じ同士ではp.Wが小さい順に並ぶ
//
//	var ps []Pos
//	SortF(ps, FuncSort(ps)) // ps が破壊的にソートされる
func FuncSort(p []Pos) func(i, j int) bool {
	return func(i, j int) bool {
		if p[i].H != p[j].H {
			return p[i].H < p[j].W
		}
		return p[i].W < p[j].W
	}
}
