package main

import "sort"

// ReadGraph は頂点数n, 辺数mを受け取り、u_i, v_iの情報からjag配列でグラフを構築する
func ReadGraph(n, m int, direction bool) [][]int {
	g := make([][]int, n)
	for i := 0; i < m; i++ {
		u, v := II()
		u--
		v--
		g[u] = append(g[u], v)
		if !direction {
			g[v] = append(g[v], u)
		}
	}
	return g
}

// WEdge は重み付き辺を表す構造体
type WEdge struct {
	From, To, Weight int
}

// NewWEdge は新しい重み付き辺を生成する
func NewWEdge(from, to, weight int) *WEdge {
	return &WEdge{
		From:   from,
		To:     to,
		Weight: weight,
	}
}

// Kruskal はクラスカル法を用いて最小全域木を求める
// 最小全域木 (MST: Minimum Spanning Tree) とは:
//   - 無向・重み付きグラフの全ての頂点を結びつける木 (辺は n-1 本)
//   - その中で「辺の重み合計」が最小になるもの
//
// 例:
//   - 村を最小コストで道路で全部つなぐ
//   - 島を橋で結んで総工費を最小化
//   - グループをK個に分けたい → MSTを作って重い辺をK−1本切る
//   - 2点間の「経路中の最大コスト」を最小化したい → MST上で見る
func Kruskal(n int, edges []*WEdge) (int, []*WEdge) {
	// はじめに辺を重みでソートする
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	// その後、Union-Findを用いて最小全域木を求める
	uf := NewUnionFind(n)
	ret := make([]*WEdge, 0)
	sum := 0

	// すべての辺を調べる
	for _, e := range edges {
		if uf.Family(e.From, e.To) {
			continue
		}
		ret = append(ret, e)
		sum += e.Weight
		uf.Union(e.From, e.To)
	}
	if uf.Size(0) != n {
		return -1, nil
	}
	return sum, ret
}
