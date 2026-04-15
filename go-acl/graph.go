package main

import "sort"

// AdjEdge は重み付き隣接リストの要素を表す。
// to は接続先頂点、weight は辺の重み。
type AdjEdge struct {
	To, Weight int
}

// FullEdge は辺集合として扱う重み付き辺を表す。
// Kruskal などで使う。
type FullEdge struct {
	From, To, Weight int
}

// ReadAdjList は頂点数 n, 辺数 m を受け取り、
// 入力された u_i, v_i から隣接リストを構築する。
// directed=false のとき無向グラフとして扱う。
func ReadAdjList(n, m int, directed bool) [][]int {
	g := make([][]int, n)
	for i := 0; i < m; i++ {
		u, v := II()
		u--
		v--

		g[u] = append(g[u], v)
		if !directed {
			g[v] = append(g[v], u)
		}
	}
	return g
}

// ReadWeightedAdjList は頂点数 n, 辺数 m を受け取り、
// 入力された u_i, v_i, w_i から重み付き隣接リストを構築する。
// directed=false のとき無向グラフとして扱う。
func ReadWeightedAdjList(n, m int, directed bool) [][]AdjEdge {
	g := make([][]AdjEdge, n)
	for i := 0; i < m; i++ {
		u, v, w := III()
		u--
		v--

		g[u] = append(g[u], AdjEdge{To: v, Weight: w})
		if !directed {
			g[v] = append(g[v], AdjEdge{To: u, Weight: w})
		}
	}
	return g
}

// Kruskal はクラスカル法を用いて最小全域木を求める。
// 返り値は (最小コスト, 採用した辺一覧)。
// 連結でない場合は (-1, nil) を返す。
func Kruskal(n int, edges []*FullEdge) (int, []*FullEdge) {
	sort.Slice(edges, func(i, j int) bool {
		return edges[i].Weight < edges[j].Weight
	})

	uf := NewUnionFind(n)
	mst := make([]*FullEdge, 0, n-1)
	total := 0

	for _, e := range edges {
		if uf.Family(e.From, e.To) {
			continue
		}
		uf.Union(e.From, e.To)
		mst = append(mst, e)
		total += e.Weight
	}

	if uf.Size(0) != n {
		return -1, nil
	}
	return total, mst
}
