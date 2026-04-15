package main

// UnionFind はUnion-Find木の実装
type UnionFind struct {
	data []int
}

// NewUnionFind は新しいUnion-Find木を生成する
func NewUnionFind(n int) *UnionFind {
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &UnionFind{
		data: data,
	}
}

func (uf *UnionFind) Root(x int) int {
	if uf.data[x] < 0 {
		return x
	} else {
		uf.data[x] = uf.Root(uf.data[x])
		return uf.data[x]
	}
}

func (uf *UnionFind) Family(x, y int) bool {
	return uf.Root(x) == uf.Root(y)
}

func (uf *UnionFind) Size(x int) int {
	return -uf.data[uf.Root(x)]
}

func (uf *UnionFind) Union(x, y int) {
	rx := uf.Root(x)
	ry := uf.Root(y)

	if rx == ry {
		return
	}

	// size(rx) >= size(ry)の状態にする
	if uf.Size(rx) < uf.Size(ry) {
		rx, ry = ry, rx
	}

	uf.data[rx] += uf.data[ry]
	uf.data[ry] = rx
}
