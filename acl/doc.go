// Package acl は AtCoder などの競技プログラミングで利用することを想定した、
// データ構造とアルゴリズムの詰め合わせライブラリです。
//
// 提供される主な機能は以下のとおりです。
//
//   - 入出力ユーティリティ: 高速な標準入力読み取りと出力ヘルパー (io.go)
//   - スライス操作: Map / Filter / Reduce などの汎用ヘルパー (slice.go)
//   - 数値計算: GCD、LCM、モジュラ演算、組合せ (math.go)
//   - データ構造: SegmentTree、LazySegmentTree、BIT、UnionFind、PartialUF、
//     PriorityQueue、DEPQ、Deque、SplayMap、Trie (segtree.go ほか)
//   - グラフ: BFS、DFS、最短路など (graph.go)
//   - 文字列: 文字列処理アルゴリズム (strings.go)
//   - 探索: 二分探索、三分探索など (bisect.go, search.go)
//   - 数論: エラトステネスの篩、素因数分解 (sieve.go)
//   - 座標圧縮、Pair、Pos などの補助型 (compress.go, pair.go, pos.go)
//
// 競技プログラミングという用途上、エラーは可能な限り panic で通知し、
// API はジェネリクスで簡潔に書けるよう設計されています。
package acl
